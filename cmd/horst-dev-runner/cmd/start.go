// Copyright © 2017 Tino Rusch
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"context"
	"encoding/json"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/coreos/etcd/clientv3"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start your pipeline",
	Long:  `start your pipeline.`,
	Run: func(cmd *cobra.Command, args []string) {
		/**
		 * Parse config
		 */
		cfg, err := getConfig()
		if err != nil {
			log.Fatal(err)
		}

		/**
		* Create network and start etcd
		 */
		docker, err := client.NewEnvClient()
		if err != nil {
			log.Fatal(err)
		}
		log.Debug("created docker client")
		if err = createNetwork(docker); err != nil {
			log.Fatal(err)
		}
		log.Debug("created horst network")
		var etcdIP string
		if etcdIP, err = startEtcd(docker); err != nil {
			log.Fatal(err)
		}
		log.Debug("started etcd")

		/**
		 * Construct etcd client.
		 * put configs and outputs into etcd.
		 */
		etcd, err := clientv3.New(clientv3.Config{
			Endpoints:   []string{etcdIP + ":2379"},
			DialTimeout: 3 * time.Second,
		})
		if err != nil {
			log.Fatal(err)
		}
		log.Debug("connected to etcd")
		for id, config := range cfg {
			if err = prepareComponent(id, config, etcd); err != nil {
				log.Fatal(err)
			}
			log.Debug("wrote config and outputs of ", id)
		}

		/**
		 * kick off containers
		 */
		for id, config := range cfg {
			if err = startComponent(id, config.Image, docker, etcdIP); err != nil {
				log.Fatal(err)
			}
			log.Debugf("started component '%v' (%v)", id, config.Image)
		}
	},
}

func init() {
	RootCmd.AddCommand(startCmd)
}

func prepareComponent(id string, config *ComponentConfig, etcd *clientv3.Client) error {
	cfgBytes, err := json.Marshal(config.Config)
	if err != nil {
		return err
	}
	outputBytes, err := json.Marshal(config.Outputs)
	if err != nil {
		return err
	}
	_, err = etcd.Put(context.Background(), "/horst/configs/"+id, string(cfgBytes))
	if err != nil {
		return err
	}
	_, err = etcd.Put(context.Background(), "/horst/outputs/"+id, string(outputBytes))
	if err != nil {
		return err
	}
	return nil
}

func createNetwork(docker *client.Client) error {
	name := "horstnet"
	_, err := docker.NetworkInspect(context.Background(), name)
	if err != nil {
		log.Debug("horst network not found, creating new")
		_, err = docker.NetworkCreate(context.Background(), name, types.NetworkCreate{Driver: "bridge"})
		return err
	}
	return nil
}

func startEtcd(docker *client.Client) (string, error) {
	containerConfig := &container.Config{
		Image: "quay.io/coreos/etcd",
		Cmd: []string{
			"/usr/local/bin/etcd",
			"--name", "node1",
			"--data-dir", "/data",
			"--initial-advertise-peer-urls", "http://0.0.0.0:2380",
			"--listen-peer-urls", "http://0.0.0.0:2380",
			"--advertise-client-urls", "http://0.0.0.0:2379",
			"--listen-client-urls", "http://0.0.0.0:2379",
			"--initial-cluster", "node1=http://0.0.0.0:2380",
		},
	}
	hostConfig := &container.HostConfig{
		AutoRemove: true,
	}

	net := "horstnet"
	networkConfig := &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			net: &network.EndpointSettings{},
		},
	}

	createResp, err := docker.ContainerCreate(context.Background(), containerConfig, hostConfig, networkConfig, "etcd")
	if err != nil {
		return "", err
	}
	if err = docker.ContainerStart(context.Background(), createResp.ID, types.ContainerStartOptions{}); err != nil {
		return "", err
	}
	return getEtcdIP(docker)
}

func getEtcdIP(docker *client.Client) (string, error) {
	res, err := docker.ContainerInspect(context.Background(), "etcd")
	if err != nil {
		return "", err
	}
	return res.NetworkSettings.Networks["horstnet"].IPAddress, nil
}

func startComponent(id, imageID string, docker *client.Client, etcdIP string) error {
	containerConfig := &container.Config{
		Image: imageID,
		Cmd: []string{
			"--id", id,
			"--etcd", etcdIP + ":2379",
		},
	}
	hostConfig := &container.HostConfig{
		AutoRemove: true,
	}
	networkConfig := &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			"horstnet": &network.EndpointSettings{},
		},
	}
	createResp, err := docker.ContainerCreate(context.Background(), containerConfig, hostConfig, networkConfig, id)
	if err != nil {
		return err
	}
	return docker.ContainerStart(context.Background(), createResp.ID, types.ContainerStartOptions{})
}

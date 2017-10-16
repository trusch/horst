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
	"io/ioutil"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/coreos/etcd/clientv3"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start your pipeline",
	Long:  `start your pipeline.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.DebugLevel)

		/**
		 * Parse config
		 */
		cfg := make(map[string]*ComponentConfig)
		bs, err := ioutil.ReadFile(viper.GetString("config"))
		if err != nil {
			log.Fatal(err)
		}
		if err = json.Unmarshal(bs, &cfg); err != nil {
			log.Fatal(err)
		}
		log.Debug("parsed config")

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

// ComponentConfig is the config of a single component instance
type ComponentConfig struct {
	Image   string            `json:"image"`
	Config  interface{}       `json:"config"`
	Outputs map[string]string `json:"outputs"`
}

func prepareComponent(id string, config *ComponentConfig, etcd *clientv3.Client) error {
	ns := viper.GetString("namespace")
	id = ns + id
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
	ns := viper.GetString("namespace")
	name := ns + "horstnet"
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

	ns := viper.GetString("namespace")
	net := ns + "horstnet"
	networkConfig := &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			net: &network.EndpointSettings{},
		},
	}

	createResp, err := docker.ContainerCreate(context.Background(), containerConfig, hostConfig, networkConfig, ns+"etcd")
	if err != nil {
		return "", err
	}
	if err = docker.ContainerStart(context.Background(), createResp.ID, types.ContainerStartOptions{}); err != nil {
		return "", err
	}
	res, err := docker.ContainerInspect(context.Background(), ns+"etcd")
	if err != nil {
		return "", err
	}
	return res.NetworkSettings.Networks[net].IPAddress, nil
}

func startComponent(id, imageID string, docker *client.Client, etcdIP string) error {
	ns := viper.GetString("namespace")
	id = ns + id

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
			ns + "horstnet": &network.EndpointSettings{},
		},
	}
	createResp, err := docker.ContainerCreate(context.Background(), containerConfig, hostConfig, networkConfig, id)
	if err != nil {
		return err
	}
	return docker.ContainerStart(context.Background(), createResp.ID, types.ContainerStartOptions{})
}

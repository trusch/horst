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
//
// ACPI output while writing this file:
// > Battery 0: Discharging, 0%, discharging at zero rate - will never fully discharge.

package cmd

import (
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/coreos/etcd/clientv3"
	log "github.com/sirupsen/logrus"

	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// reconfigureCmd represents the start command
var reconfigureCmd = &cobra.Command{
	Use:   "reconfigure",
	Short: "reconfigure",
	Long:  `reconfigure.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.DebugLevel)

		cfg := make(map[string]*ComponentConfig)
		bs, err := ioutil.ReadFile(viper.GetString("config"))
		if err != nil {
			log.Fatal(err)
		}
		if err = json.Unmarshal(bs, &cfg); err != nil {
			log.Fatal(err)
		}
		log.Debug("parsed config")

		docker, err := client.NewEnvClient()
		if err != nil {
			log.Fatal(err)
		}

		etcdIP, err := getEtcdIP(docker)
		if err != nil {
			log.Fatal(err)
		}
		etcd, err := clientv3.New(clientv3.Config{
			Endpoints:   []string{etcdIP + ":2379"},
			DialTimeout: 3 * time.Second,
		})
		if err != nil {
			log.Fatal(err)
		}

		if len(args) > 0 {
			for _, component := range args {
				if err = prepareComponent(component, cfg[component], etcd); err != nil {
					log.Fatal(err)
				}
				log.Debugf("wrote config and outputs for %v", component)
			}
			return
		}

		for id, config := range cfg {
			if err = prepareComponent(id, config, etcd); err != nil {
				log.Fatal(err)
			}
			log.Debugf("wrote config and outputs for %v", id)
		}
	},
}

func init() {
	RootCmd.AddCommand(reconfigureCmd)
}

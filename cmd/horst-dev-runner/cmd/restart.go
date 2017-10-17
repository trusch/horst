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

	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// restartCmd represents the start command
var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "restart",
	Long:  `restart.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.DebugLevel)
		docker, err := client.NewEnvClient()
		if err != nil {
			log.Fatal(err)
		}

		if len(args) > 0 {
			for _, component := range args {
				if err = restart(component, docker); err != nil {
					log.Fatal(err)
				}
			}
			return
		}

		cfg := make(map[string]*ComponentConfig)
		bs, err := ioutil.ReadFile(viper.GetString("config"))
		if err != nil {
			log.Fatal(err)
		}
		if err = json.Unmarshal(bs, &cfg); err != nil {
			log.Fatal(err)
		}
		log.Debug("parsed config")

		for id, config := range cfg {
			if err = restart(id, docker); err != nil {
				log.Fatal(err)
			}
			log.Debugf("restarted component '%v' (%v)", id, config.Image)
		}
	},
}

func restart(component string, docker *client.Client) error {
	timeout := 3 * time.Second
	if err := docker.ContainerRestart(context.Background(), component, &timeout); err != nil {
		return err
	}
	return nil
}

func init() {
	RootCmd.AddCommand(restartCmd)
}

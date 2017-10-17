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
	"time"

	"github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "stop your pipeline",
	Long:  `stop your pipeline.`,
	Run: func(cmd *cobra.Command, args []string) {
		/**
		 * Parse config
		 */
		cfg, err := getConfig()
		if err != nil {
			log.Fatal(err)
		}

		/**
		 * stop containers
		 */
		docker, err := client.NewEnvClient()
		if err != nil {
			log.Fatal(err)
		}
		log.Debug("created docker client")
		timeout := 5 * time.Second
		for id := range cfg {
			if err := docker.ContainerStop(context.Background(), id, &timeout); err != nil {
				log.Warn(err)
			}
			log.Debug("stopped ", id)
		}
		if err := docker.ContainerStop(context.Background(), "etcd", &timeout); err != nil {
			log.Warn(err)
		}
		log.Debug("stopped etcd")
	},
}

func init() {
	RootCmd.AddCommand(stopCmd)
}

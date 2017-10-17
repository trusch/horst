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
	"context"
	"os"
	"sort"
	"time"

	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"

	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

// statusCmd represents the start command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "status",
	Long:  `status.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := getConfig()
		if err != nil {
			log.Fatal(err)
		}

		docker, err := client.NewEnvClient()
		if err != nil {
			log.Fatal(err)
		}

		if len(args) > 0 {
			for _, component := range args {
				if err = status(component, docker); err != nil {
					log.Fatal(err)
				}
			}
			return
		}

		components := make([]string, len(cfg))
		i := 0
		for id := range cfg {
			components[i] = id
			i++
		}
		sort.Strings(components)

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "Status", "Running"})
		for _, id := range components {
			resp, err := docker.ContainerInspect(context.Background(), id)
			if err != nil {
				table.Append([]string{id, "stopped", ""})
				continue
			}
			startedAt, err := time.Parse(time.RFC3339, resp.State.StartedAt)
			if err != nil {
				log.Error(err)
				continue
			}
			table.Append([]string{id, resp.State.Status, time.Since(startedAt).Round(time.Second).String()})
		}
		table.Render()
	},
}

func status(component string, docker *client.Client) error {
	return nil
}

func init() {
	RootCmd.AddCommand(statusCmd)
}

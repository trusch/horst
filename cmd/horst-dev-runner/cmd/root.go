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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "dev-runner",
	Short: "start and stop your pipeline",
	Long:  `start and stop your pipeline.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringP("config", "c", "config.json", "config file")
	RootCmd.PersistentFlags().String("log-level", "INFO", "log level")
	viper.BindPFlags(RootCmd.PersistentFlags())
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.AutomaticEnv() // read in environment variables that match
	logLevel := viper.GetString("log-level")
	switch logLevel {
	case "DEBUG":
		log.SetLevel(log.DebugLevel)
	case "ERROR":
		log.SetLevel(log.ErrorLevel)
	case "WARN":
		log.SetLevel(log.WarnLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
}

// ComponentConfig is the config of a single component instance
type ComponentConfig struct {
	Image   string            `json:"image"`
	Config  interface{}       `json:"config"`
	Outputs map[string]string `json:"outputs"`
}

func getConfig() (map[string]*ComponentConfig, error) {
	cfg := make(map[string]*ComponentConfig)
	bs, err := ioutil.ReadFile(viper.GetString("config"))
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(bs, &cfg); err != nil {
		return nil, err
	}
	log.Debug("parsed config")
	return cfg, nil
}

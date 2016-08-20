package config

import (
	"encoding/json"
	"errors"
	"os"
	"strings"

	"github.com/trusch/horst/links"
)

// NodeConfig is the configuration for one processor instance
type NodeConfig struct {
	ID      string            `json:"id"`
	Class   string            `json:"class"`
	Config  interface{}       `json:"config"`
	Outputs map[string]string `json:"outputs"`
}

// Config is the config of all processors in a pipeline
type Config map[string]NodeConfig

// Load loads a config from file
func (cfg Config) Load(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(f)
	return decoder.Decode(&cfg)
}

// Save saves a config to file
func (cfg Config) Save(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(f)
	return encoder.Encode(cfg)
}

// GetLinkMap returns the resulting linkmap of this config
func (cfg Config) GetLinkMap() (links.LinkMap, error) {
	res := links.LinkMap{}
	for proc, nodeCfg := range cfg {
		for output, target := range nodeCfg.Outputs {
			targetInstance := ""
			targetInput := "default"
			targetParts := strings.Split(target, ":")
			switch len(targetParts) {
			case 1:
				targetInstance = targetParts[0]
			case 2:
				targetInstance = targetParts[0]
				targetInput = targetParts[1]
			default:
				return nil, errors.New("malformed target")
			}
			res.Add(proc, output, targetInstance, targetInput)
		}
	}
	return res, nil
}

// New loads a new config from file
func New(file string) (Config, error) {
	cfg := Config{}
	return cfg, cfg.Load(file)
}

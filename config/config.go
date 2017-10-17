package config

// ComponentConfig is the config of a single component instance
type ComponentConfig struct {
	Image   string            `json:"image"`
	Config  interface{}       `json:"config"`
	Outputs map[string]string `json:"outputs"`
}

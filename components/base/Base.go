package base

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/trusch/horst/components"
)

// Component is the most basic component we can build
type Component struct {
	Config  map[string]interface{}
	Outputs map[string]string
}

// New returns a new base.Component
func New() (components.Component, error) {
	return &Component{}, nil
}

// HandleConfigUpdate gets called when new config for this component is available
func (c *Component) HandleConfigUpdate(config map[string]interface{}) error {
	c.Config = config
	return nil
}

// HandleOutputUpdate gets called when new output endpoints are available
func (c *Component) HandleOutputUpdate(outputs map[string]string) error {
	c.Outputs = outputs
	return nil
}

// Process gets called when a new event for a specific input should be processed
func (c *Component) Process(inputID string, event interface{}) error {
	log.Printf("%v: %v", inputID, event)
	return nil
}

// Emit sends a event to a specific output
func (c *Component) Emit(outputID string, event interface{}) error {
	if targetURI, ok := c.Outputs[outputID]; ok {
		r, w := io.Pipe()
		encoder := json.NewEncoder(w)
		go func() {
			encoder.Encode(event)
			w.Close()
		}()
		resp, err := http.Post(targetURI, "application/json", r)
		if err != nil {
			return err
		}
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("http error: %v", resp.StatusCode)
		}
		return nil
	}
	return errors.New("no such output")
}

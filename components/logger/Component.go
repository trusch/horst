package logger

import (
	"log"

	"github.com/trusch/horst/components"
	"github.com/trusch/horst/components/base"
)

// Component is the most basic component we can build
type Component struct {
	base.Component
}

// New returns a new logger.Component
func New() (components.Component, error) {
	return &Component{}, nil
}

// Process gets called when a new event for a specific input should be processed
func (c *Component) Process(inputID string, event interface{}) error {
	log.Printf("%v: %v", inputID, event)
	return nil
}

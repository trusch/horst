package textsplitter

import (
	"log"
	"strings"

	"github.com/trusch/horst/components"
	"github.com/trusch/horst/components/base"
)

// Component is the most basic component we can build
type Component struct {
	base.Component
}

// New returns a new textsplitter.Component
func New() (components.Component, error) {
	return &Component{}, nil
}

// Process gets called when a new event for a specific input should be processed
func (c *Component) Process(inputID string, event interface{}) error {
	if str, ok := event.(string); ok {
		for _, word := range strings.Fields(str) {
			for id := range c.Outputs {
				if err := c.Emit(id, word); err != nil {
					log.Print(err)
				}
			}
		}
	}
	return nil
}

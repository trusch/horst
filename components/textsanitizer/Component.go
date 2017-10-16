package textsanitizer

import (
	"log"
	"strings"

	"github.com/trusch/horst/components"
	"github.com/trusch/horst/components/base"
)

// Component is the most basic component we can build
type Component struct {
	base.Component
	cutset  string
	toLower bool
	toUpper bool
}

// New returns a new textsanitizer.Component
func New() (components.Component, error) {
	return &Component{}, nil
}

// HandleConfigUpdate gets called when new config for this component is available
func (c *Component) HandleConfigUpdate(config map[string]interface{}) error {
	if cutset, ok := config["cutset"].(string); ok {
		c.cutset = cutset
	}
	if toLower, ok := config["toLower"].(bool); ok {
		c.toLower = toLower
	}
	if toUpper, ok := config["toUpper"].(bool); ok {
		c.toUpper = toUpper
	}
	return nil
}

// Process gets called when a new event for a specific input should be processed
func (c *Component) Process(inputID string, event interface{}) error {
	if str, ok := event.(string); ok {
		str = c.sanitize(str)
		for id := range c.Outputs {
			if err := c.Emit(id, str); err != nil {
				log.Print(err)
			}
		}
	}
	return nil
}

func (c *Component) sanitize(input string) string {
	if c.toLower {
		input = strings.ToLower(input)
	}
	if c.toUpper {
		input = strings.ToUpper(input)
	}
	for _, c := range c.cutset {
		input = strings.Replace(input, string(c), "", -1)
	}
	return input
}

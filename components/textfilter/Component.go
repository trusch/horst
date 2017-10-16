package textfilter

import (
	"errors"
	"log"
	"regexp"

	"github.com/trusch/horst/components"
	"github.com/trusch/horst/components/base"
)

// Component is the most basic component we can build
type Component struct {
	base.Component
	filter *regexp.Regexp
}

// New returns a new textfilter.Component
func New() (components.Component, error) {
	return &Component{}, nil
}

// HandleConfigUpdate gets called when new config for this component is available
func (c *Component) HandleConfigUpdate(config map[string]interface{}) error {
	if regexpStr, ok := config["regexp"].(string); ok {
		f, err := regexp.Compile(regexpStr)
		if err != nil {
			return err
		}
		c.filter = f
		return nil
	}
	return errors.New("require valid regexp field in config")
}

// Process gets called when a new event for a specific input should be processed
func (c *Component) Process(inputID string, event interface{}) error {
	if str, ok := event.(string); ok {
		if c.filter.MatchString(str) {
			for id := range c.Outputs {
				if err := c.Emit(id, str); err != nil {
					log.Print(err)
				}
			}
		}
	}
	return nil
}

package projector

import (
	"bytes"
	"encoding/json"
	"log"
	"strings"

	"github.com/trusch/horst/components"
	"github.com/trusch/horst/components/base"
	"github.com/trusch/jsonq"
)

// Component is the most basic component we can build
type Component struct {
	base.Component
}

// New returns a new projector.Component
func New() (components.Component, error) {
	return &Component{}, nil
}

// Process gets called when a new event for a specific input should be processed
func (c *Component) Process(inputID string, event interface{}) error {
	if doc, ok := event.(map[string]interface{}); ok {
		template := c.copyProjectionConfig()
		res := c.searchAndReplaceKeys(template, doc)
		for id := range c.Outputs {
			if err := c.Emit(id, res); err != nil {
				log.Print(err)
			}
		}
	}
	return nil
}

func (c *Component) copyProjectionConfig() interface{} {
	if str, ok := c.Config["projection"]; ok {
		return str
	}
	if obj, ok := c.Config["projection"].(map[string]interface{}); ok {
		result := make(map[string]interface{})
		var mod bytes.Buffer
		enc := json.NewEncoder(&mod)
		dec := json.NewDecoder(&mod)
		enc.Encode(obj)
		dec.Decode(&result)
		return result
	}
	return nil
}

func (c *Component) searchAndReplaceKeys(arg interface{}, inputDoc map[string]interface{}) interface{} {
	if argStr, ok := arg.(string); ok && argStr[0] == '@' {
		parts := strings.Split(argStr[1:], ".")
		jq := jsonq.NewQuery(inputDoc)
		val, _ := jq.Get(parts...)
		return val
	}
	if argMap, ok := arg.(map[string]interface{}); ok {
		for k, v := range argMap {
			argMap[k] = c.searchAndReplaceKeys(v, inputDoc)
		}
	}
	return arg
}

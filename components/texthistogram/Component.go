package texthistogram

import (
	"log"

	"github.com/trusch/horst/components"
	"github.com/trusch/horst/components/base"
)

// Component is the most basic component we can build
type Component struct {
	base.Component
	decay     float64
	min       float64
	minEmit   float64
	relative  bool
	emitEvery int64
	counter   int64
	histogram map[string]float64
}

// New returns a new texthistogram.Component
func New() (components.Component, error) {
	return &Component{}, nil
}

// HandleConfigUpdate gets called when new config for this component is available
func (c *Component) HandleConfigUpdate(config map[string]interface{}) error {
	if decay, ok := config["decay"].(float64); ok {
		c.decay = decay
	}
	if min, ok := config["min"].(float64); ok {
		c.min = min
	}
	if minEmit, ok := config["minEmit"].(float64); ok {
		c.minEmit = minEmit
	}
	if relative, ok := config["relative"].(bool); ok {
		c.relative = relative
	}
	if emitEvery, ok := config["emitEvery"].(float64); ok {
		c.emitEvery = int64(emitEvery)
	} else {
		c.emitEvery = 1
	}
	c.histogram = make(map[string]float64)
	return nil
}

// Process gets called when a new event for a specific input should be processed
func (c *Component) Process(inputID string, event interface{}) error {
	if str, ok := event.(string); ok {
		c.update(str)
		c.counter++
		c.counter %= c.emitEvery
		if c.counter == 0 {
			outputHistogram := c.copy()
			for id := range c.Outputs {
				if err := c.Emit(id, outputHistogram); err != nil {
					log.Print(err)
				}
			}
		}
	}
	return nil
}

func (c *Component) update(input string) {
	for key := range c.histogram {
		c.histogram[key] *= c.decay
		if c.histogram[key] < c.min {
			delete(c.histogram, key)
		}
	}
	c.histogram[input] += 1.0
}

func (c *Component) copy() map[string]float64 {
	res := make(map[string]float64)
	var sum float64
	for k, v := range c.histogram {
		if v >= c.minEmit {
			res[k] = v
			sum += v
		}
	}
	if c.relative {
		for k, v := range res {
			res[k] = v / sum
		}
	}
	return res
}

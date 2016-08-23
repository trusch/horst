package incrementalhistogram

import (
	"fmt"

	"github.com/trusch/horst"
	"github.com/trusch/horst/processors/base"
	"github.com/trusch/horst/registry"
)

type incrementalhistogramType struct {
	base.Base
	data      chan incrementalhistogramMessage
	decay     float64
	min       float64
	minEmit   float64
	relative  bool
	emitEvery int64

	histogram map[string]float64
}

type incrementalhistogramMessage struct {
	input string
	data  interface{}
}

func (histogram *incrementalhistogramType) parseConfig() error {
	if cfg, ok := histogram.Config.(map[string]interface{}); ok {
		if decay, ok := cfg["decay"].(float64); ok {
			histogram.decay = decay
		}
		if min, ok := cfg["min"].(float64); ok {
			histogram.min = min
		}
		if minEmit, ok := cfg["minEmit"].(float64); ok {
			histogram.minEmit = minEmit
		}
		if relative, ok := cfg["relative"].(bool); ok {
			histogram.relative = relative
		}
		if emitEvery, ok := cfg["emitEvery"].(float64); ok {
			histogram.emitEvery = int64(emitEvery)
		} else {
			histogram.emitEvery = 1
		}
	} else {
		return fmt.Errorf("config needs to be an object with 'decay':float, 'min':float, 'relative':bool")
	}
	return nil
}

func (histogram *incrementalhistogramType) update(input string) {
	for key := range histogram.histogram {
		histogram.histogram[key] *= histogram.decay
		if histogram.histogram[key] < histogram.min {
			delete(histogram.histogram, key)
		}
	}
	histogram.histogram[input] += 1.0
}

func (histogram *incrementalhistogramType) backend() {
	var i int64
	for msg := range histogram.data {
		if str, ok := msg.data.(string); ok {
			histogram.update(str)
			i++
			i %= histogram.emitEvery
			if i == 0 {
				histogram.Manager.Emit(histogram.ID, "out", histogram.copy())
			}
		}
	}
}

func (histogram *incrementalhistogramType) copy() map[string]interface{} {
	res := make(map[string]interface{})
	var sum float64
	for k, v := range histogram.histogram {
		if v >= histogram.minEmit {
			res[k] = v
			sum += v
		}
	}
	if histogram.relative {
		for k, v := range res {
			res[k] = v.(float64) / sum
		}
	}
	return res
}

func (histogram *incrementalhistogramType) Process(in string, data interface{}) {
	histogram.data <- incrementalhistogramMessage{in, data}
}

func (histogram *incrementalhistogramType) Stop() {
	close(histogram.data)
}

func init() {
	registry.Register("github.com/trusch/horst/processors/incrementalhistogram", func(id string, config interface{}, mgr horst.ProcessorManager) (horst.Processor, error) {
		incrementalhistogram := &incrementalhistogramType{
			data:      make(chan incrementalhistogramMessage, 32),
			histogram: make(map[string]float64),
		}
		incrementalhistogram.InitBase(id, config, mgr)
		err := incrementalhistogram.parseConfig()
		if err != nil {
			return nil, err
		}
		go incrementalhistogram.backend()
		return incrementalhistogram, nil
	})
}

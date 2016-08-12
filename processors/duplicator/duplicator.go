package duplicator

import (
	"fmt"
	"strconv"

	"github.com/trusch/horst"
	"github.com/trusch/horst/processors/base"
	"github.com/trusch/horst/registry"
)

type duplicatorType struct {
	base.Base
	outputs int
	data    chan duplicatorMessage
}

type duplicatorMessage struct {
	input string
	data  interface{}
}

func (duplicator *duplicatorType) backend() {
	for msg := range duplicator.data {
		for i := 0; i < duplicator.outputs; i++ {
			duplicator.Manager.Emit(duplicator.ID, "out"+strconv.Itoa(i), msg.data)
		}
	}
}

func (duplicator *duplicatorType) Process(in string, data interface{}) {
	duplicator.data <- duplicatorMessage{in, data}
}

func (duplicator *duplicatorType) Stop() {
	close(duplicator.data)
}

func init() {
	registry.Register("duplicator", func(id string, config interface{}, mgr horst.ProcessorManager) (horst.Processor, error) {
		duplicator := &duplicatorType{data: make(chan duplicatorMessage, 32)}
		duplicator.InitBase(id, config, mgr)
		numOfOutputs, ok := config.(float64)
		if !ok {
			return nil, fmt.Errorf("duplicator(%v)> config is not a number, it is %T", duplicator.ID, config)
		}
		duplicator.outputs = int(numOfOutputs)
		go duplicator.backend()
		return duplicator, nil
	})
}

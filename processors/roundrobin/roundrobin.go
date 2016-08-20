package roundrobin

import (
	"fmt"
	"strconv"

	"github.com/trusch/horst"
	"github.com/trusch/horst/processors/base"
	"github.com/trusch/horst/registry"
)

type roundrobinType struct {
	base.Base
	outputs int
	data    chan roundrobinMessage
}

type roundrobinMessage struct {
	input string
	data  interface{}
}

func (roundrobin *roundrobinType) backend() {
	current := 0
	for msg := range roundrobin.data {
		roundrobin.Manager.Emit(roundrobin.ID, "out"+strconv.Itoa(current), msg.data)
		current = (current + 1) % roundrobin.outputs
	}
}

func (roundrobin *roundrobinType) Process(in string, data interface{}) {
	roundrobin.data <- roundrobinMessage{in, data}
}

func (roundrobin *roundrobinType) Stop() {
	close(roundrobin.data)
}

func init() {
	registry.Register("github.com/trusch/horst/processors/roundrobin", func(id string, config interface{}, mgr horst.ProcessorManager) (horst.Processor, error) {
		roundrobin := &roundrobinType{data: make(chan roundrobinMessage, 32)}
		roundrobin.InitBase(id, config, mgr)
		numOfOutputs, ok := config.(float64)
		if !ok {
			return nil, fmt.Errorf("roundrobin(%v)> config is not a number, it is %T", roundrobin.ID, config)
		}
		roundrobin.outputs = int(numOfOutputs)
		go roundrobin.backend()
		return roundrobin, nil
	})
}

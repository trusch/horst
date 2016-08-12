package naturals

import (
	"time"

	"github.com/trusch/horst"
	"github.com/trusch/horst/registry"
)

type naturalsType struct {
	id  string
	mgr horst.ProcessorManager
	run bool
}

func (naturals *naturalsType) backend() {
	num := 0
	for naturals.run {
		naturals.mgr.Emit(naturals.id, "out", num)
		time.Sleep(1 * time.Second)
		num++
	}
}

func (naturals *naturalsType) Process(in string, data interface{}) {}

func (naturals *naturalsType) Stop() {
	naturals.run = false
}

func init() {
	registry.Register("naturals", func(id string, config interface{}, mgr horst.ProcessorManager) (horst.Processor, error) {
		naturals := &naturalsType{id, mgr, true}
		go naturals.backend()
		return naturals, nil
	})
}

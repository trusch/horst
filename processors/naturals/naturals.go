package naturals

import (
	"time"

	"github.com/trusch/horst"
	"github.com/trusch/horst/processors/base"
	"github.com/trusch/horst/registry"
)

type naturalsType struct {
	base.Base
	run bool
}

func (naturals *naturalsType) backend() {
	num := 0
	for naturals.run {
		naturals.Manager.Emit(naturals.ID, "out", num)
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
		naturals := &naturalsType{run: true}
		naturals.InitBase(id, config, mgr)
		go naturals.backend()
		return naturals, nil
	})
}

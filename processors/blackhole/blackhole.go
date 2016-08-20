package blackhole

import (
	"github.com/trusch/horst"
	"github.com/trusch/horst/processors/base"
	"github.com/trusch/horst/registry"
)

type blackholeType struct {
	base.Base
}

func (naturals *blackholeType) Process(in string, data interface{}) {}
func (naturals *blackholeType) Stop()                               {}

func init() {
	registry.Register("github.com/trusch/horst/processors/blackhole", func(id string, config interface{}, mgr horst.ProcessorManager) (horst.Processor, error) {
		hole := &blackholeType{}
		hole.InitBase(id, config, mgr)
		return hole, nil
	})
}

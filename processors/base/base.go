package base

import "github.com/trusch/horst"

// Base is a base processor, and provides a trivial init method
type Base struct {
	ID      string
	Config  interface{}
	Manager horst.ProcessorManager
}

// InitBase initializes the processor with basic info
func (base *Base) InitBase(id string, config interface{}, mgr horst.ProcessorManager) {
	base.ID = id
	base.Config = config
	base.Manager = mgr
}

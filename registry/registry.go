package registry

import (
	"errors"

	"github.com/trusch/horst"
)

type constructor func(id string, config interface{}, mgr horst.ProcessorManager) (horst.Processor, error)

var constructors = make(map[string]constructor)

// Register registers a processor type globally
func Register(className string, ctor constructor) {
	constructors[className] = ctor
}

// Construct returns a new processor instance for a given type
func Construct(className, id string, config interface{}, mgr horst.ProcessorManager) (horst.Processor, error) {
	if ctor, ok := constructors[className]; ok {
		return ctor(id, config, mgr)
	}
	return nil, errors.New("no such processor type")
}

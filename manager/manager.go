package manager

import (
	"github.com/trusch/horst"
	"github.com/trusch/horst/links"
)

type manager struct {
	processors map[string]horst.Processor
	links      links.LinkMap
}

func (mgr *manager) Init(links links.LinkMap) {
	mgr.links = links
	mgr.processors = make(map[string]horst.Processor)
}

func (mgr *manager) AddProcessor(processorID string, processor horst.Processor) {
	mgr.processors[processorID] = processor
}

func (mgr *manager) DelProcessor(processorID string) {
	mgr.processors[processorID].Stop()
	delete(mgr.processors, processorID)
}

func (mgr *manager) Emit(fromProcessor, fromProcessorOutput string, data interface{}) {
	if link, ok := mgr.links.Get(fromProcessor, fromProcessorOutput); ok {
		mgr.Process(link.To, link.ToInput, data)
	}
}

func (mgr *manager) Process(processorID, inputID string, data interface{}) {
	if proc, ok := mgr.processors[processorID]; ok {
		proc.Process(inputID, data)
	}
}

// New constructs a new processor manager
func New(links links.LinkMap) horst.ProcessorManager {
	return &manager{
		processors: make(map[string]horst.Processor),
		links:      links,
	}
}

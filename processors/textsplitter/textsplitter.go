package textsplitter

import (
	"strings"

	"github.com/trusch/horst"
	"github.com/trusch/horst/processors/base"
	"github.com/trusch/horst/registry"
)

type textsplitterType struct {
	base.Base
	data    chan textsplitterMessage
	cutset  string
	toLower bool
	toUpper bool
}

type textsplitterMessage struct {
	input string
	data  interface{}
}

func (textsplitter *textsplitterType) backend() {
	for msg := range textsplitter.data {
		if str, ok := msg.data.(string); ok {
			for _, word := range strings.Fields(str) {
				textsplitter.Manager.Emit(textsplitter.ID, "out", word)
			}
		}
	}
}

func (textsplitter *textsplitterType) Process(in string, data interface{}) {
	textsplitter.data <- textsplitterMessage{in, data}
}

func (textsplitter *textsplitterType) Stop() {
	close(textsplitter.data)
}

func init() {
	registry.Register("github.com/trusch/horst/processors/textsplitter", func(id string, config interface{}, mgr horst.ProcessorManager) (horst.Processor, error) {
		textsplitter := &textsplitterType{data: make(chan textsplitterMessage, 32)}
		textsplitter.InitBase(id, config, mgr)
		go textsplitter.backend()
		return textsplitter, nil
	})
}

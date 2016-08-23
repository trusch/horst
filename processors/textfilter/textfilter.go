package textfilter

import (
	"fmt"
	"regexp"

	"github.com/trusch/horst"
	"github.com/trusch/horst/processors/base"
	"github.com/trusch/horst/registry"
)

type textfilterType struct {
	base.Base
	data  chan textfilterMessage
	regex *regexp.Regexp
}

type textfilterMessage struct {
	input string
	data  interface{}
}

func (textfilter *textfilterType) parseConfig() error {
	if cfg, ok := textfilter.Config.(map[string]interface{}); ok {
		if regex, ok := cfg["regex"].(string); ok {
			textfilter.regex = regexp.MustCompile(regex)
		}
	} else {
		return fmt.Errorf("config needs to be an object with 'regex':string")
	}
	return nil
}

func (textfilter *textfilterType) backend() {
	for msg := range textfilter.data {
		if str, ok := msg.data.(string); ok {
			if textfilter.regex.MatchString(str) {
				textfilter.Manager.Emit(textfilter.ID, "out", str)
			}
		}
	}
}

func (textfilter *textfilterType) Process(in string, data interface{}) {
	textfilter.data <- textfilterMessage{in, data}
}

func (textfilter *textfilterType) Stop() {
	close(textfilter.data)
}

func init() {
	registry.Register("github.com/trusch/horst/processors/textfilter", func(id string, config interface{}, mgr horst.ProcessorManager) (horst.Processor, error) {
		textfilter := &textfilterType{data: make(chan textfilterMessage, 32)}
		textfilter.InitBase(id, config, mgr)
		err := textfilter.parseConfig()
		if err != nil {
			return nil, err
		}
		go textfilter.backend()
		return textfilter, nil
	})
}

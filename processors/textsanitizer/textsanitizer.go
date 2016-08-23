package textsanitizer

import (
	"fmt"
	"strings"

	"github.com/trusch/horst"
	"github.com/trusch/horst/processors/base"
	"github.com/trusch/horst/registry"
)

type textsanitizerType struct {
	base.Base
	data    chan textsanitizerMessage
	cutset  string
	toLower bool
	toUpper bool
}

type textsanitizerMessage struct {
	input string
	data  interface{}
}

func (textsanitizer *textsanitizerType) parseConfig() error {
	if cfg, ok := textsanitizer.Config.(map[string]interface{}); ok {
		if cutset, ok := cfg["cutset"].(string); ok {
			textsanitizer.cutset = cutset
		}
		if toLower, ok := cfg["toLower"].(bool); ok {
			textsanitizer.toLower = toLower
		}
		if toUpper, ok := cfg["toUpper"].(bool); ok {
			textsanitizer.toUpper = toUpper
		}
	} else {
		return fmt.Errorf("config needs to be an object with 'cutset':string, 'toLower':bool, 'toUpper':bool")
	}
	return nil
}

func (textsanitizer *textsanitizerType) sanitize(input string) string {
	if textsanitizer.toLower {
		input = strings.ToLower(input)
	}
	if textsanitizer.toUpper {
		input = strings.ToUpper(input)
	}
	for _, c := range textsanitizer.cutset {
		input = strings.Replace(input, string(c), "", -1)
	}
	return input
}

func (textsanitizer *textsanitizerType) backend() {
	for msg := range textsanitizer.data {
		if str, ok := msg.data.(string); ok {
			res := textsanitizer.sanitize(str)
			textsanitizer.Manager.Emit(textsanitizer.ID, "out", res)
		}
	}
}

func (textsanitizer *textsanitizerType) Process(in string, data interface{}) {
	textsanitizer.data <- textsanitizerMessage{in, data}
}

func (textsanitizer *textsanitizerType) Stop() {
	close(textsanitizer.data)
}

func init() {
	registry.Register("github.com/trusch/horst/processors/textsanitizer", func(id string, config interface{}, mgr horst.ProcessorManager) (horst.Processor, error) {
		textsanitizer := &textsanitizerType{data: make(chan textsanitizerMessage, 32)}
		textsanitizer.InitBase(id, config, mgr)
		err := textsanitizer.parseConfig()
		if err != nil {
			return nil, err
		}
		go textsanitizer.backend()
		return textsanitizer, nil
	})
}

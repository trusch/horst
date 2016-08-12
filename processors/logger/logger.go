package logger

import (
	"encoding/json"
	"fmt"

	"github.com/trusch/horst"
	"github.com/trusch/horst/processors/base"
	"github.com/trusch/horst/registry"
)

type loggerType struct {
	base.Base
	data chan loggerMessage
}

type loggerMessage struct {
	input string
	data  interface{}
}

func (logger *loggerType) backend() {
	for msg := range logger.data {
		switch msg.data.(type) {
		case map[string]interface{}:
			data, _ := json.MarshalIndent(msg.data, "", "  ")
			fmt.Printf("%v:%v> %v\n", logger.ID, msg.input, string(data))
		default:
			fmt.Printf("%v:%v> %v\n", logger.ID, msg.input, msg.data)
		}
	}
}

func (logger *loggerType) Process(in string, data interface{}) {
	logger.data <- loggerMessage{in, data}
}

func (logger *loggerType) Stop() {
	close(logger.data)
}

func init() {
	registry.Register("logger", func(id string, config interface{}, mgr horst.ProcessorManager) (horst.Processor, error) {
		logger := &loggerType{data: make(chan loggerMessage, 32)}
		logger.InitBase(id, config, mgr)
		go logger.backend()
		return logger, nil
	})
}

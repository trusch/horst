package logger

import (
	"fmt"

	"github.com/trusch/horst"
	"github.com/trusch/horst/registry"
)

type loggerType struct {
	id   string
	data chan loggerMessage
}

type loggerMessage struct {
	input string
	data  interface{}
}

func (logger *loggerType) backend() {
	for msg := range logger.data {
		fmt.Printf("%v:%v> %v\n", logger.id, msg.input, msg.data)
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
		logger := &loggerType{id, make(chan loggerMessage, 32)}
		go logger.backend()
		return logger, nil
	})
}

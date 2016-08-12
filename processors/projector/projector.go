package projector

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/trusch/horst"
	"github.com/trusch/horst/processors/base"
	"github.com/trusch/horst/registry"
	"github.com/trusch/jsonq"
)

type projectorType struct {
	base.Base
	data chan projectorMessage
}

type projectorMessage struct {
	input string
	data  interface{}
}

func (projector *projectorType) backend() {
	for msg := range projector.data {
		if obj, ok := msg.data.(map[string]interface{}); ok {
			res := projector.copyProjectionConfig()
			res = projector.searchAndReplaceKeys(res, obj).(map[string]interface{})
			projector.Manager.Emit(projector.ID, "out", res)
		}
	}
}

func (projector *projectorType) Process(in string, data interface{}) {
	projector.data <- projectorMessage{in, data}
}

func (projector *projectorType) Stop() {
	close(projector.data)
}

func (projector *projectorType) copyProjectionConfig() map[string]interface{} {
	result := make(map[string]interface{})
	var mod bytes.Buffer
	enc := json.NewEncoder(&mod)
	dec := json.NewDecoder(&mod)
	enc.Encode(projector.Config)
	dec.Decode(&result)
	return result
}

func (projector *projectorType) searchAndReplaceKeys(arg interface{}, inputDoc map[string]interface{}) interface{} {
	if argStr, ok := arg.(string); ok && argStr[0] == '@' {
		parts := strings.Split(argStr[1:], ".")
		jq := jsonq.NewQuery(inputDoc)
		val, _ := jq.Get(parts...)
		return val
	}
	if argMap, ok := arg.(map[string]interface{}); ok {
		for k, v := range argMap {
			argMap[k] = projector.searchAndReplaceKeys(v, inputDoc)
		}
	}
	return arg
}

func init() {
	registry.Register("projector", func(id string, config interface{}, mgr horst.ProcessorManager) (horst.Processor, error) {
		projector := &projectorType{data: make(chan projectorMessage, 32)}
		projector.InitBase(id, config, mgr)
		go projector.backend()
		return projector, nil
	})
}

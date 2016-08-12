package registry

import (
	"fmt"
	"testing"

	"github.com/trusch/horst"
)

type TestProcessor struct{}

func (p *TestProcessor) Init(id string, config interface{}, mgr horst.ProcessorManager) {}
func (p *TestProcessor) Process(in string, data interface{}) {
	fmt.Println(in, data)
}

func TestRegistry(t *testing.T) {
	Register("test", func(id string, config interface{}, mgr horst.ProcessorManager) (horst.Processor, error) {
		return new(TestProcessor), nil
	})
	proc, _ := Construct("test", "t1", "bla", nil)
	proc.Process("in1", "hello world")
}

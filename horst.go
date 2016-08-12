package horst

type ProcessorManager interface {
	AddProcessor(processorID string, processor Processor)
	DelProcessor(processorID string)
	Emit(fromProcessor, fromProcessorOutput string, data interface{})
	Process(processorID, inputID string, data interface{})
}

type Processor interface {
	Process(inputID string, data interface{})
	Stop()
}

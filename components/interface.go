package components

// Component is the interface for all go based components
// It enables us to write a common config and transport layer
type Component interface {
	// HandleConfigUpdate gets called when new config for this component is available
	HandleConfigUpdate(config map[string]interface{}) error
	// HandleOutputUpdate gets called when new output endpoints are available
	HandleOutputUpdate(outputs map[string]string) error
	// Process gets called when a new event for a specific input should be processed
	Process(inputID string, event interface{}) error
}

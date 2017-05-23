package device

import (
	"encoding/json"
)

type Descriptor struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Type        string `json:"type"`
}

// FlogoAppDescriptor is the descriptor for a Flogo application
type FlogoDeviceDescriptor struct {
	*Descriptor

	DeviceType  string `json:"device_type"`
	Settings map[string]string `json:"settings"`

	Actions  []*ActionConfig  `json:"actions"`
	Triggers []*TriggerConfig `json:"triggers"`
}


// ParseDeviceDescriptor parse the device descriptor
func ParseDeviceDescriptor(deviceJson string) (*FlogoDeviceDescriptor, error) {
	descriptor := &FlogoDeviceDescriptor{}

	err := json.Unmarshal([]byte(deviceJson), descriptor)

	if err != nil {
		return nil, err
	}

	return descriptor, nil
}


// TriggerConfig is the configuration for a Trigger
type TriggerConfig struct {
	Id       string                 `json:"id"`
	Ref      string                 `json:"ref"`
	Settings map[string]interface{} `json:"settings"`
	Handlers []*HandlerConfig       `json:"handlers"`
}

// HandlerConfig is the configuration for the Trigger Handler
type HandlerConfig struct {
	ActionId string                 `json:"actionId"`
	Settings map[string]interface{} `json:"settings"`
}

// Config is the configuration for the Action
type ActionConfig struct {
	Id   string          `json:"id"`
	Ref  string          `json:"ref"`
	Data DeviceActivity `json:"data"`
}

//todo hardcoded for now, should be generated from action-ref
type DeviceActivity struct {
	Id   string          `json:"id"`
	Ref  string          `json:"ref"`
	Settings map[string]interface{} `json:"settings"`
}


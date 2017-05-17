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

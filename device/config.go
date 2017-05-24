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
	MqttEnabled bool   `json:"mqtt_enabled"`
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
	Id       string            `json:"id"`
	Ref      string            `json:"ref"`
	ActionId string            `json:"actionId"`
	Settings map[string]string `json:"settings"`
}

func (tc *TriggerConfig) GetSetting(key string) string {
	return tc.Settings[key]
}

// Config is the configuration for the Action
type ActionConfig struct {
	Id   string          `json:"id"`
	Ref  string          `json:"ref"`
	Data DeviceActivity  `json:"data"`
}

//todo hardcoded for now, should be generated from action-ref
type DeviceActivity struct {
	UseTriggerVal bool          `json:"useTriggerVal"`
	Activity   *ActivityConfig  `json:"activity"`
}

type ActivityConfig struct {
	Id   string                `json:"id"`
	Ref  string                `json:"ref"`
	Settings map[string]string `json:"settings"`
}

func (ac *ActivityConfig) GetSetting(key string) string {
	return ac.Settings[key]
}

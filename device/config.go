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

	Device   *DeviceDetails    `json:"device"`
	Actions  []*ActionConfig   `json:"actions"`
	Triggers []*TriggerConfig  `json:"triggers"`
}

// FlogoAppDescriptor is the descriptor for a Flogo application
type DeviceDetails struct {
	Profile     string `json:"profile"`
	MqttEnabled bool   `json:"mqtt_enabled"`
	Settings    map[string]string `json:"settings"`

	Actions  []*ActionConfig  `json:"actions"`
	Triggers []*TriggerConfig `json:"triggers"`
}

type ActivityDescriptor struct {
	*Descriptor

	Ref           string      `json:"ref"`
	Libs          []*Lib      `json:"libs"`
	Settings      []*Setting  `json:"settings"`
	DeviceSupport []*DeviceSupportDetails `json:"device_support"`
}

type TriggerDescriptor struct {
	*Descriptor

	Ref           string      `json:"ref"`
	Libs          []*Lib      `json:"libs"`
	Settings      []*Setting  `json:"settings"`
	Outputs       []*Setting  `json:"outputs"`
	DeviceSupport []*DeviceSupportDetails `json:"device_support"`
}

type Setting struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Lib struct {
	Name    string `json:"name"`
	LibType string `json:"type"`
	Ref     string `json:"ref"`
}

type DeviceSupportDetails struct {
	Framework    string `json:"framework"`
	TemplateFile string `json:"template"`
}

type DeviceProfile struct {
	*Descriptor

	Board        string `json:"board"`
	Platform     string `json:"platform"`
	PlatformWifi string `json:"platform_wifi"`
}

type DevicePlatform struct {
	*Descriptor

	Framework    string `json:"arduino"`
	MainTemplate string `json:"main_template"`
	WifiDetails  []*PlatformFeature `json:"wifi"`
	MqttDetails  *PlatformFeature `json:"mqtt"`
}

type PlatformFeature struct {
	Name     string `json:"name"`
	Template string `json:"template"`
	Header   string `json:"header"`
	Libs     []*Lib `json:"libs"`
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
	Data json.RawMessage `json:"data"`
}

//todo consolidate parsing functions

// ParseActivityDescriptor parse the device activity descriptor
func ParseActivityDescriptor(contribJson string) (contrib *ActivityDescriptor, err error) {
	descriptor := &ActivityDescriptor{}

	err = json.Unmarshal([]byte(contribJson), descriptor)

	if err != nil {
		return nil, err
	}

	return descriptor, nil
}

// ParseTriggerDescriptor parse the device trigger descriptor
func ParseTriggerDescriptor(contribJson string) (contrib *TriggerDescriptor, err error) {
	descriptor := &TriggerDescriptor{}

	err = json.Unmarshal([]byte(contribJson), descriptor)

	if err != nil {
		return nil, err
	}

	return descriptor, nil
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

// ParseDeviceProfile parse the device profile
func ParseDeviceProfile(profileJson string) (*DeviceProfile, error) {
	profile := &DeviceProfile{}

	err := json.Unmarshal([]byte(profileJson), profile)

	if err != nil {
		return nil, err
	}

	return profile, nil
}

// ParseDeviceProfile parse the device platform
func ParseDevicePlatform(platformJson string) (*DevicePlatform, error) {
	profile := &DevicePlatform{}

	err := json.Unmarshal([]byte(platformJson), profile)

	if err != nil {
		return nil, err
	}

	return profile, nil
}

package device

import (
	"os"
	"fmt"
	"encoding/json"
)

type DevicesConfig struct {
	Devices []*DeviceConfig  `json:"devices"`
	Libs    []*LibConfig     `json:"libs"`
}

type DeviceConfig struct {
	Board    string  `json:"board"`
	Template string  `json:"template"`
	Source   string  `json:"source"`
}

type LibConfig struct {
	Name string  `json:"name"`
	ID   int     `json:"id"`
}

func LoadDevicesConfig(dir string) *DevicesConfig {

	devicesConfigFile, err := os.Open(dir + "/devices.json")

	if err != nil {
		fmt.Fprint(os.Stderr, "Error: devices.json not found.\n\n")
		os.Exit(2)
	}

	devicesConfig := &DevicesConfig{}
	jsonParser := json.NewDecoder(devicesConfigFile)

	if err = jsonParser.Decode(devicesConfig); err != nil {
		fmt.Fprint(os.Stderr, "Error: Unable to parse devices.json, file may be corrupted.\n\n")
		os.Exit(2)
	}

	devicesConfigFile.Close()

	return devicesConfig;
}

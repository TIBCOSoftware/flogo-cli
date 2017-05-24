package device

import (
	"sync"
)

var (
	devicesMu sync.Mutex
	devices   = make(map[string]*DeviceDetails)
)


type DeviceDetails struct {
	Type     string
	Board    string
	MainFile string
	MqttFiles map[string]string
	Libs map[string]int

	Files map[string]string
}

func Register(dd *DeviceDetails) {

	devicesMu.Lock()
	defer devicesMu.Unlock()

	if dd == nil {
		panic("device.Register: device details is nil")
	}

	if _, dup := devices[dd.Type]; dup {
		panic("device.Register: device already registered - " + dd.Type)
	}

	devices[dd.Type] = dd
}

func Devices() []*DeviceDetails {

	devicesMu.Lock()
	var curDevices = devices
	devicesMu.Unlock()


	list := make([]*DeviceDetails, 0, len(curDevices))

	for _, value := range curDevices {
		list = append(list, value)
	}

	return list
}

func GetDevice(deviceType string) *DeviceDetails {
	return devices[deviceType]
}
package data

import "strings"

type Devices struct {
	Devices map[string]Device `json:"devices"`
}

func NewDevices() Devices {
	return Devices{
		Devices: make(map[string]Device),
	}
}

func (devices *Devices) AddDevice(device Device) {
	//if _, ok := devices.Devices[device.name]; !ok {
	devices.Devices[device.Name] = device
	//} else {

	//}
}

func (devices *Devices) ShowDeviceNames() string {
	sb := strings.Builder{}
	sb.WriteString("Device Names: \n")
	for k, _ := range devices.Devices {
		sb.WriteString(k + "\n")
	}
	return sb.String()
}

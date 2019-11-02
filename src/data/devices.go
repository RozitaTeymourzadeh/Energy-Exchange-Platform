package data

//
//import (
//	"github.com/edgexfoundry/go-mod-core-contracts/models"
//	"strings"
//)
//
//type Devices struct {
//	//Devices map[string]Device `json:"devices"`
//	Devices map[string]models.Device `json:"devices"`
//}
//
//func NewDevices() Devices {
//	return Devices{
//		Devices: make(map[string]models.Device),
//	}
//}
//
//func (devices *Devices) AddDevice(device models.Device) {
//	//if _, ok := devices.Devices[device.name]; !ok {
//	devices.Devices[device.Name] = device
//	//} else {
//	//}
//}
//
//func (devices *Devices) ShowDeviceNames() string {
//	sb := strings.Builder{}
//	sb.WriteString("Device Names: \n")
//	for k, _ := range devices.Devices {
//		sb.WriteString(k + "\n")
//	}
//	return sb.String()
//}

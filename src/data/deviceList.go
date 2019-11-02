package data

import (
	"encoding/json"
	"fmt"
	"github.com/edgexfoundry/go-mod-core-contracts/models"
	"strings"
)

//
type Device struct {
	PeerId         string           `json:"-"` // peer Id of registered peer
	Id             string           `json:"id"`
	Name           string           `json:"name"`
	AdminState     string           `json:"adminState"`
	OperatingState string           `json:"operatingState"`
	LastConnected  int64            `json:"lastConnected"`
	LastReported   int64            `json:"lastReported"`
	Labels         []string         `json:"labels"`
	Location       string           `json:"location"`
	Commands       []models.Command `json:"commands"`
}

//
type DeviceList struct {
	Devices []Device `json:"devices"`
}

//
func NewDeviceList() DeviceList {
	return DeviceList{
		Devices: make([]Device, 0),
	}
}

//
func DeviceListFromJson(jsonBytes []byte) DeviceList {
	dl := NewDeviceList()
	err := json.Unmarshal(jsonBytes, &dl.Devices)
	if err != nil {
		fmt.Println("Error in getting device list from JSON", err)
	}
	return dl
}

//func (dl *DeviceList) GetDeviceList() []string {
//	for _, val := range dl.Devices {
//		dl = append(dl, k)
//	}
//	return dl
//
//}

//
func (dl *DeviceList) ShowDeviceInList() string {
	sb := strings.Builder{}
	sb.WriteString("Device List: \n")
	for _, val := range dl.Devices {

		sb.WriteString(val.Name + " : " + val.Id + "\n")
	}
	return sb.String()
}

//
//func (dl *DeviceList) AddAllToDevices(devices *Devices) *Devices {
//	for _, val := range dl.Devices {
//		devices.AddDevice(val)
//	}
//	return devices
//}

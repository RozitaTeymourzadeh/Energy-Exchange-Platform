package driver

import (
	"encoding/json"
	"fmt"
	"github.com/edgexfoundry/go-mod-core-contracts/models"
	"strings"
)

//
type Device struct {
	PeerId         string           `json:"-"` // peer Id of registered peer ip and port
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

type DeviceMap struct {
	Devices map[string]Device `json:"devices"`
}

//
func NewDeviceMap() DeviceMap {
	return DeviceMap{
		Devices: make(map[string]Device, 0),
	}
}

func (dm *DeviceMap) DeviceMapToList() []Device {
	devices := make([]Device, 0)
	if len(dm.Devices) > 0 {
		for _, d := range dm.Devices {
			devices = append(devices, d)
		}

	}
	return devices
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

func DeviceMapFromJson(jsonBytes []byte) DeviceMap {
	dm := NewDeviceMap()
	err := json.Unmarshal(jsonBytes, &dm.Devices)
	if err != nil {
		fmt.Println("Error in getting device list from JSON", err)
	}
	return dm
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

		sb.WriteString("@" + val.PeerId + val.Name + " : " + val.Id + "\n")
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

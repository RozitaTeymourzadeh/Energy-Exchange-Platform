package data

import (
	"encoding/json"
	"fmt"
	"github.com/edgexfoundry/go-mod-core-contracts/models"
	"strings"
)

type Device struct {
	Id             string           `json:"id"`
	Name           string           `json:"name"`
	AdminState     string           `json:"adminState"`
	OperatingState string           `json:"operatingState"`
	LastConnected  string           `json:"lastConnected"`
	LastReported   string           `json:"lastReported"`
	Labels         []string         `json:"labels"`
	Location       string           `json:"location"`
	Commands       []models.Command `json:"commands"`
}

type DeviceList struct {
	Devices []Device
}

func NewDeviceList() DeviceList {
	return DeviceList{
		Devices: make([]Device, 0),
	}
}

func DeviceListFromJson(jsonBytes []byte) DeviceList {
	dl := NewDeviceList()
	err := json.Unmarshal(jsonBytes, &dl.Devices)
	if err != nil {
		fmt.Println("Error in getting device list from JSON")
	}
	return dl
}

func (dl *DeviceList) ShowDeviceInList() string {
	sb := strings.Builder{}
	sb.WriteString("Device List: \n")
	for _, val := range dl.Devices {
		sb.WriteString(val.Name + " : " + val.Id + "\n")
	}
	return sb.String()
}

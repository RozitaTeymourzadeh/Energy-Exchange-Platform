package data

import (
	"encoding/json"
	"errors"
	"fmt"
)

// map of key = device and value = map of key = Readings.name & value = *CoreDataEvents
type DeviceEvents struct {
	//         device     name
	Events map[string]DeviceResourceNameEvents `json:"events"`
}

func NewDeviceEvents() DeviceEvents {
	return DeviceEvents{
		Events: make(map[string]DeviceResourceNameEvents),
	}
}

func (de *DeviceEvents) AddToDeviceEvents(cde CoreDataEvent) {
	deviceName := cde.Device
	drne := de.Events[deviceName]
	if drne.resourceEvents != nil {
		fmt.Println("drnes is NOT nil !!")
		drne.AddToDeviceResourceNameEvents(cde)
	} else {
		fmt.Println("drnes is nil !!")
		/*de.Events[deviceName]*/ drne = NewDeviceResourceNameEvents()
		drne.AddToDeviceResourceNameEvents(cde)
		de.Events[deviceName] = drne
	}

	//fmt.Println(de.Events[deviceName]) // todo : remember long print

}

func (de *DeviceEvents) GetLatestDeviceResourceNameEventForDevice(deviceName string, deviceResourceName string) (CoreDataEvent, error) {
	drnes := de.Events[deviceName]
	if drnes.resourceEvents != nil {
		return drnes.GetLatestDeviceResourceNameEvent(deviceResourceName)
	}
	return NewCoreDataEvent(), errors.New("no DeviceResourceNameEvents associated with deviceName exists in DeviceEvents")
}

func (des *DeviceEvents) DeviceEventsToJson() []byte { //todo : correct this method
	desJson, err := json.Marshal(des)
	if err != nil {
		fmt.Println("Error in ToJson of reading.go")
	}
	desJsonStr := (string(desJson))
	fmt.Println(desJsonStr)
	return desJson
}

///////////////

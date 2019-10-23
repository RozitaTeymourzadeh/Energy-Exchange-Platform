package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
)

// map of key = Readings.name & value = *CoreDataEvents
type DeviceResourceNameEvents struct {
	//                (deviceResourceName to CoreDataEvents)
	//         deviceResourceName     CoreDataEvent
	resourceEvents map[string]CoreDataEvents `json:"resourceEvents"`
	mux            sync.Mutex
}

func NewDeviceResourceNameEvents() DeviceResourceNameEvents {
	//cdes := NewCoreDataEvents()
	return DeviceResourceNameEvents{
		resourceEvents: make(map[string]CoreDataEvents),
		mux:            sync.Mutex{},
	}
}

func (drnes *DeviceResourceNameEvents) AddToDeviceResourceNameEvents(cde CoreDataEvent) {
	drnes.mux.Lock()
	defer drnes.mux.Unlock()

	fmt.Println("In AddToDeviceResourceNameEvents")
	if drnes.resourceEvents == nil {
		drnes.resourceEvents = make(map[string]CoreDataEvents)
	}

	deviceResourceReadings := cde.Readings
	for _, deviceResourceReading := range deviceResourceReadings {
		fmt.Println("Adding - deviceResourceReading : " + deviceResourceReading.Device + ", " + deviceResourceReading.Name + ", " + deviceResourceReading.Value)
		resourceEventsValue := drnes.resourceEvents[deviceResourceReading.Name]
		if len(resourceEventsValue.DataEvents) == 0 {
			fmt.Println("resourceEventsValue.DataEvents is zero")
		} else {
			resourceEventsValue.DataEvents = append(resourceEventsValue.DataEvents, cde)

			drnes.resourceEvents[deviceResourceReading.Name] = resourceEventsValue
		}

	}

}

func (drnes *DeviceResourceNameEvents) GetLatestDeviceResourceNameEvent(deviceResourceName string) (CoreDataEvent, error) {
	drnes.mux.Lock()
	defer drnes.mux.Unlock()

	deviceResourceNameEventsValue := drnes.resourceEvents[deviceResourceName]
	if len(deviceResourceNameEventsValue.DataEvents) > 0 {
		return deviceResourceNameEventsValue.DataEvents[0], nil
	}
	return NewCoreDataEvent(), errors.New("no event associated with deviceResourceName exist in DeviceResourceNameEvents")
}

func (drne DeviceResourceNameEvents) DeviceResourceNameEventsToJson() []byte {
	readingJson, err := json.Marshal(&drne)
	if err != nil {
		fmt.Println("Error in ToJson of DeviceResourceNameEvents")
	}
	return readingJson
}

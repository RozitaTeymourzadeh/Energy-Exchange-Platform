package data

import (
	"encoding/json"
	"fmt"
)

type CoreDataEvents struct {
	DataEvents []CoreDataEvent `json:"dataEvents"`
}

func NewCoreDataEvents() CoreDataEvents {
	return CoreDataEvents{
		DataEvents: make([]CoreDataEvent, 1),
	}
}

func CoreDataEventsToJson(cdes *CoreDataEvents) []byte {
	cdesJson, err := json.Marshal(cdes.DataEvents)
	if err != nil {
		fmt.Println("Error in ToJson of coreDataEvents.go")
	}
	return cdesJson
}

func CoreDataEventsFromJson(jsonbytes []byte) *CoreDataEvents {
	cdes := CoreDataEvents{
		DataEvents: []CoreDataEvent{},
	}
	fmt.Println("Recv in CoreDataEventsFromJson : \n", string(jsonbytes))
	err := json.Unmarshal(jsonbytes, &cdes.DataEvents)
	if err != nil {
		fmt.Println("Error in FromJson of coreDataEvents.go")
	}
	return &cdes

}

package data

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"
)

type CoreDataEvents struct {
	DataEvents []CoreDataEvent `json:"dataEvents"`
}

func NewCoreDataEvents() CoreDataEvents {
	return CoreDataEvents{
		DataEvents: make([]CoreDataEvent, 0),
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

type ByCreated []CoreDataEvent

func (a ByCreated) Len() int           { return len(a) }
func (a ByCreated) Less(i, j int) bool { return a[i].Created > a[j].Created }
func (a ByCreated) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func (cdes CoreDataEvents) Sort() CoreDataEvents {
	sort.Sort(ByCreated(cdes.DataEvents))
	unixTimeUTC := time.Unix(cdes.DataEvents[0].Created, 0)
	// todo : delete duplicate
	// todo : ================
	fmt.Println("Last created event: ", unixTimeUTC)
	return cdes
}

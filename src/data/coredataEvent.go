package data

import (
	"encoding/json"
	"fmt"
)

//[
//	{
//		"id": "fd77d573-dcf7-4946-bc36-9e45c93ed5a9",
//		"device": "Supply-Device-01",
//		"created": 1571587181556,
//		"modified": 1571587181556,
//		"origin": 1571587181527,
//		"readings": [
//			{
//				"id": "7e747dcf-dce1-4105-9bf5-b99721bddc1f",
//				"created": 1571587181530,
//				"origin": 1571587181527195000,
//				"modified": 1571587181530,
//				"device": "Supply-Device-01",
//				"name": "randomsuppliernumber",
//				"value": "81"
//			}
//		]
//	},
//]

type CoreDataEvent struct {
	Id       string    `json:"id"`
	Device   string    `json:"device"`
	Created  int64     `json:"created"`
	Modified int64     `json:"modified"`
	Origin   int64     `json:"origin"`
	Readings []Reading `json:"readings"`
}

func NewCoreDataEvent() CoreDataEvent {
	return CoreDataEvent{
		Id:       "",
		Device:   "",
		Created:  0,
		Modified: 0,
		Origin:   0,
		Readings: make([]Reading, 1),
	}
}

func (cde *CoreDataEvent) CoreDataEventToJson() []byte {
	cdeJson, err := json.Marshal(cde)
	if err != nil {
		fmt.Println("Error in ToJson of reading.go")
	}
	return cdeJson
}

func (cde *CoreDataEvent) CoreDataEventFromJson(jsonbytes []byte) {
	err := json.Unmarshal(jsonbytes, cde)
	if err != nil {
		fmt.Println("Error in FromJson of reading.go")
	}
}

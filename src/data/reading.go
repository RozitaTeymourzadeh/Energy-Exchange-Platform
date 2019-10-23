package data

import (
	"encoding/json"
	"fmt"
)

//			{
//				"id": "7e747dcf-dce1-4105-9bf5-b99721bddc1f",
//				"created": 1571587181530,
//				"origin": 1571587181527195000,
//				"modified": 1571587181530,
//				"device": "Supply-Device-01",
//				"name": "randomsuppliernumber",
//				"value": "81"
//			}

type Reading struct {
	Id       string `json:"id"`
	Created  int64  `json:"created"`
	Origin   int64  `json:"origin"`
	Modified int64  `json:"modified"`
	Device   string `json:"device"`
	Name     string `json:"name"`
	Value    string `json:"value"`
}

func NewReading() Reading {
	return Reading{
		Id: "",
		//Created:  0,
		//Origin:   0,
		//Modified: 0,
		Device: "",
		Name:   "",
		Value:  "",
	}
}

func (r *Reading) ReadingToJson() []byte {
	readingJson, err := json.Marshal(r)
	if err != nil {
		fmt.Println("Error in ToJson of reading.go")
	}
	return readingJson
}

func (r *Reading) ReadingFromJson(jsonbytes []byte) {
	err := json.Unmarshal(jsonbytes, r)
	if err != nil {
		fmt.Println("Error in FromJson of reading.go")
	}
}

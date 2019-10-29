package data

import "encoding/json"

type RegistryInfo struct {
	IpAdd string `json:"ipAdd"`
	Port  string `json:"port"`
}

//func NewRegistryInfo(string ipAdd) {
//
//}

func (ri *RegistryInfo) RegistryToJSON() []byte {
	jsonBytes, _ := json.Marshal(ri)
	return jsonBytes
}

func RegistryFromJSON(jsonBytes []byte) RegistryInfo {
	ri := RegistryInfo{}
	_ = json.Unmarshal(jsonBytes, ri)
	return ri
}

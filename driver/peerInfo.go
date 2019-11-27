package driver

import "encoding/json"

type PeerInfo struct {
	IpAdd string `json:"ipAdd"`
	Port  string `json:"port"`
}

//func NewRegistryInfo(string ipAdd) {
//
//}

func (ri *PeerInfo) PeerInfoToJSON() []byte {
	jsonBytes, _ := json.Marshal(ri)
	return jsonBytes
}

func PeerInfoFromJSON(jsonBytes []byte) PeerInfo {
	ri := PeerInfo{}
	_ = json.Unmarshal(jsonBytes, &ri)
	return ri
}

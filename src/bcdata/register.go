package bcdata

import "encoding/json"

/* RegisterDate
*
* RegisterData is only used during registration.
* You may return PeerMapJson to a new node by this data structure.
*
 */
type RegisterData struct {
	AssignedId  int32  `json:"assignedId"`
	PeerMapJson string `json:"peerMapJson"`
}

/* NewRegisterData
*
* To initial RegisterData
*
 */
func NewRegisterData(id int32, peerMapJson string) RegisterData {
	r := RegisterData{AssignedId: id, PeerMapJson: peerMapJson}
	return r
}

/* EncodeToJson
*
* To Encode RegisterData into Json format
*
 */
func (data *RegisterData) EncodeToJson() (string, error) {
	jsonBytes, err := json.Marshal(data)
	return string(jsonBytes), err
}

/* DecodeToJson
*
* To Decode RegisterData into Json format
*
 */
func (data *RegisterData) DecodeFromJSON(jsonString string) error {
	return json.Unmarshal([]byte(jsonString), data)
}

package data

import (
	"MerklePatriciaTree/p5/Blockchain_Application_P5/p4"
	"crypto/rsa"
	"encoding/json"

	//"fmt"
)

/* HeartBeatData
*
* HeartBeatData is used when two nodes are sending or receiving HeartBeats
*
 */
type HeartBeatData struct {
	IfNewBlock  		bool   			`json:"ifNewBlock"`
	Id          		int32  			`json:"id"`
	BlockJson   		string 			`json:"blockJson"`
	PeerMapJson 		string 			`json:"peerMapJson"`
	Addr        		string 			`json:"addr"`
	Hops        		int32			`json:"hops"`
	PeerPublicKey 		*rsa.PublicKey `json:"peerPublicKey"`
	IfValidTransaction 	bool			`json:"ifValidTransaction"`
	TransactionInfoJson string			`json:"transactionInfoJson"`
}


/* NewHeartBeatData()
*
* NewHeartBeatData() is a normal initial function which creates an instance
*
 */
func NewHeartBeatData(ifNewBlock bool, id int32, blockJson string, peerMapJson string, addr string, peerPublicKey *rsa.PublicKey, ifValidTransaction bool, transactionInfoJson string) HeartBeatData {
	return HeartBeatData{
		IfNewBlock:  ifNewBlock,
		Id:          id,
		BlockJson:   blockJson,
		PeerMapJson: peerMapJson,
		Addr:        addr,
		Hops:        2,
		PeerPublicKey: peerPublicKey,
		IfValidTransaction: ifValidTransaction,
		TransactionInfoJson: transactionInfoJson,
	}
}

/* PrepareHeartBeatData()
*
* PrepareHeartBeatData() is used when you want to send a HeartBeat to other peers.
* PrepareHeartBeatData would first create a new instance of HeartBeatData, then decide
* whether or not you will create a new block and send the new block to other peers.
*
 */
func PrepareHeartBeatData(sbc *SyncBlockChain, selfId int32, peerMapJson string, addr string, generateNewBlock bool,nonce string , mpt p4.MerklePatriciaTrie,peerPublicKey *rsa.PublicKey, ifValidTransaction bool, transactionInfoJson string) HeartBeatData {
	newHeartBeatData := NewHeartBeatData(false, selfId, "", peerMapJson, addr,peerPublicKey,ifValidTransaction ,transactionInfoJson)
	if generateNewBlock  {
		newBlock := sbc.GenBlock(mpt, nonce)
		blockJson, _ := newBlock.EncodeToJSON()
		newHeartBeatData = NewHeartBeatData(true, selfId, blockJson, peerMapJson, addr,peerPublicKey,ifValidTransaction ,transactionInfoJson)
	}else{
		newHeartBeatData = NewHeartBeatData(false, selfId, "", peerMapJson, addr,peerPublicKey,false ,transactionInfoJson)
	}
	return newHeartBeatData
}


/* EncodeToJson()
*
* To Encode HeartBeatData from json format
*
 */
func (data *HeartBeatData) EncodeToJson() (string, error) {
	jsonBytes, error := json.Marshal(data)
	return string(jsonBytes), error
}

/* DecodeFromJson()
*
* To Decode HeartBeatData from json format
*
 */
func (data *HeartBeatData) DecodeFromJson(jsonString string) error {
	return json.Unmarshal([]byte(jsonString), data)
}
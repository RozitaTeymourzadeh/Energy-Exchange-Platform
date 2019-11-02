package data

import (
	"fmt"
	"sync"
)

type nodeId struct {
	Address   string
	Port      string
	Separator string
	peers     map[string]Peer
}

var instanceNodeId *nodeId
var onceNodeId sync.Once

func SetNodeId(address string, port string /*, connectingAddress string*/) *nodeId {
	onceNodeId.Do(func() {
		instanceNodeId = &nodeId{
			Address:   address,
			Port:      port,
			Separator: ":",
			peers:     make(map[string]Peer),
		}
	})
	return instanceNodeId
}

func GetNodeId() *nodeId {
	return instanceNodeId
}

func (nid *nodeId) GetPeers() map[string]Peer {
	fmt.Println(nid)
	return nid.peers
}

func (nid *nodeId) AddPeer(peerObj Peer) {
	nid.peers[peerObj.ID] = peerObj
	fmt.Println("Size of peers : ", len(nid.peers))
}

//// constructor for NodeId
//func NewtNodeId(address string, port int) NodeId {
//	nid := NodeId{}
//	nid.Separator = ":"
//	nid.Address = address
//	nid.Port = string(port)
//	return nid
//}

//// address and port
//func (nid *NodeId)SplitAddressAndPort() []string {
//	return strings.Split(addressAndPort, separator)
//}

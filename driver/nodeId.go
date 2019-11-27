package driver

import (
	"fmt"
	"sync"
)

type nodeId struct {
	Address            string
	Port               string
	EdgeXAddress       string
	TaskManagerAddress string
	Separator          string
	Peers              []PeerInfo
}

var instanceNodeId *nodeId
var onceNodeId sync.Once

func SetNodeId(address string, port string, taskManagerAddress string, edgeXAddress string) *nodeId {
	onceNodeId.Do(func() {
		instanceNodeId = &nodeId{
			Address:            address,
			Port:               port,
			EdgeXAddress:       edgeXAddress,
			TaskManagerAddress: taskManagerAddress,
			Separator:          ":",
			Peers:              make([]PeerInfo, 0),
		}
	})
	fmt.Println(instanceNodeId)
	return instanceNodeId
}

func GetNodeId() *nodeId {
	return instanceNodeId
}

func (nid *nodeId) SetEdgeXAddress(edgeAddress string) {
	nid.EdgeXAddress = edgeAddress
}

func (nid *nodeId) GetPeers() []PeerInfo {
	for _, peer := range nid.Peers {
		fmt.Println("IP: " + peer.IpAdd)
	}
	fmt.Println(nid)
	return nid.Peers
}

func (nid *nodeId) AddPeer(rInfo PeerInfo) {
	nid.Peers = append(nid.Peers, rInfo)
	fmt.Println("Size of peers : ", len(nid.Peers))
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

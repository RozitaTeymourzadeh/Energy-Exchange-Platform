package data

import (
	"sync"
)

type nodeId struct {
	Address   string
	Port      string
	Separator string
}

var instanceNodeId *nodeId
var onceNodeId sync.Once

func SetNodeId(address string, port string) *nodeId {
	onceNodeId.Do(func() {
		instanceNodeId = &nodeId{
			Address:   address,
			Port:      port,
			Separator: ":",
		}
	})
	return instanceNodeId
}

func GetNodeId() *nodeId {
	return instanceNodeId
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

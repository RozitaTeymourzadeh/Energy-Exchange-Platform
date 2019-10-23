package data

import "strings"

type NodeId struct {
	Address   string
	Port      string
	Separator string
}

// constructor for NodeId
func NewtNodeId(address string, port int) NodeId {
	nid := NodeId{}
	nid.Separator = ":"
	nid.Address = address
	nid.Port = string(port)
	return nid
}

// address and port
func splitAddressAndPort(addressAndPort string, separator string) []string {
	addressAndPort = ""
	return strings.Split(addressAndPort, separator)
}

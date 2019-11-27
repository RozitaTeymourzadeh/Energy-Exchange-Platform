package bcdata

import (
	"container/ring"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"sync"
)

/* PeerList
*
* PeerList is used to store peers' addresses and IDs
*
 */
type PeerList struct {
	selfId  int32
	peerMap map[string]int32
	//peerPublicKeyMap map[*rsa.PublicKey]int32
	maxLength int32
	mux       sync.Mutex
}

/* Pair
*
* Pair addr/id
*
 */
type Pair struct {
	addr string
	id   int32
}

/* Pair
*
* A slice of Pairs that implements sort.
* Interface to sort by Value.
*
 */
type PairList []Pair

/* NewPeerList
*
* init peerMap and id, maxLength
* NewPeerList() is the initial function of PeerList structure
*
 */
func NewPeerList(id int32, maxLength int32) PeerList {
	peerList := PeerList{
		peerMap: make(map[string]int32),
		//peerPublicKeyMap: make(map[*rsa.PublicKey]int32),
		maxLength: maxLength}
	peerList.Register(id)
	return peerList
}

/* GetPeerMap()
*
* To return the peerMap
*
 */
func (peers *PeerList) GetPeerMap() map[string]int32 {
	peers.mux.Lock()
	defer peers.mux.Unlock()
	return peers.peerMap
}

/* GetPeerMap()
*
* To return the peerMap
*
 */
//func (peers *PeerList) GetPublicKeyMap() map[*rsa.PublicKey]int32{
//	peers.mux.Lock()
//	defer peers.mux.Unlock()
//	return peers.peerPublicKeyMap
//}

/* GetMaxLength()
*
* To return the peerMap
*
 */
func (peers *PeerList) GetMaxLength() int32 {
	return peers.maxLength
}

/* Add()
*
* To add ip and address into the peerMap
*
 */
func (peers *PeerList) Add(addr string, id int32) {
	peers.mux.Lock()
	peers.peerMap[addr] = id
	peers.mux.Unlock()
}

//func(peers *PeerList) AddPublicKey(publicKey *rsa.PublicKey, id int32) {
//	peers.mux.Lock()
//	peers.peerPublicKeyMap[publicKey] = id
//	peers.mux.Unlock()
//}

/* Delete()
*
* To delete ip and address from the peerMap
*
 */
func (peers *PeerList) Delete(addr string) {
	peers.mux.Lock()
	delete(peers.peerMap, addr)
	peers.mux.Unlock()
}

/* Rebalance()
*
* Rebalance func changes the PeerMap to contain take maxLength(32) closest peers (by Id)
* The PeerList can temporarily hold more than 32 nodes,
* but before sending HeartBeats, a node will first re-balance the PeerList by choosing the 32 closest peers.
* "Closest peers" is defined by this: Sort all peers' Id, insert SelfId, consider the list as a cycle,
* and choose 16 nodes at each side of SelfId. For example, if SelfId is 10, PeerList is [7, 8, 9, 15, 16], then the closest 4 nodes are [8, 9, 15, 16]. HeartBeat is sent to every peer nodes at "/heartbeat/receive".
*
* https://golang.org/pkg/container/ring
*
*
 */
func (peers *PeerList) Rebalance() {
	peers.mux.Lock()
	defer peers.mux.Unlock()
	if int32(len(peers.peerMap)) > peers.maxLength {
		peers.peerMap["selfAddr"] = peers.selfId //adding self id to peerMap
		sortedAddrIDList := sortMapByValue(peers.peerMap)
		sortedAddrIDListLength := len(sortedAddrIDList)
		peers.peerMap = peers.getBalancedPeerMap(sortedAddrIDListLength, sortedAddrIDList)
	}
}

/* Swap()
*
* To Swap element for sorting purpose
*
 */
func (pairs PairList) Swap(i, j int) {
	pairs[i], pairs[j] = pairs[j], pairs[i]
}

/* Len()
*
* To return pair Len
*
 */
func (pairs PairList) Len() int {
	return len(pairs)
}

/* Less()
*
* To conduct element comparison
*
 */
func (pairs PairList) Less(i, j int) bool {
	return pairs[i].id < pairs[j].id
}

/* sortMapByValue()
*
* To turn a map into a PairList, then sort and return it.
*
 */
func sortMapByValue(m map[string]int32) PairList {
	p := make(PairList, len(m))
	i := 0
	for k, v := range m {
		//fmt.Println("in sortMapByValue : k, v :", k, v)
		p[i] = Pair{
			addr: k,
			id:   int32(v),
		}
		i++
	}
	sort.Sort(p)
	return p
}

/* getBalancedPeerMap()
*
* To Balance PeerMap
*
 */
func (peers *PeerList) getBalancedPeerMap(sortedAddrIDListLength int, sortedAddrIDList PairList) map[string]int32 {

	r := ring.New(sortedAddrIDListLength) // new ring
	useRingPtr := r
	for i := 0; i < sortedAddrIDListLength; i++ {
		r.Value = sortedAddrIDList[i]
		if sortedAddrIDList[i].id == peers.selfId {
			useRingPtr = r
		}
		r = r.Next()
	}
	newPeerMap := make(map[string]int32)
	r = useRingPtr
	for i := 1; i <= int(peers.maxLength/2); i++ {
		r = r.Prev()
		pair := r.Value.(Pair)
		newPeerMap[pair.addr] = pair.id
	}
	r = useRingPtr
	for i := 1; i <= int(peers.maxLength/2); i++ {
		r = r.Next()
		pair := r.Value.(Pair)
		newPeerMap[pair.addr] = pair.id
	}
	return newPeerMap
}

/* Show()
*
* To shows all addresses and their corresponding IDs.
* For example, it returns "This is PeerMap: \n addr=127.0.0.1, id=1".
*
 */
func (peers *PeerList) Show() string {
	rs := ""
	peers.mux.Lock()
	defer peers.mux.Unlock()
	for addr, id := range peers.peerMap {
		rs += fmt.Sprintf("addr= %s ", addr)
		rs += fmt.Sprintf(", id= %d \n", id)
	}
	rs += "\n"
	//for publicKey, id := range peers.peerPublicKeyMap {
	//	rs += fmt.Sprintf("public Key= %s ", publicKey)
	//	rs += fmt.Sprintf(", id= %d \n", id)
	//}
	//rs += "\n"
	//rs = fmt.Sprintf("This is the PeerMap: %s\n", hex.EncodeToString(sum[:])) + rs
	rs = fmt.Sprintf("This is the PeerMap: \n") + rs
	fmt.Print(rs)
	return rs
}

/* Register()
*
* Register() is used to set ID.
* You can consider it as "SetId()".
*
 */
func (peers *PeerList) Register(id int32) {
	peers.mux.Lock()
	peers.selfId = id
	fmt.Printf("SelfId=%v\n", id)
	peers.mux.Unlock()
}

/* Copy()
*
* Copy func returns a copy of the peerMap
*
 */
func (peers *PeerList) Copy() map[string]int32 {

	peers.mux.Lock()
	copyOfPeerMap := make(map[string]int32)
	for k := range peers.peerMap {
		copyOfPeerMap[k] = peers.peerMap[k]
	}
	peers.mux.Unlock()
	return copyOfPeerMap
}

/* GetSelfId()
*
* Return peerList.SelfId
*
 */
func (peers *PeerList) GetSelfId() int32 {
	return peers.selfId
}

/* PeerMapToJson()
*
* The "PeerMapJson" in HeartBeatData is the JSON format of "PeerList.peerMap"
* It is the result of "PeerList.PeerMapToJSON()" function.
*
 */
func (peers *PeerList) PeerMapToJson() (string, error) {
	jsonBytes, err := json.Marshal(peers.peerMap)
	return string(jsonBytes), err
}

/* InjectPeerMapJson()
*
* To inject PeerMap to PeerMap List
*
 */
func (peers *PeerList) InjectPeerMapJson(peerMapJsonStr string, selfAddr string) {
	var newPeerMap map[string]int32
	err := json.Unmarshal([]byte(peerMapJsonStr), &newPeerMap)
	if err == nil {
		peers.mux.Lock()
		for k := range newPeerMap {
			if k != selfAddr {
				peers.peerMap[k] = newPeerMap[k]
			}
		}
		peers.mux.Unlock()
	}
}

/* EncodePeerMapToJSON()
*
* To encode PeerMap into Json
*
 */
func (peers *PeerList) EncodePeerMapToJSON() (string, error) {
	jsonBytes, err := json.Marshal(peers.peerMap)
	return string(jsonBytes), err
}

/* TestPeerListRebalance()
*
* To Test Rebalance() function
*
 */
func TestPeerListRebalance() {
	peers := NewPeerList(5, 4)
	peers.Add("1111", 1)
	peers.Add("4444", 4)
	peers.Add("-1-1", -1)
	peers.Add("0000", 0)
	peers.Add("2121", 21)
	peers.Rebalance()
	expected := NewPeerList(5, 4)
	expected.Add("1111", 1)
	expected.Add("4444", 4)
	expected.Add("2121", 21)
	expected.Add("-1-1", -1)
	fmt.Println(reflect.DeepEqual(peers, expected))

	peers = NewPeerList(5, 2)
	peers.Add("1111", 1)
	peers.Add("4444", 4)
	peers.Add("-1-1", -1)
	peers.Add("0000", 0)
	peers.Add("2121", 21)
	peers.Rebalance()
	expected = NewPeerList(5, 2)
	expected.Add("4444", 4)
	expected.Add("2121", 21)
	fmt.Println(reflect.DeepEqual(peers, expected))

	peers = NewPeerList(5, 4)
	peers.Add("1111", 1)
	peers.Add("7777", 7)
	peers.Add("9999", 9)
	peers.Add("11111111", 11)
	peers.Add("2020", 20)
	peers.Rebalance()
	expected = NewPeerList(5, 4)
	expected.Add("1111", 1)
	expected.Add("7777", 7)
	expected.Add("9999", 9)
	expected.Add("2020", 20)
	fmt.Println(reflect.DeepEqual(peers, expected))
	peers.Show()
}

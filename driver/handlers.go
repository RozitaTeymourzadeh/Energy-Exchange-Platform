package driver

import (
	//"MerklePatriciaTree/p3/cs686-blockchain-p3-RozitaTeymourzadeh/p2"
	//"MerklePatriciaTree/p5/Blockchain_Application_P5/data"
	//"MerklePatriciaTree/p5/Blockchain_Application_P5/p4"
	"bytes"
	"crypto"
	"crypto/rsa"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/crypto/sha3"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var SBC SyncBlockChain
var Peers PeerList
var Heart HeartBeatData
var ifStarted bool

//var FRONT_SERVER = "http://192.168.1.106:7787"
var TA_SERVER = "http://192.168.1.106:6686"
var FIRST_PEER = "192.168.1.106:6686" //
var SELF_ADDR = "192.168.1.106:6686"  //"localhost:6686"
var REGISTER_SERVER = TA_SERVER + "/peer"
var ASK_PEERS = "/block"
var EDGEX_IP = "localhost"
var SELF_ID = 0

var BC_DOWNLOAD_SERVER = FIRST_PEER + "/upload"
var RECEIVE_PATH = "/heartbeat/receive"
var STOP_GEN_BLOCK = false
var NONCE_ZERO = "000"

//var newTransactionObject Transaction
var TxPool TransactionPool
var Balance int = 10000
var TransactionMap map[string]Transaction
var minerKey *rsa.PrivateKey

//var Verified bool = false
//var enoughBalance bool = false
var transactionReady = false
var SignatureV []byte
var hashed []byte
var h crypto.Hash

/* init()
*
* Initialization
*
 */
func init() {

	fmt.Println("......Initialization ....")
	TxPool = NewTransactionPool()
	SBC = NewSyncBlockChain()
	Peers = NewPeerList(Peers.GetSelfId(), 32)
	ifStarted = false
	TransactionMap = make(map[string]Transaction)

	/*Init Block*/
	mpt := MerklePatriciaTrie{}
	mpt.Initial()
	mpt.Insert(StringRandom(2), StringRandom(5))
	block := Block{}
	block.Initial(1, "gensis", mpt, NONCE_ZERO)
	block.Header.Nonce = NONCE_ZERO
	SBC.Insert(block)
	if len(os.Args) > 1 {
		responseString := os.Args[1]
		fmt.Println(responseString)
		result, err := strconv.ParseInt(responseString, 10, 32)
		if err != nil {
			panic(err)
		}
		id := int32(result)
		Peers.Register(id)
		//SELF_ADDR="localhost"+os.Args[1]
		SELF_ADDR = SystemIp() + ":" + os.Args[1]
		Peers.Add(FIRST_PEER, 6686)
		//publicKey,_:=data.ParseRsaPublicKeyFromPemStr("HARD_CODED_PEER1")
		//Peers.AddPublicKey(publicKey,6686)
	} else {
		Peers.Register(6686)
		//SELF_ADDR="localhost:6686"
		SELF_ADDR = SystemIp() + ":6686"
	}

}

/* Start
*
* To start the application
* Register ID, download BlockChain, start HeartBeat
*
 */
func Start(w http.ResponseWriter, r *http.Request) {

	/*Register*/
	if ifStarted == true {
		fmt.Println("Register on TA Server!")
		Register()
	} else {
		fmt.Println("Port Register!", Peers.GetSelfId())
	}
	fmt.Println("Host is : ", r.Host)
	//SELF_ADDR = r.Host
	//if r.Host != FIRST_PEER {
	if SELF_ADDR != FIRST_PEER {
		fmt.Println("Not First node: download().")
		Download()
	} else {
		fmt.Println("First node: Skip downloading.")
	}

	/* Start Trying Nonce */
	minerKey = GenerateKeyPair(4096)
	fmt.Println("Public Key:", minerKey.PublicKey)
	fmt.Println("Private Key:", minerKey)
	//Peers.AddPublicKey(&minerKey.PublicKey,Peers.GetSelfId())
	go StartTryingNonce()

	fmt.Println("Calling Periodic")
	go Periodic() // MAKING DECISION

	/*Timer to send heartBeat periodically*/
	ticker := time.NewTicker(10 * time.Second)
	quit := make(chan struct{})
	for {
		select {
		case <-ticker.C:
			StartHeartBeat()
		case <-quit:
			ticker.Stop()
			return
		}
	}
}

/* Show
*
* Shows the PeerMap and the BlockChain
*
 */
func Show(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s\n %s", Peers.Show(), SBC.Show())
}

/* Canonical
*
* Shows the BlockChain after POW
*
 */
func Canonical(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s\n", SBC.Canonical())
}

/* Register
*
* Register to TA's server
* get an ID and register to PeerList
*
 */
func Register() {
	response, err := http.Get(REGISTER_SERVER)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	responseString := string(responseData)
	fmt.Println(responseString)
	result, err := strconv.ParseInt(responseString, 10, 32)
	if err != nil {
		panic(err)
	}
	id := int32(result)
	Peers = NewPeerList(id, 32)
}

/* Download
*
* Download blockchain from TA server
*
 */
func Download() {
	uploadAddress := "http://" + FIRST_PEER + "/upload"
	fmt.Println("Download from URL:" + uploadAddress)
	peerMapStringValue, _ := Peers.EncodePeerMapToJSON()
	registerData := NewRegisterData(Peers.GetSelfId(), peerMapStringValue)
	jsonBytes, err := json.Marshal(registerData)
	req, err := http.NewRequest("POST", uploadAddress, bytes.NewBuffer(jsonBytes))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(req)
	fmt.Println("---> Upload Request from:[", FIRST_PEER, "] --->")
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	fmt.Println("response Status:", response.Status)
	fmt.Println("response Headers:", response.Header)
	fmt.Println("response from peer:", response.Body)
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	responseString := string(responseData)
	SBC.UpdateEntireBlockChain(responseString)
}

/* Upload()
*
* Return the BlockChain's JSON for any peer who call this function
*
 */
func Upload(w http.ResponseWriter, r *http.Request) {
	fmt.Println(".....Upload .....")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	log.Println(string(body))
	registerData := NewRegisterData(0, "")
	registerData.DecodeFromJSON(string(body))
	fmt.Println("AssignedId:", registerData.AssignedId)
	fmt.Println("PeerMapJson:", registerData.PeerMapJson)
	blockChainJSON, err := SBC.BlockChainToJson()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
	fmt.Fprint(w, blockChainJSON)
}

/* UploadBlock
*
* Upload a block to whoever called this method, return jsonStr
* Return the Block's JSON
*
 */
func UploadBlock(w http.ResponseWriter, r *http.Request) {
	fmt.Println("<--- UploadBlock Received  From :[", r.Host, "] <---")
	param := strings.Split(r.URL.Path, "/")
	h, err := strconv.ParseInt(param[2], 10, 32)
	fmt.Println("param0:", param[0])
	fmt.Println("param1:", param[1])
	fmt.Println("param2:", param[2])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: InternalServerError. " + err.Error()))
	} else {
		encode := param[3]
		block, flag := SBC.GetBlock(int32(h), encode)
		if flag == false {
			w.WriteHeader(http.StatusNoContent)
		} else {
			blockStr, err := block.EncodeToJSON()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("HTTP 500: InternalServerError. " + err.Error()))
			} else {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(blockStr))
			}
		}
	}
}

/* HeartBeatReceive
*
* Add the remote address, and the PeerMapJSON into local PeerMap
* Then check if the HeartBeatData contains a new block
*
 */
func HeartBeatReceive(w http.ResponseWriter, r *http.Request) {
	var mutex = sync.Mutex{}
	mutex.Lock()
	defer mutex.Unlock()

	transaction := Transaction{}

	defer r.Body.Close()
	data, _ := ioutil.ReadAll(r.Body)
	fmt.Println("<---- HeartBeat Received  From:[", r.Host, "] <<<<<<<<<<<<")
	fmt.Fprintf(w, "%s\n", r.Host)
	fmt.Fprintf(w, "%s", string(data))
	error := json.Unmarshal(data, &Heart)
	if error != nil {
		fmt.Println("Error occured in HeartBeatReceive: ", error)
	} else {

	}
	//fmt.Println("HeartBeatVariable.Addr", HeartBeatVariable.Addr)
	//fmt.Println("SELF_ADDR", SELF_ADDR)
	if Heart.Addr == SELF_ADDR {
		return
	}

	transaction.DecodeFromJson(Heart.TransactionInfoJson)

	//Peers.AddPublicKey(Heart.PeerPublicKey, Heart.Id)
	Peers.Add(Heart.Addr, Heart.Id)
	Peers.InjectPeerMapJson(Heart.PeerMapJson, SELF_ADDR)
	if Heart.IfNewBlock {
		fmt.Println("HeartBeat flag is true!!")

		heartBlock := Block{}
		heartBlock.DecodeFromJSON(Heart.BlockJson)
		fmt.Println("Received block! Root:", heartBlock.Value.GetRoot())
		receivedPuzzle := heartBlock.Header.ParentHash + heartBlock.Header.Nonce + heartBlock.Value.GetRoot()
		sum := sha3.Sum256([]byte(receivedPuzzle))
		if strings.HasPrefix(hex.EncodeToString(sum[:]), NONCE_ZERO) {
			fmt.Println("Block with SPECIAL PREFIX arrived from:[", r.Host, "]")
			latestBlocks := SBC.GetLatestBlocks()
			for i := 0; i < len(latestBlocks); i++ {
				if latestBlocks[i].Header.Hash == heartBlock.Header.ParentHash {
					STOP_GEN_BLOCK = true
					break
				}
			}
			if heartBlock.Header.Height == 1 {
				if TxPool.CheckConfirmedPool(transaction) == false {
					SBC.Insert(heartBlock)
					TxPool.DeleteFromTransactionPool(transaction.EventId)
					TxPool.AddToConfirmedPool(transaction)
				} else {
					TxPool.DeleteFromTransactionPool(transaction.EventId)
				}
			} else {
				_, flag := SBC.GetBlock(heartBlock.Header.Height-1, heartBlock.Header.ParentHash)
				if flag {
					if TxPool.CheckConfirmedPool(transaction) == false {
						SBC.Insert(heartBlock)
						TxPool.DeleteFromTransactionPool(transaction.EventId)
						TxPool.AddToConfirmedPool(transaction)
					} else {
						TxPool.DeleteFromTransactionPool(transaction.EventId)
					}
				} else {
					fmt.Println("Gap.Inserting Heart Beat Block:", heartBlock)
					AskForBlock(heartBlock.Header.Height-1, heartBlock.Header.ParentHash)
					if TxPool.CheckConfirmedPool(transaction) == false {
						SBC.Insert(heartBlock)
						TxPool.DeleteFromTransactionPool(transaction.EventId)
						TxPool.AddToConfirmedPool(transaction)
					} else {
						TxPool.DeleteFromTransactionPool(transaction.EventId)
					}
				}
			}
		} else {
			fmt.Println("Unmatched Puzzle! Calculated Puzzle:", hex.EncodeToString(sum[:]))
			fmt.Println("Incoming Heart Beat Block.Hash::", heartBlock.Header.Hash)
			fmt.Println("Incoming Heart Beat Nonce:", heartBlock.Header.Nonce)
			fmt.Println("Incoming Heart Beat mpt.Root:", heartBlock.Value.GetRoot())
			fmt.Println("Calculated Incoming Hash Puzzle:", hex.EncodeToString(sum[:]))
		}
	} else {
		fmt.Println("There is no block in heartBeat!")
	}
	Heart.Hops -= 1
	if Heart.Hops > 0 {
		Heart.Addr = SELF_ADDR
		Heart.Id = Peers.GetSelfId()
		ForwardHeartBeat(Heart)
	}
}

/* AskForBlock
*
* Ask another server to return a block of certain height and hash
* http get to /localhost:port/block/{height}/{hash} (UploadBlock) to get the Block
*
 */

func AskForBlock(height int32, hash string) {
	PeerMap := Peers.GetPeerMap()
	for key, value := range PeerMap {
		fmt.Printf("key[%s] value[%d]\n", key, value)
		fmt.Println("height:", height)
		heightString := ConvertIntToString(height)
		//http://localhost:8863/upload/1/xxxx
		prepareRequest := "http://" + key + ASK_PEERS + "/" + heightString + "/" + hash
		//http://localhost:8845:/block/5/hashvalue
		fmt.Println("PrepareRequest:", prepareRequest)
		response, err := http.Get(prepareRequest)
		fmt.Println("---> AskForBlock Sent  --->")
		if response.StatusCode == 204 {
			fmt.Println("There is no block available on Peer:", key)
		} else if response.StatusCode == 500 {
			fmt.Println("Internal Server Error on Peer:", key)
		} else {
			if err != nil {
				log.Fatal(err)
			}
			defer response.Body.Close()
			responseData, err := ioutil.ReadAll(response.Body)
			if err != nil {
				log.Fatal(err)
			}
			responseString := string(responseData)
			fmt.Println("Missing Block's JSON Response:", responseString)
			missingBlock := Block{}
			missingBlock.DecodeFromJSON(responseString)
			//SBC.Insert(missingBlock)
			fmt.Println("Get Block after AskForBlock:", missingBlock)
			fmt.Println("From peer:", key)
			if !SBC.CheckParentHash(missingBlock) {
				AskForBlock(missingBlock.Header.Height-1, missingBlock.Header.ParentHash)
				fmt.Println("Get Block after recursively AskForBlock:", missingBlock)
				fmt.Println("From peer:", key)
			}
			SBC.Insert(missingBlock)
		}
	}
}

/* ForwardHeartBeat
*
* Send the HeartBeatData to all peers in local PeerMap after receiving
* Let other peers know about newly generated block
*
 */
func ForwardHeartBeat(heartBeatData HeartBeatData) {
	if heartBeatData.Hops != 0 {
		heartBeatData, _ := json.Marshal(heartBeatData)
		Peers.Rebalance()
		for addr, _ := range Peers.Copy() {
			resp, err := http.Post("http://"+addr+"/heartbeat/receive", "application/json; charset=UTF-8", strings.NewReader(string(heartBeatData)))
			if err != nil || resp.StatusCode != 200 {
				fmt.Println("FATAL : ", err)
				Peers.Delete(addr)
			}
		}
	}
}

func SendHeartBeatToSelf(hbd HeartBeatData) {
	hbd.Hops = 1
	heartBeatData, _ := json.Marshal(hbd)
	resp, err := http.Post("http://"+GetNodeId().GetAddressPort()+"/heartbeat/receive",
		"application/json; charset=UTF-8", strings.NewReader(string(heartBeatData)))
	if err != nil || resp.StatusCode != 200 {
		fmt.Println("FATAL : in SendHeartBeatToSelf : ", err)
	}
}

/* StartHeartBeat
*
* Start a while loop. Inside the loop, sleep for randomly 5~10 seconds,
* then use PrepareHeartBeatData() to create a
* HeartBeatData, and send it to all peers in the local PeerMap.
*
 */
func StartHeartBeat() {
	fmt.Println("Start Heart Beat!")
	Peers.Rebalance()
	PeerMap := Peers.GetPeerMap()
	fmt.Println("Size of PeerMap:", len(PeerMap))

	for key, value := range PeerMap {
		fmt.Printf("key[%s] value[%d]\n", key, value)
		uploadAddress := "http://" + key + "/heartbeat/receive"
		fmt.Println("/heartbeat/receive Request will be sent to :" + uploadAddress)
		//destination := "http://localhost:6688" +/heartbeat/receive"
		peerMapToJson, err := Peers.PeerMapToJson()
		if err != nil {
			log.Fatal(err)
		}

		mpt := MerklePatriciaTrie{}
		heartBearData := PrepareHeartBeatData(&SBC, Peers.GetSelfId(), peerMapToJson, SELF_ADDR, false, "", mpt, &minerKey.PublicKey, false, "")
		jsonBytes, err := json.Marshal(heartBearData)
		req, err := http.NewRequest("POST", uploadAddress, bytes.NewBuffer(jsonBytes))
		req.Header.Set("X-Custom-Header", "myvalue")
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		fmt.Println("---> HeartBeatSent To:[", key, "]")
		if err != nil {
			fmt.Println("Peer issue:  ", key)
			Peers.Delete(key)
			return
		}
		defer resp.Body.Close()
		fmt.Println("response Status:", resp.Status)
		fmt.Println("response Headers:", resp.Header)
		body, _ := ioutil.ReadAll(resp.Body)
		if heartBearData.IfNewBlock || heartBearData.IfValidTransaction {
			fmt.Println("response Body:", string(body))
		}

	}
}

/* StartTryingNonce()
*
* Generate block after solving hash puzzel
*
 */
func StartTryingNonce() {
	var mutex = sync.Mutex{}
	mutex.Lock()
	defer mutex.Unlock()
	isValidTransaction := false
	for {
	GetLatestBlock:
		newMpt := MerklePatriciaTrie{}
		newMpt.Initial()
		blocks := SBC.GetLatestBlocks()
		STOP_GEN_BLOCK = false
		var transactionJSON string
		transaction := TxPool.GetOneTxFromPool(TxPool, Balance)
		if transaction != nil {
			transactionJSON, _ = transaction.EncodeToJson()

			//if true /*Verified && enoughBalance*/ {
			//
			//	encryptedPKCS1v15 := []byte(transactionJSON)
			//	decryptedPKCS1v15 := DecryptPKCS(minerKey, encryptedPKCS1v15)
			//	Verified, _ := VerifyPKCS(&minerKey.PublicKey, h, hashed, SignatureV)
			//	fmt.Println("Decrypted message is: ", decryptedPKCS1v15)
			//	fmt.Println("Signature verified: ", Verified)
			//	fmt.Println("User has enough money: ", enoughBalance)
			//
			//}
			if len(transactionJSON) > 0 {
				newMpt.Insert(transaction.EventId, transactionJSON)
				validateNonce := StringRandom(16)
				hashPuzzle := string(blocks[0].Header.Hash) + string(validateNonce) + string(newMpt.GetRoot())
				sum := sha3.Sum256([]byte(hashPuzzle))
				if strings.HasPrefix(hex.EncodeToString(sum[:]), NONCE_ZERO) {
					peerMapJson, _ := Peers.PeerMapToJson()
					transactionJSON, _ = transaction.EncodeToJson()
					heartBeatData := PrepareHeartBeatData(&SBC, Peers.GetSelfId(), peerMapJson, SELF_ADDR, true, validateNonce,
						newMpt, &minerKey.PublicKey, isValidTransaction, transactionJSON)
					heartBlock := Block{}
					heartBlock.DecodeFromJSON(heartBeatData.BlockJson)
					fmt.Println("heartBlock.Root:", string(heartBlock.Value.GetRoot()))
					testPuzzle := string(heartBlock.Header.ParentHash) + string(heartBlock.Header.Nonce) + string(heartBlock.Value.GetRoot())
					sum = sha3.Sum256([]byte(testPuzzle))
					//SendHeartBeatToSelf(heartBeatData)  // todo :test
					SBC.Insert(heartBlock) //todo: test
					ForwardHeartBeat(heartBeatData)
					isValidTransaction = false
					TxPool.DeleteFromTransactionPool(transaction.EventId)
					if STOP_GEN_BLOCK {
						fmt.Println("Stop Generating Block.")
						goto GetLatestBlock
					}
				}
			}

		} else {
		}
	}
}

/* Event()
*
* /getevent API
* /postevent API
* To enter the event info
 */
func Event(w http.ResponseWriter, r *http.Request) {

	log.Println(".....Event method .....")

	switch r.Method {
	case "POST":
		//read request body
		body, err := readRequestBody(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}

		//convert json in body to tx
		tx := TransactionFromJSON(body)

		fmt.Println("Received body : " + string(body))

		if tx.EventId != "" {
			go TxPool.AddToTransactionPool(tx)
			fmt.Fprintf(w, "Tx added to pool !-) Balance = %d\n", Balance)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	default:
		fmt.Fprintf(w, "FATAL: Wrong HTTP Request!")
	}
}

/* Event()
*
* /getQueryEvent API
* /postQueryEvent API
* To search the event info
 */
func QueryEvent(w http.ResponseWriter, r *http.Request) {
	log.Println("Event method is triggered!")

	switch r.Method {
	case "GET":
		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("PWD:", dir)

		http.ServeFile(w, r, "QueryEvent.html")
	case "POST":
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		eventId := r.FormValue("eventId")
		fmt.Fprintf(w, "Event ID = %s\n", eventId)
		fmt.Fprintf(w, "SEARCH RESULT = %s\n", SBC.GetEventInfornation(eventId))
	default:
		fmt.Fprintf(w, "Wrong method !")
	}
}

func readRequestBody(r *http.Request) ([]byte, error) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return []byte{}, errors.New("cannot read request body")
	}
	defer r.Body.Close()
	return reqBody, nil
}

func GetTxpool(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	sb := strings.Builder{}
	sb.WriteString("Tx pool:\n\n")
	for _, tx := range TxPool.Pool {
		txJson, _ := tx.TransactionToJSON()
		sb.WriteString(string(txJson) + "\n\n")
	}
	w.Write([]byte(sb.String()))
}

//func Showsbc(w http.ResponseWriter, r *http.Request) {
//
//	fmt.Println("showsbc : ")
//
//	tempBlockSlice := make([]Block, 0)
//	height := SBC.bc.Length // length of blockchain
//
//	for height > 1 {
//		forkBlocks, _ := SBC.Get(height)
//		if len(forkBlocks) == 1 {
//			tempBlock := forkBlocks[0]
//			tempBlockSlice = append(tempBlockSlice, tempBlock) // appending
//			if height > 2 {
//				tempParentBlock := SBC.GetParentBlock(tempBlock)
//				//tempBlockSlice = append([]Block{tempParentBlock}, tempBlockSlice...) // prepending
//				tempBlock = tempParentBlock
//				height = tempBlock.Header.Height
//			}
//		}
//	}
//
//	//tempTxsSlice := make([]Transaction, 0)
//	str := strings.Builder{}
//	str.WriteString("show sbc :")
//
//	for _, block := range tempBlockSlice {
//		txList := block.Value.LeafList()
//		for _, txStr := range txList {
//			//tx := TransactionFromJSON([]byte(txStr))
//			str.WriteString(txStr + "\n")
//		}
//	}
//
//	w.WriteHeader(http.StatusOK)
//	w.Write([]byte(str.String()))
//
//}

func ShowTxList(w http.ResponseWriter, r *http.Request) {
	p := PageVars{
		Title:                 APPNAME,
		IpPort:                GetNodeId().Address + ":" + GetNodeId().Port,
		DeviceMap:             DEVICELIST.Devices, //SELFDEVICES.DeviceMapToList(),
		SupplyDevicesDetails:  SUPPLYDEVICEDETAILS,
		ConsumeDevicesDetails: CONSUMEDEVICEDETAILS,
		Transactions:          ALLTRANSACTIONS,
		SdReadings:            GetLast100SDReadings(),
		CdReadings:            GetLast100CDReadings(),
	}

	//x := p.SupplyDevicesDetails
	//fmt.Println(len(x))
	render(w, "txns.html", p)
}

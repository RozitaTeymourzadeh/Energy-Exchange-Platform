package driver

import (
	"bytes"
	"fmt"
	"strconv"

	"io/ioutil"
	"net/http"
	"time"
)

//////////////////////////
//  Device Manager periodic tasks
/////////////////////////

func Periodic() {
	fmt.Println("Periodic")
	time.Sleep(1 * time.Second)
	for true {
		readBlockchainAndUpdateStates(LASTREADFORHEIGHT) // updated open tx map
		makeSupplyDecision()

		generateVarsForIndex()

		time.Sleep(10 * time.Second)
	}
}

// read Blockchain And Update States
func readBlockchainAndUpdateStates(lastReadForHeight int) {
	tempBlockSlice := make([]Block, 0)
	fmt.Println("readBlockchainAndUpdateStates : lastReadForHeight : ", lastReadForHeight)
	sbcLength := int(SBC.bc.Length) //canonical 2 height below current length
	if sbcLength >= lastReadForHeight {
		tempBlockSlice = fillTempBlockSlice(tempBlockSlice, lastReadForHeight, sbcLength)

		updateDevicesStatesHelper(tempBlockSlice)

		LASTREADFORHEIGHT = sbcLength + 1
	}

}

// get slice of blocks
func fillTempBlockSlice(tempBlockSlice []Block, lastReadForHeight int, sbcLength int) []Block {
	height := int32(sbcLength)
	forkBlocks, _ := SBC.Get(height)
	if len(forkBlocks) == 1 {
		tempBlock := forkBlocks[0]
		tempBlockSlice = append([]Block{tempBlock}, tempBlockSlice...)
		for height > int32(lastReadForHeight) {
			tempParentBlock := SBC.GetParentBlock(tempBlock)
			tempBlockSlice = append([]Block{tempParentBlock}, tempBlockSlice...) // prepending

			tempBlock = tempParentBlock
			height = tempBlock.Header.Height
		}
	}
	return tempBlockSlice
}

//update device states
func updateDevicesStatesHelper(tempBlockSlice []Block) {
	for _, block := range tempBlockSlice {
		txList := block.Value.LeafList()
		for _, txStr := range txList {
			tx := TransactionFromJSON([]byte(txStr))

			updateDeviceStateForTx(tx)

			fmt.Println("TxStr : " + txStr)
		}
	}
}

// update device states for a tx
func updateDeviceStateForTx(tx Transaction) {
	if tx.EventId != "" {
		ALLTRANSACTIONS = append([]Transaction{tx}, ALLTRANSACTIONS...)

		if tx.EventType == "require" {
			updateDeviceStateForRequireTx(tx)
		}
		if tx.EventType == "supply" {
			updateDeviceStateForSupplyTx(tx)

		}
	}
}

// update device states for require tx
func updateDeviceStateForRequireTx(tx Transaction) {
	if _, ok := OPENCONSUMETXS.Pool[tx.EventId]; !ok {
		// deduct balance for self tx
		if tx.ConsumerId == GetConsumeDeviceId() &&
			tx.ConsumerAddress == GetConsumeDeviceAddress() { // self tx
			Balance -= tx.EventFee
			TRANSACTIONS = append([]Transaction{tx}, TRANSACTIONS...)
		}
		fmt.Println("Adding tx to OPENCONSUMETXS")
		OPENCONSUMETXS.Pool[tx.EventId] = tx
	}
}

// update device states for supply tx
func updateDeviceStateForSupplyTx(tx Transaction) {
	OPENCONSUMETXS.DeleteFromTransactionPool(tx.RequireEventId)
	toSupply, _ := strconv.Atoi(tx.SupplierToSupply)
	buyRate, _ := strconv.Atoi(tx.BuyRate)
	supplierSupplyRate, _ := strconv.Atoi(tx.SupplierSupplyRate)

	selfInSupplierField := false
	selfInConsumerField := false

	if tx.SupplierAddress == GetSupplyDeviceAddress() &&
		tx.SupplierId == GetSupplyDeviceId() { // Supplier - self in supply tx
		// deduct balance for self tx
		Balance -= tx.EventFee
		SetIsSupplying(1)
		SetToSupply(toSupply)
		fmt.Println("proxy of  : send energy to consumer")
		Balance += toSupply * buyRate

		selfInSupplierField = true
	}

	if tx.ConsumerAddress == GetConsumeDeviceAddress() &&
		tx.ConsumerId == GetConsumeDeviceId() { // consumer - self in supply tx
		SetIsReceiving(1)
		SetToReceive(toSupply)
		SetToReceiveRate(supplierSupplyRate)
		fmt.Println("proxy of  : receive energy from supplier") // todo
		Balance -= toSupply * buyRate

		selfInConsumerField = true
	}

	if selfInSupplierField || selfInConsumerField {
		TRANSACTIONS = append([]Transaction{tx}, TRANSACTIONS...)
	}
}

/*
makeSupplyDecision checks if it has supplydevice
if it has buyRate >= sellRate && surplus > 0
*/
func makeSupplyDecision() {
	fmt.Println("In makeSupplyDecision")
	isSupplying := GetIsSupplying()
	hasOffered := GetHasOffered()
	offeredAtTime := GetHasOfferedAtTime()
	timeNow := time.Now()
	duration := timeNow.Sub(offeredAtTime).Seconds()

	if isSupplying == 0 && hasOffered == false {
		fmt.Println("if GetIsSupplying() == 0")
		surplus := GetSurplus()
		sellRate := GetSellRate()

		for _, cnTx := range OPENCONSUMETXS.Pool {
			fmt.Println("in for _, cnTx := range OPENCONSUMETXS.Pool")
			require, _ := strconv.Atoi(cnTx.ConsumerRequire)
			buyRate, _ := strconv.Atoi(cnTx.BuyRate)
			if buyRate >= sellRate && surplus > 0 { // condition to sell
				fmt.Println("if buyRate >= sellRate && surplus > 0")
				if require <= surplus {
					spTx := createSupplyTx(cnTx, strconv.Itoa(require))
					sendSpTxToAll(spTx)
					break
				} else {
					spTx := createSupplyTx(cnTx, strconv.Itoa(surplus))
					sendSpTxToAll(spTx)
					break
				}
			}
		}
	}

	if duration > 45 && hasOffered && isSupplying == 0 {
		SetHasOffered(false)
	}

}

// creates supply tx
func createSupplyTx(cnTx Transaction, toSupply string) Transaction {
	newTx := NewTransaction("supply", cnTx.ConsumerName, cnTx.ConsumerId, cnTx.ConsumerAddress,
		cnTx.ConsumerRequire, cnTx.ConsumerCharge, cnTx.ConsumerDischargeRate, cnTx.BuyRate, cnTx.EventId,
		GetSupplyDeviceName(), GetSupplyDeviceId(), GetSupplyDeviceAddress(), toSupply, /*cnTx.ConsumerRequire,*/
		strconv.Itoa(GetSupplyRate()), Balance)

	return newTx
}

// send supply tx to all peers
func sendSpTxToAll(newTx Transaction) {
	SetHasOffered(true)               // setting has offered to true
	SetHasOfferedAtHeight(time.Now()) // setting has offered at time

	body, err := newTx.TransactionToJSON()
	if err == nil {
		uri := "http://" + GetNodeId().Address + ":" + GetNodeId().Port + "/postevent"
		fmt.Println("supply tx to : " + uri)
		http.Post(uri, "application/json", bytes.NewBuffer(body))
		//SetHasOffered(true)
		for peer, _ := range Peers.Copy() {
			uri := "http//:" + peer + "/postevent"
			fmt.Println("supply tx to : " + uri)
			http.Post(uri, "application/json", bytes.NewBuffer(body))
		}
	}

}

// generates Vars For Index
func generateVarsForIndex() {
	DEVICELIST = getAllDevices( /*data.GetNodeId().ConnectingAddress*/ )
	SUPPLYDEVICEDETAILS = generateSupplyDeviceTypeBoard("supply")
	CONSUMEDEVICEDETAILS = generateConsumeDeviceTypeBoard("consume")

	fmt.Println("Decision : ")
	if len(CONSUMEDEVICEDETAILS) < 1 {
		//_, _ = w.Write([]byte("No consume device"))
		fmt.Println("No consume device")

	} else if len(SUPPLYDEVICEDETAILS) < 1 {
		//_, _ = w.Write([]byte("No supply device"))
		fmt.Println("No supply device")

	} else {

		str := "Dummy makeDecisionHandlerHelper : todo : makeDecisionHandlerHelper()" //
		//str := makeDecisionHandlerHelper()
		fmt.Println(str)

	} // end of else

}

// get all devices
func getAllDevices() DeviceList {
	dl := NewDeviceList()
	// start of get self devices ////
	devices := updateDeviceListWithSelfDevices()
	for _, device := range devices.Devices {
		device.PeerId = GetNodeId().Address + ":" + GetNodeId().Port
		dl.Devices = append(dl.Devices, device) /// to read SBC and create board
	}
	return dl
}

// update DeviceList With SelfDevices
func updateDeviceListWithSelfDevices() DeviceList { //
	//data.GetSupplyDeviceBoard()
	//todo : update by reading canonical SBC
	resp, err := http.Get("http://" + GetNodeId().Address + ":" + GetNodeId().Port + "/" + "getallselfdevices")
	if err != nil {
		fmt.Println("Error in getting all devices : in : updateDeviceTypeBoards")
	}
	defer resp.Body.Close()
	bytesRead, _ := ioutil.ReadAll(resp.Body)

	deviceList := DeviceListFromJson(bytesRead)
	for _, device := range deviceList.Devices {
		fmt.Println("In updateDeviceTypeBoards : " + device.PeerId + " - " + device.Id + " - " + device.Name)
	}
	return deviceList

}

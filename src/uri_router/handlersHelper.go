package uri_router

import (
	"bytes"
	"fmt"
	"github.com/edgexfoundry/device-simple/src/data"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

//for return code 503
func returnCode503(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Server Error", http.StatusServiceUnavailable)
	http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
}

//for return code 500
func returnCode500(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Server Error", http.StatusInternalServerError)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

//for return code 204
func returnCode204(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Block does not exists", http.StatusNoContent)
	http.Error(w, http.StatusText(http.StatusNoContent), http.StatusNoContent)
}

func getAllDevices( /*ip string*/ ) data.DeviceList {
	dl := data.NewDeviceList()

	for _, peer := range data.GetNodeId().GetPeers() {
		//uri := "http://" + peer.Address + ":48082/api/v1/device"
		uri := "http://" + peer.Address + ":9999/sendDeviceList"
		fmt.Println("Sending device req to : ", peer.Address)
		resp, err := http.Get(uri)
		if err != nil {
			fmt.Println("Error in getting all devices")
		}
		defer resp.Body.Close()
		bytesRead, _ := ioutil.ReadAll(resp.Body)
		peerDeviceList := data.DeviceListFromJson(bytesRead)
		for _, val := range peerDeviceList.Devices {
			val.PeerId = peer.Address
			dl.Devices = append(dl.Devices, val)
		}

	}

	return dl

	//
	//
	//uri := "http://"+data.GetNodeId().Address+":48082/api/v1/device"
	//
	//resp, err := http.Get(uri)
	//if err != nil {
	//	fmt.Println("Error in getting all devices")
	//}
	//defer resp.Body.Close()
	//bytesRead, _ := ioutil.ReadAll(resp.Body)
	//
	//return data.DeviceListFromJson(bytesRead)
}

//1. get list of peers
//	2. iterate and get list of devices for the peer
//		3. for device with name
//			4. if name contains "supply"
//				5. add to supplyDeviceList
//
//use by PageVars
func generateDeviceTypeBoard(deviceType string) []data.DeviceTypeDetails { //todo : here build suppy board

	sl := make([]data.DeviceTypeDetails, 0)

	for _, d := range DeviceList.Devices {

		uri := "http://" + d.PeerId + ":48080/api/v1/event/device/" + d.Name + "/" + "1"
		resp, err := http.Get(uri)

		if err != nil {
			fmt.Println("error in reading response body in startreading")
		}
		defer resp.Body.Close()
		bytesRead, _ := ioutil.ReadAll(resp.Body)
		//fmt.Println(string(bytesRead))

		cdes := data.CoreDataEventsFromJson(bytesRead)
		// iterate through cdes to get only "supply" device
		//for _, cde := range cdes.DataEvents {
		if strings.Contains(strings.ToLower(cdes.DataEvents[0].Device), strings.ToLower(deviceType)) {
			sl = append(sl, data.DeviceTypeDetails{
				DeviceAddress: d.PeerId,
				DeviceName:    cdes.DataEvents[0].Device,
				Id:            cdes.DataEvents[0].Id,
				Reading:       cdes.DataEvents[0].Readings[0].Value,
			})
		}
		//}

	}
	return sl
}

func sendTransactionToSupplier(tx data.Transaction) {
	txJson, err := tx.TransactionToJSON()
	if err != nil {
		log.Print("Cannot create transaction")
		return
	}

	//uri := "http://localhost:48081/api/v1/deviceprofile/id/" + vars["deviceId"]
	uri := "http://" + tx.SupplierAddress + ":9999/suppliertx"
	log.Println("sending post req to : " + uri)
	client := &http.Client{}
	// creating request
	req, _ := http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(txJson))
	// fetching response
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error in reading response body in start reading")
	}
	defer resp.Body.Close()
}

func sendTransactionToConsumer(tx data.Transaction) {
	txJson, err := tx.TransactionToJSON()
	if err != nil {
		log.Print("Cannot create transaction")
		return
	}

	//uri := "http://localhost:48081/api/v1/deviceprofile/id/" + vars["deviceId"]
	uri := "http://" + tx.ConsumerAddress + ":9999/consumertx"
	log.Println("sending post req to : " + uri)
	client := &http.Client{}
	// creating request
	req, _ := http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(txJson))
	// fetching response
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error in reading response body in start reading")
	}
	defer resp.Body.Close()
}

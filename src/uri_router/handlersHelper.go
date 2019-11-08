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

func getAllDevices() data.DeviceList {
	dl := data.NewDeviceList()
	if len(data.GetNodeId().GetPeers()) > 0 {
		for _, peer := range data.GetNodeId().GetPeers() {
			uri := "http://" + peer.IpAdd + ":" + peer.Port + "/sendDeviceList"
			fmt.Println("Sending device req to : ", uri)
			resp, err := http.Get(uri)
			if err != nil {
				fmt.Println("Error in getting all devices")
			}
			defer resp.Body.Close()
			bytesRead, _ := ioutil.ReadAll(resp.Body)
			peerDeviceList := data.DeviceListFromJson(bytesRead)
			for _, val := range peerDeviceList.Devices {
				val.PeerId = peer.IpAdd
				dl.Devices = append(dl.Devices, val)
			}

		}
	}

	return dl
}

//1. get list of peers
//	2. iterate and get list of devices for the peer
//		3. for device with name
//			4. if name contains "supply"
//				5. add to supplyDeviceList
//
//use by PageVars
func generateDeviceTypeBoard(deviceType string) []data.DeviceTypeDetails {

	sl := make([]data.DeviceTypeDetails, 0)

	for _, d := range DEVICELIST.Devices {

		//uri := "http://" + d.PeerId + ":48080/api/v1/event/device/" + d.Name + "/" + "10"
		uri := "http://" + d.PeerId + ":9999/sendDeviceEvents/" + d.Name + "/" + "10"
		resp, err := http.Get(uri)

		if err != nil {
			fmt.Println("error in reading response body in startreading")
		}
		defer resp.Body.Close()
		bytesRead, _ := ioutil.ReadAll(resp.Body)
		//fmt.Println(string(bytesRead))

		cdes := data.CoreDataEventsFromJson(bytesRead)
		// iterate through cdes to get only "supply" or "consume" device
		if strings.Contains(strings.ToLower(cdes.DataEvents[0].Device), strings.ToLower(deviceType)) {
			updatedCharge := false
			updatedRate := false
			dl := data.DeviceTypeDetails{}
			dl.DeviceAddress = d.PeerId
			dl.DeviceName = cdes.DataEvents[0].Device
			dl.Id = cdes.DataEvents[0].Id
			for _, cde := range cdes.DataEvents {

				for _, reading := range cde.Readings {
					if strings.Contains(reading.Name, "Charge") && updatedCharge == false {
						updatedCharge = true
						//dl.Charge = cdes.DataEvents[0].Readings[0].Value
						dl.Charge = reading.Value
					} else if strings.Contains(reading.Name, "Rate") && updatedRate == false {
						updatedRate = true
						//dl.Rate = cdes.DataEvents[0].Readings[0].Value
						dl.Rate = reading.Value
					} else {
						continue
					}
				}
				if updatedCharge && updatedRate {
					sl = append(sl, dl)
					break
				}
				//if strings.Contains(strings.ToLower(cdes.DataEvents[0].Device), strings.ToLower(deviceType)) {

				//}
			}
		}

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

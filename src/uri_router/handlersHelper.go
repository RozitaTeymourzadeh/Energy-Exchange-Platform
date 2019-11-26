package uri_router

import (
	"bytes"
	"fmt"
	"github.com/edgexfoundry/device-simple/src/data"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
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
				val.PeerId = peer.IpAdd + ":" + peer.Port
				dl.Devices = append(dl.Devices, val)
			}

		}
	}

	return dl
}

func GetSelfDevices() {
	uri := "http://" + data.GetNodeId().EdgeXAddress + ":48082/api/v1/device"

	resp, err := http.Get(uri)
	if err != nil {
		fmt.Println("Error in getting all devices")
	}
	defer resp.Body.Close()
	bytesRead, _ := ioutil.ReadAll(resp.Body)

	deviceList := data.DeviceListFromJson(bytesRead)
	for _, device := range deviceList.Devices {
		device.PeerId = data.GetNodeId().Address + ":" + data.GetNodeId().Port
		SELFDEVICES.Devices[device.Id] = device
	}
}

//1. get list of peers
//	2. iterate and get list of devices for the peer
//		3. for device with name
//			4. if name contains "supply"
//				5. add to supplyDeviceList
//
//use by PageVars
func generateSupplyDeviceTypeBoard(deviceType string) []data.DeviceTypeDetails {

	sl := make([]data.DeviceTypeDetails, 0)

	for _, d := range DEVICELIST.Devices {

		if strings.HasPrefix(d.Name, "Supply") {
			//uri := "http://" + d.PeerId + ":48080/api/v1/event/device/" + d.Name + "/" + "10"
			//uri := "http://" + d.PeerId + ":9999/sendDeviceEvents/" + d.Name + "/" + "10"
			uri := "http://" + d.PeerId + "/sendDeviceEvents/" + d.Name + "/" + "10"
			resp, err := http.Get(uri)

			if err != nil {
				fmt.Println("error in reading response body in startreading")
			}
			defer resp.Body.Close()
			bytesRead, _ := ioutil.ReadAll(resp.Body)
			//fmt.Println(string(bytesRead))

			cdes := data.CoreDataEventsFromJson(bytesRead)
			// iterate through cdes to get only "supply" device
			//fmt.Println("strings.ToLower(cdes.DataEvents[0].Device), strings.ToLower(deviceType) : "+ strings.ToLower(deviceType))
			if strings.Contains(strings.ToLower(cdes.DataEvents[0].Device), strings.ToLower(deviceType)) {
				updatedSupplierCharge := false
				updatedSupplierRate := false
				dl := data.DeviceTypeDetails{}
				dl.DeviceAddress = d.PeerId
				dl.DeviceName = cdes.DataEvents[0].Device
				dl.Id = cdes.DataEvents[0].Id
				for _, cde := range cdes.DataEvents {

					for _, reading := range cde.Readings {
						if strings.Contains(reading.Name, "supplierCharge") && updatedSupplierCharge == false {
							updatedSupplierCharge = true
							//dl.Charge = cdes.DataEvents[0].Readings[0].Value
							dl.SupplierCharge, err = strconv.Atoi(reading.Value)
						} else if strings.Contains(reading.Name, "supplyRate") && updatedSupplierRate == false {
							updatedSupplierRate = true
							//dl.Rate = cdes.DataEvents[0].Readings[0].Value
							dl.SupplyRate, err = strconv.Atoi(reading.Value)
						} else {
							continue
						}
					}
					if updatedSupplierCharge && updatedSupplierRate {
						sl = append(sl, dl)
						break
					}
					//if strings.Contains(strings.ToLower(cdes.DataEvents[0].Device), strings.ToLower(deviceType)) {

					//}
				}
			}

		}
	}
	return sl
}

func generateConsumeDeviceTypeBoard(deviceType string) []data.DeviceTypeDetails {

	sl := make([]data.DeviceTypeDetails, 0)

	for _, d := range DEVICELIST.Devices {

		if strings.HasPrefix(d.Name, "Consume") {
			//uri := "http://" + d.PeerId + ":48080/api/v1/event/device/" + d.Name + "/" + "10"
			//uri := "http://" + d.PeerId + ":9999/sendDeviceEvents/" + d.Name + "/" + "10"
			uri := "http://" + d.PeerId + "/sendDeviceEvents/" + d.Name + "/" + "10"
			resp, err := http.Get(uri)

			if err != nil {
				fmt.Println("error in reading response body in startreading")
			}
			defer resp.Body.Close()
			bytesRead, _ := ioutil.ReadAll(resp.Body)
			//fmt.Println(string(bytesRead))

			cdes := data.CoreDataEventsFromJson(bytesRead)
			// iterate through cdes to get only "supply" device
			//fmt.Println("strings.ToLower(cdes.DataEvents[0].Device), strings.ToLower(deviceType) : "+ strings.ToLower(deviceType))
			if strings.Contains(strings.ToLower(cdes.DataEvents[0].Device), strings.ToLower(deviceType)) {
				updatedConsumerCharge := false
				updatedconsumerDischargeRate := false
				dl := data.DeviceTypeDetails{}
				dl.DeviceAddress = d.PeerId
				dl.DeviceName = cdes.DataEvents[0].Device
				dl.Id = cdes.DataEvents[0].Id
				for _, cde := range cdes.DataEvents {

					for _, reading := range cde.Readings {
						if strings.Contains(reading.Name, "consumerCharge") && updatedConsumerCharge == false {
							updatedConsumerCharge = true
							//dl.Charge = cdes.DataEvents[0].Readings[0].Value
							dl.ConsumerCharge, err = strconv.Atoi(reading.Value)
						} else if strings.Contains(reading.Name, "consumerDischargeRate") && updatedconsumerDischargeRate == false {
							updatedconsumerDischargeRate = true
							//dl.Rate = cdes.DataEvents[0].Readings[0].Value
							dl.ConsumerDischargeRate, err = strconv.Atoi(reading.Value)
						} else {
							continue
						}
					}
					if updatedConsumerCharge && updatedconsumerDischargeRate {
						sl = append(sl, dl)
						break
					}
					//if strings.Contains(strings.ToLower(cdes.DataEvents[0].Device), strings.ToLower(deviceType)) {

					//}
				}
			}

		}
	}
	return sl
}

func updateDeviceTypeBoards() {
	//data.GetSupplyDeviceBoard()
	for _, d := range DEVICELIST.Devices {
		fmt.Println(": : : Device in DEVICELIST.Devices : : :")
		fmt.Println("Device d.PeerId : " + d.PeerId)
		fmt.Println("Device d.Name : " + d.Name)
		fmt.Println("Device d.Id : " + d.Id)

	}
}

func sendTransactionToSupplier(tx data.Transaction) {
	txJson, err := tx.TransactionToJSON()
	if err != nil {
		log.Print("Cannot create transaction")
		return
	}

	//uri := "http://localhost:48081/api/v1/deviceprofile/id/" + vars["deviceId"]
	//uri := "http://" + tx.SupplierAddress + ":9999/suppliertx"
	uri := "http://" + tx.SupplierAddress + "/suppliertx"
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
	//uri := "http://" + tx.ConsumerAddress + ":9999/consumertx"
	uri := "http://" + tx.ConsumerAddress + "/consumertx"
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

//func generateDeviceTypeBoard(deviceType string) []data.DeviceTypeDetails {
//
//	sl := make([]data.DeviceTypeDetails, 0)
//
//	for _, d := range DEVICELIST.Devices {
//
//		//uri := "http://" + d.PeerId + ":48080/api/v1/event/device/" + d.Name + "/" + "10"
//		//uri := "http://" + d.PeerId + ":9999/sendDeviceEvents/" + d.Name + "/" + "10"
//		uri := "http://" + d.PeerId + "/sendDeviceEvents/" + d.Name + "/" + "10"
//		resp, err := http.Get(uri)
//
//		if err != nil {
//			fmt.Println("error in reading response body in startreading")
//		}
//		defer resp.Body.Close()
//		bytesRead, _ := ioutil.ReadAll(resp.Body)
//		//fmt.Println(string(bytesRead))
//
//		cdes := data.CoreDataEventsFromJson(bytesRead)
//		// iterate through cdes to get only "supply" or "consume" device
//		if strings.Contains(strings.ToLower(cdes.DataEvents[0].Device), strings.ToLower(deviceType)) {
//			updatedCharge := false
//			updatedRate := false
//			dl := data.DeviceTypeDetails{}
//			dl.DeviceAddress = d.PeerId
//			dl.DeviceName = cdes.DataEvents[0].Device
//			dl.Id = cdes.DataEvents[0].Id
//			for _, cde := range cdes.DataEvents {
//
//				for _, reading := range cde.Readings {
//					if strings.Contains(reading.Name, "Charge") && updatedCharge == false {
//						updatedCharge = true
//						//dl.Charge = cdes.DataEvents[0].Readings[0].Value
//						dl.Charge = reading.Value
//					} else if strings.Contains(reading.Name, "Rate") && updatedRate == false {
//						updatedRate = true
//						//dl.Rate = cdes.DataEvents[0].Readings[0].Value
//						dl.Rate = reading.Value
//					} else {
//						continue
//					}
//				}
//				if updatedCharge && updatedRate {
//					sl = append(sl, dl)
//					break
//				}
//				//if strings.Contains(strings.ToLower(cdes.DataEvents[0].Device), strings.ToLower(deviceType)) {
//
//				//}
//			}
//		}
//
//	}
//	return sl
//}

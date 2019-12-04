package driver

import (
	"fmt"
	//"github.com/edgexfoundry/device-simple/driver/data"
	"io/ioutil"
	"net/http"
	"strings"
)

func GetSelfDevices() {
	uri := "http://" + GetNodeId().EdgeXAddress + ":48082/api/v1/device"

	resp, err := http.Get(uri)
	if err != nil {
		fmt.Println("Error in getting all devices")
	}
	defer resp.Body.Close()
	bytesRead, _ := ioutil.ReadAll(resp.Body)

	deviceList := DeviceListFromJson(bytesRead)
	for _, device := range deviceList.Devices {
		device.PeerId = GetNodeId().Address + ":" + GetNodeId().Port
		SELFDEVICES.Devices[device.Id] = device
		if strings.HasPrefix(device.Name, "Supply") {
			SetSupplyDeviceName(device.Name)
			SetSupplyDeviceId(device.Id)
			SetSupplyDeviceAddress(GetNodeId().Address + ":" + GetNodeId().Port)
		}
		if strings.HasPrefix(device.Name, "Consume") {
			SetConsumeDeviceName(device.Name)
			SetConsumeDeviceId(device.Id)
			SetConsumeDeviceAddress(GetNodeId().Address + ":" + GetNodeId().Port)
		}
	}
}

//1. get list of peers
//	2. iterate and get list of devices for the peer
//		3. for device with name
//			4. if name contains "supply"
//				5. add to supplyDeviceList
//
//use by PageVars
func generateSupplyDeviceTypeBoard(deviceType string) []DeviceTypeDetails {

	sl := make([]DeviceTypeDetails, 0)
	sd := DeviceTypeDetails{}
	sd.DeviceAddress = GetNodeId().GetAddressPort()
	sd.DeviceName = GetSupplyDeviceName()
	sd.DeviceId = GetSupplyDeviceId()
	sd.SupplierCharge = GetSupplierCharge()
	sd.SupplierChargeRate = GetSupplierChargeRate()
	sd.SupplyRate = GetSupplyRate()
	sd.Surplus = GetSurplus()
	sd.IsSupplying = GetIsSupplying()
	sd.ToSupply = GetToSupply()
	sd.SellRate = GetSellRate()
	sd.HasOffered = GetHasOffered()
	sd.SellThreshold = GetSellThreshold()
	sd.SupplierMaxCharge = GetSupplierMaxCharge()

	sl = append(sl, sd)
	//for _, d := range DEVICELIST.Devices {
	//
	//	if strings.HasPrefix(d.Name, "Supply") {
	//		//uri := "http://" + d.PeerId + ":48080/api/v1/event/device/" + d.Name + "/" + "10"
	//		//uri := "http://" + d.PeerId + ":9999/sendDeviceEvents/" + d.Name + "/" + "10"
	//		uri := "http://" + d.PeerId + "/sendDeviceEvents/" + d.Name + "/" + "10"
	//		fmt.Println("calling for device events : " + uri)
	//		resp, err := http.Get(uri)
	//
	//		if err != nil {
	//			fmt.Println("error in reading response body in startreading")
	//		}
	//		defer resp.Body.Close()
	//		bytesRead, _ := ioutil.ReadAll(resp.Body)
	//		//fmt.Println(string(bytesRead))
	//
	//		cdes := CoreDataEventsFromJson(bytesRead)
	//		// iterate through cdes to get only "supply" device
	//		//fmt.Println("strings.ToLower(cdes.DataEvents[0].Device), strings.ToLower(deviceType) : "+ strings.ToLower(deviceType))
	//		if len(cdes.DataEvents) > 0 && strings.Contains(strings.ToLower(cdes.DataEvents[0].Device), strings.ToLower(deviceType)) {
	//			updatedSupplierCharge := false
	//			updatedSupplierRate := false
	//			dl := DeviceTypeDetails{}
	//			dl.DeviceAddress = d.PeerId
	//			dl.DeviceName = cdes.DataEvents[0].Device
	//			dl.Id = cdes.DataEvents[0].Id
	//			for _, cde := range cdes.DataEvents {
	//
	//				for _, reading := range cde.Readings {
	//					if strings.Contains(reading.Name, "supplierCharge") && updatedSupplierCharge == false {
	//						updatedSupplierCharge = true
	//						//dl.Charge = cdes.DataEvents[0].Readings[0].Value
	//						dl.SupplierCharge, err = strconv.Atoi(reading.Value)
	//					} else if strings.Contains(reading.Name, "supplyRate") && updatedSupplierRate == false {
	//						updatedSupplierRate = true
	//						//dl.Rate = cdes.DataEvents[0].Readings[0].Value
	//						dl.SupplyRate, err = strconv.Atoi(reading.Value)
	//					} else {
	//						continue
	//					}
	//				}
	//				if updatedSupplierCharge && updatedSupplierRate {
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
	//}
	return sl
}

func generateConsumeDeviceTypeBoard(deviceType string) []DeviceTypeDetails {

	sl := make([]DeviceTypeDetails, 0)
	cd := DeviceTypeDetails{}
	cd.DeviceAddress = GetNodeId().GetAddressPort()
	cd.DeviceName = GetConsumeDeviceName()
	cd.DeviceId = GetConsumeDeviceId()
	cd.ConsumerCharge = GetConsumerCharge()
	cd.ConsumerDischargeRate = GetConsumerDischargeRate()
	cd.Require = GetRequire()
	cd.IsReceiving = GetIsReceiving()
	cd.ToReceive = GetToReceive()
	cd.BuyRate = GetBuyRate()
	cd.BuyBaseRate = GetBuyBaseRate()
	cd.ToReceiveRate = GetToReceiveRate()
	cd.HasAsked = GetHasAsked()
	cd.ConsumerMaxCharge = GetConsumerMaxCharge()
	cd.BuyThreshold = GetBuyThreshold()

	sl = append(sl, cd)
	//for _, d := range DEVICELIST.Devices {
	//
	//	if strings.HasPrefix(d.Name, "Consume") {
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
	//		//fmt.Println("CoreDataEventsFromJson : !!! ::: " + string(bytesRead))
	//
	//		cdes := CoreDataEventsFromJson(bytesRead)
	//		if len(cdes.DataEvents) > 0 && strings.Contains(strings.ToLower(cdes.DataEvents[0].Device), strings.ToLower(deviceType)) {
	//			updatedConsumerCharge := false
	//			updatedconsumerDischargeRate := false
	//			dl := DeviceTypeDetails{}
	//			dl.DeviceAddress = d.PeerId
	//			dl.DeviceName = cdes.DataEvents[0].Device
	//			dl.Id = cdes.DataEvents[0].Id
	//			for _, cde := range cdes.DataEvents {
	//
	//				for _, reading := range cde.Readings {
	//					if strings.Contains(reading.Name, "consumerCharge") && updatedConsumerCharge == false {
	//						updatedConsumerCharge = true
	//						//dl.Charge = cdes.DataEvents[0].Readings[0].Value
	//						dl.ConsumerCharge, err = strconv.Atoi(reading.Value)
	//					} else if strings.Contains(reading.Name, "consumerDischargeRate") && updatedconsumerDischargeRate == false {
	//						updatedconsumerDischargeRate = true
	//						//dl.Rate = cdes.DataEvents[0].Readings[0].Value
	//						dl.ConsumerDischargeRate, err = strconv.Atoi(reading.Value)
	//					} else {
	//						continue
	//					}
	//				}
	//				if updatedConsumerCharge && updatedconsumerDischargeRate {
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
	//}
	return sl
}

//func sendTransactionToSupplier(tx Transaction) {
//	txJson, err := tx.TransactionToJSON()
//	if err != nil {
//		log.Print("Cannot create transaction")
//		return
//	}
//
//	//uri := "http://localhost:48081/api/v1/deviceprofile/id/" + vars["deviceId"]
//	//uri := "http://" + tx.SupplierAddress + ":9999/suppliertx"
//	uri := "http://" + tx.SupplierAddress + "/suppliertx"
//	log.Println("sending post req to : " + uri)
//	client := &http.Client{}
//	// creating request
//	req, _ := http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(txJson))
//	// fetching response
//	resp, err := client.Do(req)
//	if err != nil {
//		fmt.Println("error in reading response body in start reading")
//	}
//	defer resp.Body.Close()
//}

//func sendTransactionToConsumer(tx Transaction) {
//	txJson, err := tx.TransactionToJSON()
//	if err != nil {
//		log.Print("Cannot create transaction")
//		return
//	}
//
//	//uri := "http://localhost:48081/api/v1/deviceprofile/id/" + vars["deviceId"]
//	//uri := "http://" + tx.ConsumerAddress + ":9999/consumertx"
//	uri := "http://" + tx.ConsumerAddress + "/consumertx"
//	log.Println("sending post req to : " + uri)
//	client := &http.Client{}
//	// creating request
//	req, _ := http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(txJson))
//	// fetching response
//	resp, err := client.Do(req)
//	if err != nil {
//		fmt.Println("error in reading response body in start reading")
//	}
//	defer resp.Body.Close()
//}

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

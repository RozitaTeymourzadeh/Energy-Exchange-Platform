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
			SetSupplyDeviceAddress(GetNodeId().GetAddressPort())
		}
		if strings.HasPrefix(device.Name, "Consume") {
			SetConsumeDeviceName(device.Name)
			SetConsumeDeviceId(device.Id)
			SetConsumeDeviceAddress(GetNodeId().GetAddressPort())
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
	return sl
}

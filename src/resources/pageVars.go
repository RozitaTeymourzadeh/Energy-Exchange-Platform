package resources

import "github.com/edgexfoundry/device-simple/src/data"

type PageVars struct {
	Title                 string
	DeviceList            []data.Device
	SupplyDevicesDetails  []data.SupplyDeviceDetails
	ConsumeDevicesDetails []data.ConsumeDeviceDetails
	//
	// transaction history :
	Body string
}

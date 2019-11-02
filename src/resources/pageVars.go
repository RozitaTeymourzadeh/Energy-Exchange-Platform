package resources

import "github.com/edgexfoundry/device-simple/src/data"

type PageVars struct {
	Title                 string
	DeviceList            []data.Device
	SupplyDevicesDetails  []data.DeviceTypeDetails
	ConsumeDevicesDetails []data.DeviceTypeDetails
	Transactions          []data.Transaction
	// transaction history :
	Body string
}

func NewPageVars() PageVars {
	pv := PageVars{
		Title:                 "",
		DeviceList:            make([]data.Device, 0),
		SupplyDevicesDetails:  make([]data.DeviceTypeDetails, 0),
		ConsumeDevicesDetails: make([]data.DeviceTypeDetails, 0),
		Body:                  "",
	}
	return pv
}

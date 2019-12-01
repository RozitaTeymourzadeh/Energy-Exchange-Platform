package driver

type PageVars struct {
	Title                 string
	DeviceMap             []Device
	SupplyDevicesDetails  []DeviceTypeDetails
	ConsumeDevicesDetails []DeviceTypeDetails
	Transactions          []Transaction

	// transaction history :
	Body string
}

func NewPageVars() PageVars {
	pv := PageVars{
		Title:                 "EEP",
		DeviceMap:             make([]Device, 0),
		SupplyDevicesDetails:  make([]DeviceTypeDetails, 0),
		ConsumeDevicesDetails: make([]DeviceTypeDetails, 0),
		Body:                  "",
	}
	return pv
}

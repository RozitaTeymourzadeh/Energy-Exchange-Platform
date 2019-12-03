package driver

type PageVars struct {
	Title                 string
	IpPort                string
	DeviceMap             []Device
	SupplyDevicesDetails  []DeviceTypeDetails
	ConsumeDevicesDetails []DeviceTypeDetails
	Transactions          []Transaction
	SdReadings            []int
	CdReadings            []int

	// transaction history :
	Body string
}

func NewPageVars() PageVars {
	pv := PageVars{
		Title:                 "EEP",
		IpPort:                SELF_ADDR,
		DeviceMap:             make([]Device, 0),
		SupplyDevicesDetails:  make([]DeviceTypeDetails, 0),
		ConsumeDevicesDetails: make([]DeviceTypeDetails, 0),
		Body:                  "",
	}
	return pv
}

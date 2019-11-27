package driver

//being used in pageVars.go in resources
type DeviceTypeDetails struct {
	DeviceAddress string // ip and port
	DeviceName    string
	Id            string
	//Reading
	//supplier//
	SupplierCharge     int
	SupplierChargeRate int
	SupplyRate         int
	Surplus            int
	IsSupplying        int
	ToSupply           int
	SellRate           int
	SellBaseRate       int
	HasOffered         bool
	SupplierMaxCharge  int
	SellThreshold      int
	//consumer //
	ConsumerCharge        int
	ConsumerDischargeRate int
	Require               int
	IsReceiving           int
	ToReceive             int
	BuyRate               int
	BuyBaseRate           int
	ToReceiveRate         int
	HasAsked              bool
	ConsumerMaxCharge     int
	BuyThreshold          int
}

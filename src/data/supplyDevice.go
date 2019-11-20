package data

import (
	"sync"
)

type SupplyDevice struct {
	Charge             int
	SupplyRate         int
	SupplierChargeRate int

	IsSupplying int
	ToSupply    int

	SellRate int
}

var supplyDevice *SupplyDevice
var sonce sync.Once

func GetSupplyDevice() *SupplyDevice {
	sonce.Do(func() {
		supplyDevice = &SupplyDevice{}
		supplyDevice.SupplierChargeRate = 5
	})
	return supplyDevice
}

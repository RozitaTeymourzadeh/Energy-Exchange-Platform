package driver

import (
	"fmt"
	"math/rand"
	"sync"
)

type SupplyDevice struct {

	// todo add
	SupplyDeviceName    string
	SupplyDeviceId      string
	SupplyDeviceAddress string
	//
	supplierCharge     int
	supplierChargeRate int
	supplyRate         int
	surplus            int
	isSupplying        int
	toSupply           int
	sellRate           int
	sellBaseRate       int

	//toSupplyAddress string
	hasOffered    bool
	maxCharge     int
	sellThreshold int

	mux sync.RWMutex

	last100SDReadings    []int
	last100SDReadingsMux sync.RWMutex
}

var supplyDevice *SupplyDevice
var sonce sync.Once

func GetSupplyDevice() *SupplyDevice {
	sonce.Do(func() {
		fmt.Println("Init SupplyDevice")
		supplyDevice = &SupplyDevice{
			last100SDReadings: []int{},
		}

		supplyDevice.maxCharge = 1000
		supplyDevice.sellThreshold = supplyDevice.maxCharge / 4
		supplyDevice.hasOffered = false
		supplyDevice.supplierCharge = supplyDevice.maxCharge/2 + rand.Intn(supplyDevice.maxCharge/2)
		supplyDevice.supplierChargeRate = 13 + rand.Intn(5)
		supplyDevice.supplyRate = 11 + rand.Intn(5)
		supplyDevice.surplus = 0
		supplyDevice.isSupplying = 0
		supplyDevice.toSupply = 0
		supplyDevice.sellRate = 10
		supplyDevice.sellBaseRate = 10

	})
	return supplyDevice
}

func GetSupplyDeviceName() string {
	return supplyDevice.SupplyDeviceName
}

func GetSupplyDeviceId() string {
	return supplyDevice.SupplyDeviceId
}

func GetSupplyDeviceAddress() string {
	return supplyDevice.SupplyDeviceAddress
}

func GetSupplierCharge() int {
	return supplyDevice.supplierCharge
}

func GetSupplierChargeRate() int {
	return supplyDevice.supplierChargeRate
}

func GetSupplyRate() int {
	return supplyDevice.supplyRate
}

func GetSurplus() int {
	return supplyDevice.surplus
}

func GetIsSupplying() int {
	return supplyDevice.isSupplying
}

func GetToSupply() int {
	return supplyDevice.toSupply
}

//func GetToSupplyAddress() string {
//	return supplyDevice.toSupplyAddress
//}

func GetSellRate() int {
	return supplyDevice.sellRate
}

func GetSellBaseRate() int {
	return supplyDevice.sellBaseRate
}

func GetHasOffered() bool {
	return supplyDevice.hasOffered
}

func GetSupplierMaxCharge() int {
	return supplyDevice.maxCharge
}

func GetSellThreshold() int {
	return supplyDevice.sellThreshold
}

func GetLast100SDReadings() []int {
	supplyDevice.last100SDReadingsMux.Lock()
	defer supplyDevice.last100SDReadingsMux.Unlock()
	return supplyDevice.last100SDReadings
}

func SetSupplyDeviceName(change string) {
	supplyDevice.mux.Lock()
	defer supplyDevice.mux.Unlock()
	supplyDevice.SupplyDeviceName = change
}
func SetSupplyDeviceId(change string) {
	supplyDevice.mux.Lock()
	defer supplyDevice.mux.Unlock()
	supplyDevice.SupplyDeviceId = change
}
func SetSupplyDeviceAddress(change string) {
	supplyDevice.mux.Lock()
	defer supplyDevice.mux.Unlock()
	supplyDevice.SupplyDeviceAddress = change
}

func SetSupplierCharge(change int) {
	supplyDevice.mux.Lock()
	defer supplyDevice.mux.Unlock()

	if change < 0 {
		supplyDevice.supplierCharge = 0
	} else if change > supplyDevice.maxCharge {
		supplyDevice.supplierCharge = supplyDevice.maxCharge
	} else {
		supplyDevice.supplierCharge = change
	}
}

func SetSupplierChargeRate(change int) {
	supplyDevice.mux.Lock()
	defer supplyDevice.mux.Unlock()
	supplyDevice.supplierChargeRate = change
}

func SetSupplyRate(change int) {
	supplyDevice.mux.Lock()
	defer supplyDevice.mux.Unlock()
	supplyDevice.supplyRate = change
}

func SetSurplus(change int) {
	supplyDevice.mux.Lock()
	defer supplyDevice.mux.Unlock()
	supplyDevice.surplus = change
}

func SetIsSupplying(change int) {
	supplyDevice.mux.Lock()
	defer supplyDevice.mux.Unlock()
	supplyDevice.isSupplying = change
}

func SetToSupply(change int) {
	supplyDevice.mux.Lock()
	defer supplyDevice.mux.Unlock()
	supplyDevice.toSupply = change
}

func SetSellRate(change int) {
	supplyDevice.mux.Lock()
	defer supplyDevice.mux.Unlock()
	supplyDevice.sellRate = change
}

//func SetToSupplyAddress(change string) string {
//	supplyDevice.mux.Lock()
//	defer supplyDevice.mux.Unlock()
//	supplyDevice.toSupplyAddress = change
//}

func SetHasOffered(change bool) {
	supplyDevice.mux.Lock()
	defer supplyDevice.mux.Unlock()
	supplyDevice.hasOffered = change
}

func SetSupplierMaxCharge(change int) {
	supplyDevice.mux.Lock()
	defer supplyDevice.mux.Unlock()
	supplyDevice.maxCharge = change
}

func AddToLast100SDReadings(value int) {
	supplyDevice.last100SDReadingsMux.Lock()
	defer supplyDevice.last100SDReadingsMux.Unlock()
	supplyDevice.last100SDReadings = append([]int{value}, supplyDevice.last100SDReadings...)
	if len(supplyDevice.last100SDReadings) > 100 {
		supplyDevice.last100SDReadings = supplyDevice.last100SDReadings[:100]
	}
}

func SupplyCompleteCleanup() {
	supplyDevice.mux.Lock()
	defer supplyDevice.mux.Unlock()
	supplyDevice.isSupplying = 0
	supplyDevice.toSupply = 0
	supplyDevice.hasOffered = false
	//supplyDevice.toSupplyAddress = ""
}

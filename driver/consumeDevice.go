package driver

import (
	"fmt"
	"math/rand"
	"sync"
)

type ConsumeDevice struct {
	consumerCharge        int
	consumerDischargeRate int
	require               int
	isReceiving           int
	toReceive             int
	buyRate               int
	buyBaseRate           int

	toReceiveRate int
	hasAsked      bool
	maxCharge     int
	buyThreshold  int

	mux sync.RWMutex
}

var consumeDevice *ConsumeDevice
var conce sync.Once

func GetConsumeDevice() *ConsumeDevice {
	conce.Do(func() {
		fmt.Println("Init ConsumeDevice")
		consumeDevice = &ConsumeDevice{}

		consumeDevice.maxCharge = 1000
		consumeDevice.buyThreshold = consumeDevice.maxCharge - 200 //todo : change
		consumeDevice.hasAsked = false
		consumeDevice.consumerCharge = consumeDevice.maxCharge/2 + rand.Intn(consumeDevice.maxCharge/2)
		consumeDevice.consumerDischargeRate = 5 + rand.Intn(5)
		consumeDevice.require = 0
		consumeDevice.isReceiving = 0
		consumeDevice.toReceive = 0
		consumeDevice.buyRate = 10
		consumeDevice.buyBaseRate = 10
		consumeDevice.toReceiveRate = 0

	})
	return consumeDevice
}

func GetConsumerCharge() int {
	return consumeDevice.consumerCharge
}

func GetConsumerDischargeRate() int {
	return consumeDevice.consumerDischargeRate
}

func GetRequire() int {
	return consumeDevice.require
}

func GetIsReceiving() int {
	return consumeDevice.isReceiving
}

func GetToReceive() int {
	return consumeDevice.toReceive
}

func GetBuyRate() int {
	return consumeDevice.buyRate
}

func GetBuyBaseRate() int {
	return consumeDevice.buyBaseRate
}

func GetToReceiveRate() int {
	return consumeDevice.toReceiveRate
}

func GetHasAsked() bool {
	return consumeDevice.hasAsked
}

func GetConsumerMaxCharge() int {
	return consumeDevice.maxCharge
}

func GetBuyThreshold() int {
	return consumeDevice.buyThreshold
}

func SetConsumerCharge(change int) {
	consumeDevice.mux.Lock()
	defer consumeDevice.mux.Unlock()

	if change < 0 {
		consumeDevice.consumerCharge = 0
	} else if change > 1000 {
		consumeDevice.consumerCharge = 1000
	} else {
		consumeDevice.consumerCharge = change
	}
}

func SetConsumerDischargeRate(change int) {
	consumeDevice.mux.Lock()
	defer consumeDevice.mux.Unlock()
	consumeDevice.consumerDischargeRate = change
}

func SetRequire(change int) {
	consumeDevice.mux.Lock()
	defer consumeDevice.mux.Unlock()
	consumeDevice.require = change
}

func SetIsReceiving(change int) {
	consumeDevice.mux.Lock()
	defer consumeDevice.mux.Unlock()
	consumeDevice.isReceiving = change
}

func SetToReceive(change int) {
	consumeDevice.mux.Lock()
	defer consumeDevice.mux.Unlock()
	consumeDevice.toReceive = change
}

func SetBuyRate(change int) {
	consumeDevice.mux.Lock()
	defer consumeDevice.mux.Unlock()
	consumeDevice.buyRate = change
}

func SetHasAsked(change bool) {
	consumeDevice.mux.Lock()
	defer consumeDevice.mux.Unlock()
	consumeDevice.hasAsked = change
}

func SetConsumerMaxCharge(change int) {
	consumeDevice.mux.Lock()
	defer consumeDevice.mux.Unlock()
	consumeDevice.maxCharge = change
}

func SetToReceiveRate(change int) {
	consumeDevice.mux.Lock()
	defer consumeDevice.mux.Unlock()
	consumeDevice.toReceiveRate = change
}

func ConsumeCompleteCleanup() {
	consumeDevice.mux.Lock()
	defer consumeDevice.mux.Unlock()
	consumeDevice.isReceiving = 0
	consumeDevice.toReceive = 0
	consumeDevice.toReceiveRate = 0
	consumeDevice.hasAsked = false

}

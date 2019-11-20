package data

import (
	"sync"
)

type ConsumeDevice struct {
	Charge      int
	ConsumeRate int
	Require     int

	IsReceiving int
	ToReceive   int

	BuyRate int
}

var consumeDevice *ConsumeDevice
var conce sync.Once

func GetConsumeDevice() *ConsumeDevice {
	conce.Do(func() {
		consumeDevice = &ConsumeDevice{}
	})
	return consumeDevice
}

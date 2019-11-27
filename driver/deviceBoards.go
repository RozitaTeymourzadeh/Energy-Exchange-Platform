package driver

import (
	"fmt"
	"sync"
)

type DeviceBoards struct {
	supplyBoard  map[string]DeviceTypeDetails
	consumeBoard map[string]DeviceTypeDetails

	sbMux sync.RWMutex
	cbMux sync.RWMutex
}

var deviceBoards *DeviceBoards
var dbonce sync.Once

func GetDeviceBoards() *DeviceBoards {
	conce.Do(func() {
		fmt.Println("Init DeviceBoards")
		deviceBoards = &DeviceBoards{
			supplyBoard:  make(map[string]DeviceTypeDetails),
			consumeBoard: make(map[string]DeviceTypeDetails),
			sbMux:        sync.RWMutex{},
			cbMux:        sync.RWMutex{},
		}

	})
	return deviceBoards
}

func GetSupplyDeviceBoard() map[string]DeviceTypeDetails {
	deviceBoards.sbMux.RLock()
	defer deviceBoards.sbMux.RUnlock()
	return GetDeviceBoards().supplyBoard
}

func GetConsumeDeviceBoard() map[string]DeviceTypeDetails {
	deviceBoards.cbMux.RLock()
	defer deviceBoards.cbMux.RUnlock()
	return GetDeviceBoards().consumeBoard
}

func GetFromSupplyBoard(key string) (DeviceTypeDetails, bool) {
	deviceBoards.sbMux.RLock()
	defer deviceBoards.sbMux.RUnlock()
	dtd, ok := deviceBoards.supplyBoard[key]
	return dtd, ok
}

func InsertInSupplyBoard(key string, dtd DeviceTypeDetails) {
	deviceBoards.sbMux.Lock()
	defer deviceBoards.sbMux.Unlock()

	deviceBoards.supplyBoard[key] = dtd

}

func DeleteInSupplyBoard(key string) {
	deviceBoards.sbMux.Lock()
	defer deviceBoards.sbMux.Unlock()

	delete(deviceBoards.supplyBoard, key)

}

func GetFromConsumeBoard(key string) (DeviceTypeDetails, bool) {
	deviceBoards.cbMux.RLock()
	defer deviceBoards.cbMux.RUnlock()
	dtd, ok := deviceBoards.consumeBoard[key]
	return dtd, ok
}

func InsertInConsumeBoard(key string, dtd DeviceTypeDetails) {
	deviceBoards.cbMux.Lock()
	defer deviceBoards.cbMux.Unlock()

	deviceBoards.consumeBoard[key] = dtd

}

func DeleteInConsumeBoard(key string) {
	deviceBoards.cbMux.Lock()
	defer deviceBoards.cbMux.Unlock()

	delete(deviceBoards.consumeBoard, key)

}

package uri_router

import (
	"fmt"
	"github.com/edgexfoundry/device-simple/src/data"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func MakeDecision() {
	// making decision
	for true {

		///
		DEVICELIST = getAllDevices( /*data.GetNodeId().ConnectingAddress*/ )
		SUPPLYDEVICEDETAILS = generateDeviceTypeBoard("supply")
		CONSUMEDEVICEDETAILS = generateDeviceTypeBoard("consume")
		///

		fmt.Println("Decision : ")
		if len(CONSUMEDEVICEDETAILS) < 1 {
			//_, _ = w.Write([]byte("No consume device"))
			fmt.Println("No consume device")

		} else if len(SUPPLYDEVICEDETAILS) < 1 {
			//_, _ = w.Write([]byte("No supply device"))
			fmt.Println("No supply device")

		} else {

			str := makeDecisionHandlerHelper()
			//_, _ = w.Write([]byte(str))
			fmt.Println(str)

		} // end of else

		time.Sleep(10 * time.Second) // every 10 seconds
	}

}

func makeDecisionHandlerHelper() string {
	sb := strings.Builder{}
	sb.WriteString("Pairing consume and supply devices:\n")
	for _, cv := range CONSUMEDEVICEDETAILS {
		matched := false
		sb.WriteString(">>> \n")
		for _, sv := range SUPPLYDEVICEDETAILS {
			sval, _ := strconv.Atoi(sv.Charge)
			cval, _ := strconv.Atoi(cv.Charge)
			//if sval >= cval { // one supply device supplying all energy needed by the consume device
			if cval <= 40 && sval >= 40 {
				matched = true
				//SupplyDeviceDetails[sk] = string(sval - cval)

				//generate random number between 15 and 30
				rand.Seed(time.Now().UnixNano())
				min := 10
				max := 30
				//rand.Intn(max - min + 1) + min)
				randPowerUnits := rand.Intn(max-min+1) + min
				newTx := data.Transaction{
					SupplierName:    sv.DeviceName,
					SupplierId:      sv.Id,
					SupplierAddress: sv.DeviceAddress,
					ConsumerName:    cv.DeviceName,
					ConsumerId:      cv.Id,
					ConsumerAddress: cv.DeviceAddress,
					PowerUnits:      strconv.Itoa(randPowerUnits),
				}
				TRANSACTIONS = append([]data.Transaction{newTx}, TRANSACTIONS...) // prepend
				sb.WriteString(cv.DeviceName + " will receive " + strconv.Itoa(randPowerUnits) + " units from " + sv.DeviceName)

				go sendTransactionToSupplier(newTx) // spawning new thread

			}
		}
		if matched == false {
			sb.WriteString("Could not match " + cv.DeviceName + " to any Supply device")
		}
	}

	return sb.String()
}

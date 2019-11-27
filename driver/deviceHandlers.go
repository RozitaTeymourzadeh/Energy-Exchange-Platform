package driver

import (
	"errors"
	"fmt"
	//"github.com/edgexfoundry/device-simple/driver/data"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

//////////////////////////
//////////////////////////
//  Device Manager apis
/////////////////////////
/////////////////////////

//var TASKMANAGER_ADDRESS = "http://localhost:6686"

func Periodic() {
	fmt.Println("Periodic")
	time.Sleep(1 * time.Second)
	for true {
		makeDecision()
		time.Sleep(10 * time.Second)
	}
}

func makeDecision() {
	DEVICELIST = getAllDevices( /*data.GetNodeId().ConnectingAddress*/ )
	SUPPLYDEVICEDETAILS = generateSupplyDeviceTypeBoard("supply")
	CONSUMEDEVICEDETAILS = generateConsumeDeviceTypeBoard("consume")

	fmt.Println("Decision : ")
	if len(CONSUMEDEVICEDETAILS) < 1 {
		//_, _ = w.Write([]byte("No consume device"))
		fmt.Println("No consume device")

	} else if len(SUPPLYDEVICEDETAILS) < 1 {
		//_, _ = w.Write([]byte("No supply device"))
		fmt.Println("No supply device")

	} else {

		//str := "Dummy makeDecisionHandlerHelper : todo : makeDecisionHandlerHelper()" //todo : makeDecisionHandlerHelper()
		str := makeDecisionHandlerHelper()
		fmt.Println(str)

	} // end of else

}

func getAllDevices() DeviceList {
	dl := NewDeviceList()
	// start of get self devices ////
	devices := updateDeviceListWithSelfDevices()
	for _, device := range devices.Devices {
		device.PeerId = GetNodeId().Address + ":" + GetNodeId().Port
		dl.Devices = append(dl.Devices, device) /// to read SBC and create board
	}
	// end of get self devices ////
	if len(GetNodeId().GetPeers()) > 0 {
		for _, peer := range GetNodeId().GetPeers() {
			uri := "http://" + peer.IpAdd + ":" + peer.Port + "/sendDeviceList"
			fmt.Println("Sending device req to : ", uri)
			resp, err := http.Get(uri)
			if err != nil {
				fmt.Println("Error in getting all devices")
			} else {
				defer resp.Body.Close()
				bytesRead, _ := ioutil.ReadAll(resp.Body)
				peerDeviceList := DeviceListFromJson(bytesRead)
				for _, val := range peerDeviceList.Devices {
					val.PeerId = peer.IpAdd + ":" + peer.Port
					dl.Devices = append(dl.Devices, val)
				}
			}

		}
	}

	return dl
}

func updateDeviceListWithSelfDevices() DeviceList { // todo : read SBC and make board
	//data.GetSupplyDeviceBoard()
	//todo : update by reading canonical SBC
	resp, err := http.Get("http://" + GetNodeId().Address + ":" + GetNodeId().Port + "/" + "getallselfdevices")
	if err != nil {
		fmt.Println("Error in getting all devices : in : updateDeviceTypeBoards")
	}
	defer resp.Body.Close()
	bytesRead, _ := ioutil.ReadAll(resp.Body)

	deviceList := DeviceListFromJson(bytesRead)
	for _, device := range deviceList.Devices {
		fmt.Println("In updateDeviceTypeBoards : " + device.PeerId + " - " + device.Id + " - " + device.Name)
	}
	return deviceList

}

func SendDeviceList(w http.ResponseWriter, r *http.Request) {
	//uri := "http://localhost:48082/api/v1/device"
	uri := "http://" + GetNodeId().EdgeXAddress + ":" + "48082" + "/api/v1/device"

	resp, err := http.Get(uri)
	if err != nil {
		fmt.Println("Error in getting all devices")
	}
	defer resp.Body.Close()
	bytesRead, _ := ioutil.ReadAll(resp.Body)

	//deviceList := data.DeviceListFromJson(bytesRead)
	//DEVICELIST = deviceList //deviceList.AddAllToDevices(&Devices)

	w.WriteHeader(200)
	_, err = w.Write(bytesRead)
	if err != nil {
		log.Println("Error in getting devices")
		w.WriteHeader(405)
		_, _ = w.Write([]byte("No Device found"))
	}

}

func SendDeviceEvents(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	uri := "http://localhost:48080/api/v1/event/device/" + vars["deviceName"] + "/" + vars["noOfReadings"]

	resp, err := http.Get(uri)
	if err != nil {
		fmt.Println("error in reading response body in startreading")
	}
	defer resp.Body.Close()
	bytesRead, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(bytesRead))

	w.WriteHeader(http.StatusOK)
	w.Write(bytesRead)

}

func SupplierTx(w http.ResponseWriter, r *http.Request) {
	log.Println("Supplier tx recv'ed")

	bytesRead, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(errors.New("error in reading req body of SupplierTx"))
	}
	defer r.Body.Close()

	tx := TransactionFromJSON(bytesRead)

	//pwd, _ := os.Getwd()
	//fmt.Println("Current working dir is : " + pwd)

	changeValue, err := strconv.Atoi(tx.PowerUnits)
	if err != nil {
		log.Println(errors.New("Cannot read Change value in param: " + tx.PowerUnits))
	}

	//parser.UpdateValueInFile("../../cmd/device-simple/supplierChargeValue.txt", -changeValue)
	newVal := GetSupplierCharge() - changeValue
	SetSupplierCharge(newVal)

	sendTransactionToConsumer(tx)
}

func ConsumerTx(w http.ResponseWriter, r *http.Request) {
	log.Println("Consumer tx recv'ed")

	bytesRead, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(errors.New("error in reading req body of ConsumerTx"))
	}
	defer r.Body.Close()

	tx := TransactionFromJSON(bytesRead)

	pwd, _ := os.Getwd()
	fmt.Println("Current working dir is : " + pwd)

	changeValue, err := strconv.Atoi(tx.PowerUnits)
	if err != nil {
		log.Println(errors.New("Cannot read Change value in param: " + tx.PowerUnits))
	}

	//parser.UpdateValueInFile("../../cmd/device-simple/consumerChargeValue.txt", changeValue)
	newVal := GetConsumerCharge() + changeValue
	SetConsumerCharge(newVal)
}

//// moved from handlers.go

func GetAllSelfDevices(w http.ResponseWriter, r *http.Request) {
	//uri := "http://localhost:48082/api/v1/device"
	uri := "http://" + GetNodeId().EdgeXAddress + ":48082/api/v1/device"

	resp, err := http.Get(uri)
	if err != nil {
		fmt.Println("Error in getting all devices")
	}
	defer resp.Body.Close()
	bytesRead, _ := ioutil.ReadAll(resp.Body)

	w.WriteHeader(http.StatusOK)
	w.Write(bytesRead)
	//deviceList := data.DeviceListFromJson(bytesRead)
	//
	//w.Write([]byte(deviceList.ShowDeviceInList()))

}

func DeleteDevice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	//uri := "http://localhost:48081/api/v1/device/id/" + vars["deviceId"]
	uri := "http://" + GetNodeId().EdgeXAddress + ":48081/api/v1/device/id/" + vars["deviceId"]

	//creating a new client
	client := http.Client{}
	// creating request
	req, _ := http.NewRequest(http.MethodDelete, uri, nil)
	// fetching response
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error in reading response body in startreading")
	}
	defer resp.Body.Close()

	bytesRead, _ := ioutil.ReadAll(resp.Body)

	w.Write(bytesRead)
}

func GetAllDeviceProfiles(w http.ResponseWriter, r *http.Request) {
	//uri := "http://localhost:48081/api/v1/deviceprofile"
	uri := "http://" + GetNodeId().EdgeXAddress + ":48081/api/v1/deviceprofile"

	resp, err := http.Get(uri)
	if err != nil {
		fmt.Println("Error in getting all devices")
	}
	defer resp.Body.Close()
	bytesRead, _ := ioutil.ReadAll(resp.Body)

	deviceProfiles := DeviceProfilesFromJson(bytesRead)

	w.Write([]byte(deviceProfiles.ShowDeviceProfiles()))

}

func DeleteDeviceProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	//uri := "http://localhost:48081/api/v1/deviceprofile/id/" + vars["deviceId"]
	uri := "http://" + GetNodeId().EdgeXAddress + ":48081/api/v1/deviceprofile/id/" + vars["deviceId"]

	//creating a new client
	client := http.Client{}
	// creating request
	req, _ := http.NewRequest(http.MethodDelete, uri, nil)
	// fetching response
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error in reading response body in startreading")
	}
	defer resp.Body.Close()

	bytesRead, _ := ioutil.ReadAll(resp.Body)

	w.Write(bytesRead)
}

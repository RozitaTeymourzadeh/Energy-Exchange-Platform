package uri_router

import (
	"errors"
	"fmt"
	"github.com/edgexfoundry/device-simple/src/data"
	"github.com/edgexfoundry/device-simple/src/parser"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

//////////////////////////
//////////////////////////
//  Device Manager apis
/////////////////////////
/////////////////////////

//var TASKMANAGER_ADDRESS = "http://localhost:6686"

func SendDeviceList(w http.ResponseWriter, r *http.Request) {
	//uri := "http://localhost:48082/api/v1/device"
	uri := "http://" + data.GetNodeId().EdgeXAddress + ":" + "48082" + "/api/v1/device"

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

	tx := data.TransactionFromJSON(bytesRead)

	//pwd, _ := os.Getwd()
	//fmt.Println("Current working dir is : " + pwd)

	changeValue, err := strconv.Atoi(tx.PowerUnits)
	if err != nil {
		log.Println(errors.New("Cannot read Change value in param: " + tx.PowerUnits))
	}

	parser.UpdateValueInFile("../../cmd/device-simple/supplierChargeValue.txt", -changeValue)

	sendTransactionToConsumer(tx)
}

func ConsumerTx(w http.ResponseWriter, r *http.Request) {
	log.Println("Consumer tx recv'ed")

	bytesRead, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(errors.New("error in reading req body of ConsumerTx"))
	}
	defer r.Body.Close()

	tx := data.TransactionFromJSON(bytesRead)

	pwd, _ := os.Getwd()
	fmt.Println("Current working dir is : " + pwd)

	changeValue, err := strconv.Atoi(tx.PowerUnits)
	if err != nil {
		log.Println(errors.New("Cannot read Change value in param: " + tx.PowerUnits))
	}

	parser.UpdateValueInFile("../../cmd/device-simple/consumerChargeValue.txt", changeValue)
}

//// moved from handlers.go

func GetAllDevices(w http.ResponseWriter, r *http.Request) {
	//uri := "http://localhost:48082/api/v1/device"
	uri := "http://" + data.GetNodeId().EdgeXAddress + ":48082/api/v1/device"

	resp, err := http.Get(uri)
	if err != nil {
		fmt.Println("Error in getting all devices")
	}
	defer resp.Body.Close()
	bytesRead, _ := ioutil.ReadAll(resp.Body)

	deviceList := data.DeviceListFromJson(bytesRead)
	DEVICELIST = deviceList
	//deviceList.AddAllToDevices(&Devices)

	w.Write([]byte(deviceList.ShowDeviceInList()))

}

func DeleteDevice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	//uri := "http://localhost:48081/api/v1/device/id/" + vars["deviceId"]
	uri := "http://" + data.GetNodeId().EdgeXAddress + ":48081/api/v1/device/id/" + vars["deviceId"]

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
	uri := "http://" + data.GetNodeId().EdgeXAddress + ":48081/api/v1/deviceprofile"

	resp, err := http.Get(uri)
	if err != nil {
		fmt.Println("Error in getting all devices")
	}
	defer resp.Body.Close()
	bytesRead, _ := ioutil.ReadAll(resp.Body)

	deviceProfiles := data.DeviceProfilesFromJson(bytesRead)

	w.Write([]byte(deviceProfiles.ShowDeviceProfiles()))

}

func DeleteDeviceProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	//uri := "http://localhost:48081/api/v1/deviceprofile/id/" + vars["deviceId"]
	uri := "http://" + data.GetNodeId().EdgeXAddress + ":48081/api/v1/deviceprofile/id/" + vars["deviceId"]

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

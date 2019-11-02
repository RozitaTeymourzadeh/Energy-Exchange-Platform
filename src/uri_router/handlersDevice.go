package uri_router

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/edgexfoundry/device-simple/src/data"
	"github.com/edgexfoundry/device-simple/src/parser"
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
func On(w http.ResponseWriter, r *http.Request) {

	rInfo := data.RegistryInfo{
		IpAdd: "localhost", // todo - ngrok ip
		Port:  "9999",
	}

	rJson := rInfo.RegistryToJSON()

	//creating a new client
	client := http.Client{}
	// creating request
	req, _ := http.NewRequest(http.MethodPost, "http://localhost:6686/register" /*ngrok*/, bytes.NewBuffer(rJson)) //todo - ngrok ip
	// fetching response
	_, err := client.Do(req)
	if err != nil {
		log.Println(errors.New("Error in device registration : " + err.Error()))
	}

	w.Write([]byte("On--ing device"))
}

func SendDeviceList(w http.ResponseWriter, r *http.Request) {
	uri := "http://localhost:48082/api/v1/device"

	resp, err := http.Get(uri)
	if err != nil {
		fmt.Println("Error in getting all devices")
	}
	defer resp.Body.Close()
	bytesRead, _ := ioutil.ReadAll(resp.Body)

	//deviceList := data.DeviceListFromJson(bytesRead)
	//DeviceList = deviceList //deviceList.AddAllToDevices(&Devices)

	w.WriteHeader(200)
	_, err = w.Write(bytesRead)
	if err != nil {
		log.Println("Error in getting devices")
		w.WriteHeader(405)
		_, _ = w.Write([]byte("No Device found"))
	}

}

func SupplierTx(w http.ResponseWriter, r *http.Request) {
	log.Println([]byte("Supplier tx recv'ed"))

	bytesRead, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(errors.New("error in reading req body of SupplierTx"))
	}
	defer r.Body.Close()

	tx := data.TransactionFromJSON(bytesRead)

	pwd, _ := os.Getwd()
	fmt.Println("Current working dir is : " + pwd)

	changeValue, err := strconv.Atoi(tx.PowerUnits)
	if err != nil {
		log.Println(errors.New("Cannot read Change value in param: " + tx.PowerUnits))
	}

	parser.UpdateValueInFile("../../cmd/device-simple/randomsuppliernumberValue.txt", -changeValue)

	sendTransactionToConsumer(tx)
}

func ConsumerTx(w http.ResponseWriter, r *http.Request) {
	log.Println([]byte("Consumer tx recv'ed"))

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

	parser.UpdateValueInFile("../../cmd/device-simple/randomconsumenumberValue.txt", changeValue)
}

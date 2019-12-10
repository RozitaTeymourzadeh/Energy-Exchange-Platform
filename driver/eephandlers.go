package driver

import (
	"encoding/json"
	"strings"

	//"bytes"
	"fmt"
	"runtime"
	"strconv"

	//"github.com/edgexfoundry/device-simple/driver"
	"html/template"
	"log"
	"os"
	//"os"
	//"time"

	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

// data structure to hold readings
var DeviceEventsDS = NewDeviceEvents()

//var Devices = ds.NewDevices()
var APPNAME = "PowerFlow : Energy Exchange Platform"
var DEVICELIST = NewDeviceList()
var SELFDEVICES = NewDeviceMap()

var SUPPLYDEVICEDETAILS = make([]DeviceTypeDetails, 0)
var CONSUMEDEVICEDETAILS = make([]DeviceTypeDetails, 0)
var TRANSACTIONS = make([]Transaction, 0)
var ALLTRANSACTIONS = make([]Transaction, 0)

/// new datastructures
var LASTREADFORHEIGHT = 1
var OPENCONSUMETXS = NewTransactionPool() //make(map[string]Transaction)

//var PageVars = resources.NewPageVars()

func init() {
	//	// This function will be executed before everything else.
}

// eep handler
func Eep(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("PowerFlow : Energy Exchange Platform"))
}

// Index handler
func Index(w http.ResponseWriter, r *http.Request) {

	p := PageVars{
		Title:                 APPNAME,
		IpPort:                GetNodeId().Address + ":" + GetNodeId().Port,
		DeviceMap:             DEVICELIST.Devices, //SELFDEVICES.DeviceMapToList(),
		SupplyDevicesDetails:  SUPPLYDEVICEDETAILS,
		ConsumeDevicesDetails: CONSUMEDEVICEDETAILS,
		Transactions:          TRANSACTIONS,
		SdReadings:            GetLast100SDReadings(),
		CdReadings:            GetLast100CDReadings(),
	}

	render(w, "home.html", p)
}

// render func to serve html in templates dir
func render(w http.ResponseWriter, tmpl string, pageVars PageVars) {

	tmpl = fmt.Sprintf("../../driver/resources/templates/%s", tmpl) // prefix the name passed in with templates/
	t, err := template.ParseFiles(tmpl)                             // parse the template file held in the templates folder

	if err != nil { // if there is an error
		pwd, _ := os.Getwd()
		log.Println("Current working dir : " + pwd)
		log.Print("template parsing error: ", err) // log it
	}

	err = t.Execute(w, pageVars) //execute the template and pass in the variables to fill the gaps

	if err != nil { // if there is an error
		log.Print("template executing error: ", err) //log it
	}
}

// Show Latest Device Data
func ShowLatestDeviceData(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	deviceName := vars["deviceName"]
	resourceName := vars["resourceName"]

	latestCde, err := DeviceEventsDS.GetLatestDeviceResourceNameEventForDevice(deviceName /*"Supply-Device-01"*/, resourceName /*"randomsuppliernumber"*/)
	if err != nil {
		fmt.Println("Error in getting latest CoreEventData for a device")
		w.WriteHeader(203)
		w.Write([]byte("No data found"))
	} else {
		_, _ = w.Write(([]byte)(latestCde.Readings[0].Device + " -" + latestCde.Readings[0].Name + " = " + latestCde.Readings[0].Value))
	}

}

// Show All Latest Device Data
func ShowAllLatestDeviceData(w http.ResponseWriter, r *http.Request) {

	_, _ = w.Write([]byte(DeviceEventsDS.Show()))

}

// Show Peers
func ShowPeers(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	str := "["
	for _, peer := range GetNodeId().GetPeers() {
		str += string(peer.PeerInfoToJSON()) + ","
	}
	str = str[:len(str)-1]
	str += "]"
	w.Write([]byte(str))
}

// Read Device Data
func ReadDeviceData(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	uri := "http://localhost:48080/api/v1/event/device/" + vars["deviceName"] + "/" + vars["noOfReadings"]

	resp, err := http.Get(uri)
	if err != nil {
		fmt.Println("error in reading response body in startreading")
	}
	defer resp.Body.Close()
	bytesRead, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(bytesRead))

	cdes := CoreDataEventsFromJson(bytesRead)

	fmt.Println("coreDataEvent:")
	for _, coreDataEvent := range cdes.DataEvents {
		//fmt.Println(string(coreDataEvent.CoreDataEventToJson()))
		DeviceEventsDS.AddToDeviceEvents(coreDataEvent)
	}
	_, _ = w.Write([]byte(DeviceEventsDS.ShowDeviceEvents(vars["deviceName"])))

}

// get the count of number of go routines in the system.
func countGoRoutines() int {
	return runtime.NumGoroutine()
}
func getGoroutinesCountHandler(w http.ResponseWriter, r *http.Request) {
	// Get the count of number of go routines running.
	count := countGoRoutines()
	w.Write([]byte(strconv.Itoa(count)))
}

// show Open Consumer Tx
func OpenConsumerTx(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	sb := strings.Builder{}
	sb.WriteString("OpenConsumerTxs:\n\n")
	for _, tx := range OPENCONSUMETXS.Pool {
		sb.WriteString(tx.EventId)
	}
	w.Write([]byte(sb.String()))
}

// Send Last 100 SD Readings
func SendLast100SDReadings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Headers", "Authorization")
	readings := GetLast100SDReadings()
	barr, err := json.Marshal(&readings)
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		w.Write([]byte("no content"))
	} else {
		w.WriteHeader(http.StatusOK)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		w.Write(barr)
	}
}

// Send Last 100 CD Readings
func SendLast100CDReadings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	readings := GetLast100CDReadings()
	barr, err := json.Marshal(&readings)
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		w.Write([]byte("no content"))
	} else {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		//w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write(barr)
	}
}

// Send Sd Device Info
func SendSdDeviceInfo(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Access-Control-Allow-Origin", "*")
	devicesInfo := generateSupplyDeviceTypeBoard("supply")
	deviceInfo := devicesInfo[0]
	deviceInfo.DeviceAddress = GetNodeId().Address + ":" + GetNodeId().Port
	barr, err := json.Marshal(&deviceInfo)
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		w.Write([]byte("no content"))
	} else {
		w.WriteHeader(http.StatusOK)
		//w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write(barr)
	}
}

// Send Cd DeviceInfo
func SendCdDeviceInfo(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Access-Control-Allow-Origin", "*")
	devicesInfo := generateConsumeDeviceTypeBoard("supply")
	deviceInfo := devicesInfo[0]
	deviceInfo.DeviceAddress = GetNodeId().Address + ":" + GetNodeId().Port
	barr, err := json.Marshal(&deviceInfo)
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		w.Write([]byte("no content"))
	} else {
		w.WriteHeader(http.StatusOK)
		//w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write(barr)
	}
}

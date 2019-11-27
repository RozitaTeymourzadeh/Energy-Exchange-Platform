package uri_router

import (
	//"bytes"
	"fmt"
	"github.com/edgexfoundry/device-simple/src/data"
	"github.com/edgexfoundry/device-simple/src/resources"
	"html/template"
	"log"
	"os"
	//"os"
	//"time"

	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

//var INIT_SERVER_ADDRESS = "http://localhost:6686"
//changes in init for arg of port provided
//var SELFID = ds.NewtNodeId("localhost", 6686)

// data structure to hold readings
var DeviceEventsDS = data.NewDeviceEvents()

//var Devices = ds.NewDevices()
var APPNAME = "PowerFlow : Energy Exchange Platform"
var DEVICELIST = data.NewDeviceList() //todo is blockchain variable
var SELFDEVICES = data.NewDeviceMap()

// //DEVICEBOARD - // GetSupplyDeviceBoard() // GetConsumeDeviceBoard() - init in
var SUPPLYDEVICEDETAILS = make([]data.DeviceTypeDetails, 0)
var CONSUMEDEVICEDETAILS = make([]data.DeviceTypeDetails, 0)
var TRANSACTIONS = make([]data.Transaction, 0)

//var PageVars = resources.NewPageVars()

func init() {
	//	// This function will be executed before everything else.
	//
	//	SELF_ADDR = SELF_ADDR_PREFIX + os.Args[1]
	//	fmt.Println("Node : ", SELF_ADDR)

}

// eep handler
func Eep(w http.ResponseWriter, r *http.Request) {

	//DEVICELIST = getAllDevices( /*data.GetNodeId().ConnectingAddress*/ )
	//SUPPLYDEVICEDETAILS = generateDeviceTypeBoard("supply")
	//CONSUMEDEVICEDETAILS = generateDeviceTypeBoard("consume")

	_, _ = w.Write([]byte("PowerFlow : Energy Exchange Platform"))
}

// Index handler
func Index(w http.ResponseWriter, r *http.Request) {
	//pageVars := resources.PageVars{
	//	Title: "Energy Trading Platform",
	//}

	p := resources.PageVars{
		Title:                 APPNAME,
		DeviceList:            DEVICELIST.Devices, //SELFDEVICES.DeviceMapToList(),
		SupplyDevicesDetails:  SUPPLYDEVICEDETAILS,
		ConsumeDevicesDetails: CONSUMEDEVICEDETAILS,
		Transactions:          TRANSACTIONS,
	}

	//x := p.SupplyDevicesDetails
	//fmt.Println(len(x))
	render(w, "home.html", p)
}

// render func to serve html in templates dir
func render(w http.ResponseWriter, tmpl string, pageVars resources.PageVars) {

	tmpl = fmt.Sprintf("../../src/resources/templates/%s", tmpl) // prefix the name passed in with templates/
	t, err := template.ParseFiles(tmpl)                          //parse the template file held in the templates folder

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

func ShowAllLatestDeviceData(w http.ResponseWriter, r *http.Request) {

	_, _ = w.Write([]byte(DeviceEventsDS.Show()))

}

//func Register(w http.ResponseWriter, r *http.Request) {
//	fmt.Println("In registration service")
//	if r.Method == http.MethodPost {
//		defer r.Body.Close()
//		bytesRead, _ := ioutil.ReadAll(r.Body)
//		rInfo := data.PeerInfoFromJSON(bytesRead)
//		data.GetNodeId().AddPeer(rInfo)
//		fmt.Println(data.GetNodeId().GetPeers())
//		w.WriteHeader(200)
//	} else {
//		w.WriteHeader(405)
//	}
//
//}

func ShowPeers(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	str := "["
	for _, peer := range data.GetNodeId().GetPeers() {
		str += string(peer.PeerInfoToJSON()) + ","
	}
	str = str[:len(str)-1]
	str += "]"
	w.Write([]byte(str))
}

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

	cdes := data.CoreDataEventsFromJson(bytesRead)

	fmt.Println("coreDataEvent:")
	for _, coreDataEvent := range cdes.DataEvents {
		fmt.Println(string(coreDataEvent.CoreDataEventToJson()))
		DeviceEventsDS.AddToDeviceEvents(coreDataEvent)
	}

	//// todo: remove below code to  different endpoint
	//latestCde, err := DeviceEventsDS.GetLatestDeviceResourceNameEventForDevice("Supply-Device-01", "randomsuppliernumber")
	//if err != nil {
	//	fmt.Println("Error in getting latest CoreEventData for a device")
	//}
	//
	//_, _ = w.Write(([]byte)(latestCde.Readings[0].Device + " : " + latestCde.Readings[0].Value))

	//_ , _ = w.Write([]byte(DeviceEventsDS.ShowDevice(vars["deviceName"])))
	_, _ = w.Write([]byte(DeviceEventsDS.ShowDeviceEvents(vars["deviceName"])))

}

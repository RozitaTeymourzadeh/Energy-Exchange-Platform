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
var DEVICELIST = data.NewDeviceList()
var SUPPLYDEVICEDETAILS = make([]data.DeviceTypeDetails, 0)
var CONSUMEDEVICEDETAILS = make([]data.DeviceTypeDetails, 0)
var TRANSACTIONS = make([]data.Transaction, 0)

//var PageVars = resources.NewPageVars()

//func init() {
//	// This function will be executed before everything else.
//
//	SELF_ADDR = SELF_ADDR_PREFIX + os.Args[1]
//	fmt.Println("Node : ", SELF_ADDR)
//}

// Start handler
func Start(w http.ResponseWriter, r *http.Request) {

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
		DeviceList:            DEVICELIST.Devices,
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

	tmpl = fmt.Sprintf("src/resources/templates/%s", tmpl) // prefix the name passed in with templates/
	t, err := template.ParseFiles(tmpl)                    //parse the template file held in the templates folder

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

//func MakeDecision(w http.ResponseWriter, r *http.Request) {
//
//	go MakeDecisionHelper()
//	//// making decision
//	//if len(CONSUMEDEVICEDETAILS) < 1 {
//	//	_, _ = w.Write([]byte("No consume device"))
//	//
//	//} else if len(SupplyDeviceDetails) < 1 {
//	//	_, _ = w.Write([]byte("No supply device"))
//	//
//	//} else {
//	//
//	//	str := makeDecisionHandlerHelper()
//	//	_, _ = w.Write([]byte(str))
//	//
//	//} // end of else
//
//}

/* Event()
*
* /getevent API
* /postevent API
* To enter the event info
 */
func TaskManagerFrontend(w http.ResponseWriter, r *http.Request) {

	log.Println(".....Task Manager Frontend Method .....")

	switch r.Method {
	case "GET":
		//dir, err := os.Getwd()
		resp, err := http.Get("http://taskmanager.com/")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("response:", resp)
		http.ServeFile(w, r, "ControllerFrontend.html")
	case "POST":
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		fmt.Fprintf(w, "HTTP Post = %v\n", r.PostForm)
		fmt.Fprintf(w, "Hello astaxie!") // write data to response
		//eventName := r.FormValue("eventName")
		//eventDate := r.FormValue("eventDate")
		//eventDescription := r.FormValue("eventDescription")
		//fmt.Fprintf(w, "Event ID: %s\n", eventId)
		//fmt.Fprintf(w, "Event Name: %s\n", eventName)
		//fmt.Fprintf(w, "Event Date: %d\n", eventDate)
		//fmt.Fprintf(w, "Event Description: %s\n", eventDescription)

		//eventId := p4.StringRandom(16)
		//newTimestamp := time.Now().Unix()
		//buf := bytes.Buffer{}
		//buf.WriteString(eventId)
		//buf.WriteString(eventName)
		//buf.WriteString(eventDescription)

		//result := buf.String()
		//transactionFee:= data.TransactionFeeCalculation(result)
		///*Block Validation */
		//if userBalance - transactionFee >= 0 {
		//	userBalance = userBalance - transactionFee
		//	//minershortKey:= rsa.PublicKey{}
		//	newTransactionObject := data.NewTransaction(eventId, &minerKey.PublicKey, eventName, newTimestamp, eventDescription, transactionFee, userBalance)
		//	fmt.Println("Transaction:", newTransactionObject)
		//	transactionJSON, _ := newTransactionObject.EncodeToJson()
		//	fmt.Println("Transaction JSON:", transactionJSON)
		//	if transactionReady {
		//		encryptedPKCS1v15 := data.EncryptPKCS(&minerKey.PublicKey, transactionJSON)
		//		fmt.Println("encryptedPKCS1v15 is:", encryptedPKCS1v15)
		//		encryptedPKCS1v15Str := string(encryptedPKCS1v15)
		//		h, hashed, signature := data.SignPKCS(encryptedPKCS1v15Str, minerKey) //Private Key
		//		fmt.Println("User Signature is:", signature)
		//		fmt.Println("h is:", h)
		//		fmt.Println("hashed is:", hashed)
		//	}
		//	go TxPool.AddToTransactionPool(newTransactionObject)
		//
		//} else {
		//	fmt.Fprintf(w, "User's has not got enough balance to add Transaction! Sorry!Balance = %d\n", userBalance)
		//}
	default:
		fmt.Fprintf(w, "FATAL: Wrong HTTP Request!")
	}
}

func Register(w http.ResponseWriter, r *http.Request) { //todo
	fmt.Println("In registration service")
	if r.Method == http.MethodPost {
		defer r.Body.Close()
		bytesRead, _ := ioutil.ReadAll(r.Body)
		rInfo := data.PeerInfoFromJSON(bytesRead)
		data.GetNodeId().AddPeer(rInfo)
		fmt.Println(data.GetNodeId().GetPeers())
		w.WriteHeader(200)
	} else {
		w.WriteHeader(405)
	}

}

func DeviceFront(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok in device Front"))
}

//func DevicesInSys() {
//
//
//	uri := "http://localhost:48082/api/v1/device"
//
//	resp, err := http.Get(uri)
//	if err != nil {
//		fmt.Println("Error in getting all devices")
//	}
//	defer resp.Body.Close()
//	bytesRead, _ := ioutil.ReadAll(resp.Body)
//
//	deviceProfiles := ds.DeviceProfilesFromJson(bytesRead)
//
//
//	w.Write([]byte(deviceProfiles.ShowDeviceProfiles()))
//}

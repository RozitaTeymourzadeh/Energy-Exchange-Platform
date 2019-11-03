package uri_router

import (
	//"bytes"
	"fmt"
	"github.com/edgexfoundry/device-simple/src/data"
	"github.com/edgexfoundry/device-simple/src/resources"
	"html/template"
	"log"
	"math/rand"
	"os"
	"time"

	//"os"
	//"time"

	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

// data structure to hold readings
var DeviceEventsDS = data.NewDeviceEvents()

//var Devices = ds.NewDevices()
var DeviceList = data.NewDeviceList()
var SupplyDeviceDetails = make([]data.DeviceTypeDetails, 0)
var ConsumeDeviceDetails = make([]data.DeviceTypeDetails, 0)
var Transactions = make([]data.Transaction, 0)

//var PageVars = resources.NewPageVars()

//func init() {
//	// This function will be executed before everything else.
//
//	SELF_ADDR = SELF_ADDR_PREFIX + os.Args[1]
//	fmt.Println("Node : ", SELF_ADDR)
//}

// Start handler
func Start(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("PowerFlow : Energy Exchange Platform"))
}

// Index handler
func Index(w http.ResponseWriter, r *http.Request) {
	//pageVars := resources.PageVars{
	//	Title: "Energy Trading Platform",
	//}
	title := "Energy Trading Platform"
	DeviceList = getAllDevices( /*data.GetNodeId().ConnectingAddress*/ )
	SupplyDeviceDetails = generateDeviceTypeBoard("supply")
	ConsumeDeviceDetails = generateDeviceTypeBoard("consume")

	p := resources.PageVars{
		Title:                 title,
		DeviceList:            DeviceList.Devices,
		SupplyDevicesDetails:  SupplyDeviceDetails,
		ConsumeDevicesDetails: ConsumeDeviceDetails,
		Transactions:          Transactions,
		Body:                  "",
	}
	//PageVars.DeviceList = DeviceList.Devices //append(PageVars.DeviceList, "A", "B")

	x := p.SupplyDevicesDetails
	fmt.Println(len(x))
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

func GetAllDevices(w http.ResponseWriter, r *http.Request) {
	uri := "http://localhost:48082/api/v1/device"

	resp, err := http.Get(uri)
	if err != nil {
		fmt.Println("Error in getting all devices")
	}
	defer resp.Body.Close()
	bytesRead, _ := ioutil.ReadAll(resp.Body)

	deviceList := data.DeviceListFromJson(bytesRead)
	DeviceList = deviceList
	//deviceList.AddAllToDevices(&Devices)

	w.Write([]byte(deviceList.ShowDeviceInList()))

}

func DeleteDevice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uri := "http://localhost:48081/api/v1/device/id/" + vars["deviceId"]

	//creating a new client
	client := http.Client{}
	// creating request
	req, _ := http.NewRequest(http.MethodDelete, uri, nil)
	// fetching response
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error in reading response body in DeleteDevice")
	}
	defer resp.Body.Close()

	bytesRead, _ := ioutil.ReadAll(resp.Body)

	w.Write(bytesRead)
}

func GetAllDeviceProfiles(w http.ResponseWriter, r *http.Request) {
	uri := "http://localhost:48081/api/v1/deviceprofile"

	resp, err := http.Get(uri)
	if err != nil {
		fmt.Println("Error in getting all devices")
	}
	defer resp.Body.Close()
	bytesRead, _ := ioutil.ReadAll(resp.Body)

	deviceProfiles := data.DeviceProfilesFromJson(bytesRead)

	w.Write([]byte(deviceProfiles.ShowDeviceProfiles()))

}

func SwitchButton(w http.ResponseWriter, r *http.Request) {
	uri := "http://localhost:48081/api/v1/device/" + + vars["deviceId"] + "/SwitchButton"

	resp, err := http.Get(uri)
	if err != nil {
		fmt.Println("Error on switching")
	}
	defer resp.Body.Close()
	bytesRead, _ := ioutil.ReadAll(resp.Body)

	deviceProfiles := data.DeviceProfilesFromJson(bytesRead)//ToDo

	w.Write([]byte(deviceProfiles.ShowDeviceProfiles()))

}


func DeleteDeviceProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uri := "http://localhost:48081/api/v1/deviceprofile/id/" + vars["deviceId"]

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

func MakeDecision(w http.ResponseWriter, r *http.Request) {

	//forDecisionSupply := make(map[string]string)
	//forDecisionConsume := make(map[string]string)

	//forDevices := [2]string{"Supply-Device-01/randomsuppliernumber", "Consume-Device01/randomconsumenumber"}
	//for _, ford := range forDevices {
	//	resp, err := http.Get("http://localhost:6686/showLatestDeviceData/" + ford /*/Supply-Device-01/randomsuppliernumber"*/)
	//	if err != nil {
	//		fmt.Println("Could not fetch data for Supply-Device-01")
	//	}
	//	defer resp.Body.Close()
	//
	//	value, _ := ioutil.ReadAll(resp.Body)
	//	values := string(value[:])
	//	vals := strings.Split(values, "=")
	//	if len(vals) > 1 {
	//		if strings.Contains(vals[0], "Supply") {
	//			forDecisionSupply[vals[0]] = vals[1]
	//		} else if strings.Contains(vals[0], "Consume") {
	//			forDecisionConsume[vals[0]] = vals[1]
	//		}
	//	}
	//}

	// making decision
	if len(ConsumeDeviceDetails) < 1 {
		_, _ = w.Write([]byte("No consume device"))

	} else if len(SupplyDeviceDetails) < 1 {
		_, _ = w.Write([]byte("No supply device"))

	} else {
		sb := strings.Builder{}
		sb.WriteString("Pairing consume and supply devices:\n")
		for _, cv := range ConsumeDeviceDetails {
			matched := false
			sb.WriteString(">>> \n")
			for _, sv := range SupplyDeviceDetails {
				sval, _ := strconv.Atoi(sv.Reading)
				cval, _ := strconv.Atoi(cv.Reading)
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
					Transactions = append(Transactions, newTx)
					sb.WriteString(cv.DeviceName + " will receive " + strconv.Itoa(randPowerUnits) + " units from " + sv.DeviceName)

					go sendTransactionToSupplier(newTx)

				}
			}
			if matched == false {
				sb.WriteString("Could not match " + cv.DeviceName + " to any Supply device")
			}
		}

		_, _ = w.Write([]byte(sb.String()))

	} // end of else

}

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
		peerObj := data.PeerFromJson(bytesRead)
		data.GetNodeId().AddPeer(peerObj)
		w.WriteHeader(200)
	}
	w.WriteHeader(405)
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

package uri_router

import (
	"fmt"
	ds "github.com/edgexfoundry/device-simple/src/data"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

//var INIT_SERVER_ADDRESS = "http://localhost:6686"
//changes in init for arg of port provided
var SELFID = ds.NewtNodeId("localhost", 6686)

// data structure to hold readings
var DeviceEventsDS = ds.NewDeviceEvents()
var Devices = ds.NewDevices()

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

func GetAllDevices(w http.ResponseWriter, r *http.Request) {
	uri := "http://localhost:48082/api/v1/device"

	resp, err := http.Get(uri)
	if err != nil {
		fmt.Println("Error in getting all devices")
	}
	defer resp.Body.Close()
	bytesRead, _ := ioutil.ReadAll(resp.Body)

	deviceList := ds.DeviceListFromJson(bytesRead)

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
		fmt.Println("error in reading response body in startreading")
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

	deviceProfiles := ds.DeviceProfilesFromJson(bytesRead)

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

	cdes := ds.CoreDataEventsFromJson(bytesRead)

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

	forDecisionSupply := make(map[string]string)
	forDecisionConsume := make(map[string]string)

	forDevices := [2]string{"Supply-Device-01/randomsuppliernumber", "Consume-Device01/randomconsumenumber"}
	for _, ford := range forDevices {
		resp, err := http.Get("http://localhost:6686/showLatestDeviceData/" + ford /*/Supply-Device-01/randomsuppliernumber"*/)
		if err != nil {
			fmt.Println("Could not fetch data for Supply-Device-01")
		}
		defer resp.Body.Close()

		value, _ := ioutil.ReadAll(resp.Body)
		values := string(value[:])
		vals := strings.Split(values, "=")
		if len(vals) > 1 {
			if strings.Contains(vals[0], "Supply") {
				forDecisionSupply[vals[0]] = vals[1]
			} else if strings.Contains(vals[0], "Consume") {
				forDecisionConsume[vals[0]] = vals[1]
			}
		}
	}

	// making decision
	if len(forDecisionConsume) < 1 {
		_, _ = w.Write([]byte("No consume device"))

	} else if len(forDecisionSupply) < 1 {
		_, _ = w.Write([]byte("No consume device"))

	} else {
		sb := strings.Builder{}
		sb.WriteString("Pairing consume and supply devices:\n")
		for ck, cv := range forDecisionConsume {
			matched := false
			sb.WriteString(">>>\n")
			for sk, sv := range forDecisionSupply {
				sval, _ := strconv.Atoi(sv)
				cval, _ := strconv.Atoi(cv)
				if sval >= cval { // one supply device supplying all energy needed by the consume device
					matched = true
					forDecisionSupply[sk] = string(sval - cval)
					sb.WriteString(ck + "will receive " + cv + " units from " + sk)
				}
			}
			if matched == false {
				sb.WriteString("Could not match" + ck + "to any Supply device")
			}
		}

		_, _ = w.Write([]byte(sb.String()))

	} // end of else

}

package driver

import (
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Eep",
		"GET",
		"/eep",
		Eep,
	},
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	//Route{
	//	"Register",
	//	"POST",
	//	"/register",
	//	Register,
	//},
	Route{
		"ShowPeers",
		"GET",
		"/showpeers",
		ShowPeers,
	},

	Route{
		"ShowLatestDeviceData",
		"GET",
		"/showLatestDeviceData/{deviceName}/{resourceName}",
		ShowLatestDeviceData,
		//	localhost:6686/showLatestDeviceData/Supply-Device-01/randomsuppliernumber
		//	localhost:6686/showLatestDeviceData/Consume-Device01/randomconsumenumber
	},
	Route{
		"ShowAllLatestDeviceData",
		"GET",
		"/showAllLatestDeviceData",
		ShowAllLatestDeviceData,
		//	localhost:6686/showAllLatestDeviceData
	},
	Route{
		"ReadDeviceData",
		"GET",
		"/readDeviceData/{deviceName}/{noOfReadings}",
		ReadDeviceData,
		//	http://localhost:6686/readDeviceData/Supply-Device-01/10
		//  http://localhost:6686/readDeviceData/Consume-Device01/10
	},
	///////////////////// !!!!!!!!! ///////////////////
	//device manager apis
	///////////////////// !!!!!!!!! ///////////////////
	Route{
		"GetAllSelfDevices",
		"GET",
		"/getallselfdevices",
		GetAllSelfDevices,
		//	localhost:48082/api/v1/device
	},
	Route{
		"DeleteDevice",
		"DELETE",
		"/device/{deviceId}",
		DeleteDevice,
		//	localhost:48081/api/v1/device/id/ce13abf3-fd29-453b-9707-df679cbb60a5
	},
	Route{
		"AllDeviceProfiles",
		"GET",
		"/alldeviceprofiles",
		GetAllDeviceProfiles,
		//	localhost:48081/api/v1/deviceprofile
	},
	Route{
		"DeleteDeviceProfile",
		"DELETE",
		"/deviceprofile/{deviceId}",
		DeleteDeviceProfile,
		//	localhost:48081/api/v1/deviceprofile/id/33793ade-9e34-4e24-a8a1-7936ba693f7a
	},
	Route{
		"SendDeviceList",
		"GET",
		"/sendDeviceList",
		SendDeviceList,
		//	localhost:9999/sendDeviceList
	},
	Route{
		"SendDeviceEvents",
		"GET",
		"/sendDeviceEvents/{deviceName}/{noOfReadings}",
		SendDeviceEvents,
	},
	Route{
		"SupplierTx",
		"POST",
		"/suppliertx",
		SupplierTx,
		//	localhost:9999/suppliertx
	},
	Route{
		"ConsumerTx",
		"POST",
		"/consumertx",
		ConsumerTx,
		//	localhost:9999/consumertx
	},

	/////////blockchain
	Route{
		"Canonical",
		"GET",
		"/canonical",
		Canonical,
	},
	Route{
		"Show",
		"GET",
		"/show",
		Show,
	},
	Route{
		"Upload",
		"POST",
		"/upload",
		Upload,
	},
	Route{
		"UploadBlock",
		"GET",
		"/block/{height}/{hash}",
		UploadBlock,
	},
	Route{
		"HeartBeatReceive",
		"POST",
		"/heartbeat/receive",
		HeartBeatReceive,
	},
	Route{
		"Start",
		"GET",
		"/start",
		Start,
	},
	//Route{
	//	"GetEvent",
	//	"GET",
	//	"/getevent",
	//	Event,
	//},
	//ToDo Post Transaction
	Route{
		"PostEvent",
		"POST",
		"/postevent",
		Event,
	},
	Route{
		"GetEvent",
		"GET",
		"/getQueryEvent",
		QueryEvent,
	},
	Route{
		"PostEvent",
		"POST",
		"/postQueryEvent",
		QueryEvent,
	},
}

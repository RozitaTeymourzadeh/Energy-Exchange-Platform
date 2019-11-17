package p5

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
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
	Route{
		"GetEvent",
		"GET",
		"/getevent",
		Event,
	},
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
package main

import (
	"fmt"
	"github.com/edgexfoundry/device-simple/src/data"
	"github.com/edgexfoundry/device-simple/src/uri_router"
	"log"
	"net/http"
	"os"
)

func main() {
	fmt.Print("starting TaskManager at ")
	router := uri_router.NewRouter()

	ip := uri_router.GetIP()
	port := ""
	if len(os.Args) > 1 {
		port = os.Args[1]
	} else {
		port = "6686"
	}
	data.SetNodeId(ip, port)
	fmt.Println(ip + " : " + port)

	log.Fatal(http.ListenAndServe(data.GetNodeId().Address+":"+data.GetNodeId().Port, router))
}

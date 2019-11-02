package devicePkg

import (
	"fmt"
	"github.com/edgexfoundry/device-simple/src/data"
	"github.com/edgexfoundry/device-simple/src/uri_router"
	"log"
	"net/http"
	"os"
)

func RunDeviceManager() {
	fmt.Println("Task Manager is listening ....")
	router := uri_router.NewRouter()

	ip := "localhost" //uri_router.GetIP() //
	port := ""
	if len(os.Args) > 1 {
		port = os.Args[1]
	} else {
		port = "9999"
	}
	data.SetNodeId(ip, port)
	fmt.Println("http://" + ip + ":" + port + "/on")

	// serve everything in the css folder, the img folder and mp3 folder as a file
	pwd, _ := os.Getwd()
	fmt.Println("Current working dir : ", pwd)
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	//http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("img"))))
	//http.Handle("/mp3/", http.StripPrefix("/mp3/", http.FileServer(http.Dir("mp3"))))

	// listen and serve at ip and port
	log.Fatal(http.ListenAndServe(data.GetNodeId().Address+":"+data.GetNodeId().Port, router))
}

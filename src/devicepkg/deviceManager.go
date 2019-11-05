package devicePkg

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/edgexfoundry/device-simple/src/data"
	"github.com/edgexfoundry/device-simple/src/uri_router"
	"log"
	"net/http"
	"os"
)

func Run() {
	// // // // // // // //
	edgeXAddress := "localhost"
	taskManagerAddress := "localhost"
	taskManagerPort := "6686"
	// // // // // // // //

	fmt.Println("Device Manager is listening ....")
	router := uri_router.NewRouter()

	ip := data.SystemIp() //
	port := ""
	if len(os.Args) > 1 {
		port = os.Args[1]
	} else {
		port = "9999"
	}
	data.SetNodeId(ip, port, taskManagerAddress, taskManagerPort)
	data.GetNodeId().SetEdgeXAddress(edgeXAddress) // setting edgex address
	fmt.Println("EdgeX address is : " + edgeXAddress)

	go On()

	// serve everything in the css folder, the img folder and mp3 folder as a file
	pwd, _ := os.Getwd()
	fmt.Println("Current working dir : ", pwd)
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	//http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("img"))))
	//http.Handle("/mp3/", http.StripPrefix("/mp3/", http.FileServer(http.Dir("mp3"))))

	// listen and serve at ip and port
	log.Fatal(http.ListenAndServe(data.GetNodeId().Address+":"+data.GetNodeId().Port, router))
}

func On() {

	uri := "http://" + data.GetNodeId().TaskManagerAddress + ":" + data.GetNodeId().TaskManagerPort + "/register"
	fmt.Println(uri)

	//data.SetNodeId(SystemIp(), )
	pInfo := data.PeerInfo{
		IpAdd: data.GetNodeId().Address,
		Port:  data.GetNodeId().Port,
	}

	rJson := pInfo.PeerInfoToJSON()

	//creating a new client
	client := http.Client{}
	// creating request
	req, _ := http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(rJson))
	// fetching response
	_, err := client.Do(req)
	if err != nil {
		log.Println(errors.New("Error in device registration : " + err.Error()))
	}

	fmt.Println("On--ing device")
}

package driver

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	fmt.Println("Task Manager is listening ....")
	router := NewRouter()

	ip := "localhost" //uri_router.SystemIp() //
	port := ""
	if len(os.Args) > 1 {
		port = os.Args[1]
	} else {
		port = "6686"
	}
	SetNodeId(ip, port, "", "")
	fmt.Println("http://" + ip + ":" + port)

	// serve everything in the css folder, the img folder and mp3 folder as a file
	pwd, _ := os.Getwd()
	fmt.Println("Current working dir : ", pwd)
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	//http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("img"))))
	//http.Handle("/mp3/", http.StripPrefix("/mp3/", http.FileServer(http.Dir("mp3"))))

	/////////////////////
	// todo : if using task manager go uri_router.MakeDecision()
	/////////////////////

	// listen and serve at ip and port
	log.Fatal(http.ListenAndServe(GetNodeId().Address+":"+GetNodeId().Port, router))
}

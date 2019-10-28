package main

import (
	"fmt"
	"github.com/edgexfoundry/device-simple/src/uri_router"
	"log"
	"net/http"
	"os"
)

func main() {
	fmt.Println("Task Manager is listening ....")
	router := uri_router.NewRouter()
	if len(os.Args) > 1 {
		log.Fatal(http.ListenAndServe(":"+os.Args[1], router))
		fmt.Println("Task Manager is listening from port: ", os.Args[1])
	} else {
		log.Fatal(http.ListenAndServe(":6686", router))
		fmt.Println("Task Manager is listening from port: 6686")
	}

}

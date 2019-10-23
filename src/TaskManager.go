package main

import (
	"fmt"
	"github.com/edgexfoundry/device-simple/src/uri_router"
	"log"
	"net/http"
	"os"
)

func main() {
	fmt.Println("ok ok ")

	router := uri_router.NewRouter()
	if len(os.Args) > 1 {
		log.Fatal(http.ListenAndServe(":"+os.Args[1], router))
	} else {
		log.Fatal(http.ListenAndServe(":6686", router))
	}

}

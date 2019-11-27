package uri_router

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"github.com/edgexfoundry/device-simple/src/bcdata"
	"github.com/edgexfoundry/device-simple/src/data"
	"log"
	"net/http"
	"os"
	"time"
)

func Run() {
	// generate private key
	privatekey, err := rsa.GenerateKey(rand.Reader, 1024)

	if err != nil {
		fmt.Println(err.Error)
		os.Exit(1)
	}

	privatekey.Precompute()
	err = privatekey.Validate()
	if err != nil {
		fmt.Println(err.Error)
		os.Exit(1)
	}

	var publickey *rsa.PublicKey
	publickey = &privatekey.PublicKey

	msg := "The secret message!"

	// EncryptPKCS1v15
	encryptedPKCS1v15 := bcdata.EncryptPKCS(publickey, msg)

	// DecryptPKCS1v15
	decryptedPKCS1v15 := bcdata.DecryptPKCS(privatekey, encryptedPKCS1v15)
	fmt.Printf("PKCS1v15 decrypted [%x] to \n[%s]\n", encryptedPKCS1v15, decryptedPKCS1v15)
	// SignPKCS1v15
	message := "message to be signed"
	h, hashed, signature := bcdata.SignPKCS(message, privatekey)
	fmt.Printf("PKCS1v15 Signature : %x\n", signature)
	//VerifyPKCS1v15
	verified, err := bcdata.VerifyPKCS(publickey, h, hashed, signature)
	fmt.Println("Signature verified: ", verified)

	// // // // // // // //
	edgeXAddress := "localhost"
	// // // // // // // //
	router := NewRouter()
	ip := data.SystemIp() //
	port := ""

	if len(os.Args) > 1 {
		port = os.Args[1]

		// // // // // // // //
		data.SetNodeId(ip, port, TA_SERVER, edgeXAddress)
		//data.GetNodeId().SetEdgeXAddress(edgeXAddress) // setting edgex address
		//fmt.Println("EdgeX address is : " + edgeXAddress)
		// // // // // // // //
		data.GetSupplyDevice()  // init singleton supplydevice.go // used in simpleDriver
		data.GetConsumeDevice() // init singleton consumedevice.go // used in simpleDriver
		data.GetDeviceBoards()  // init singleton GetDeviceBoards // store for all things bc
		// // // // // // // //
		GetSelfDevices() // init self devices in router
		for _, device := range SELFDEVICES.Devices {
			fmt.Println(": : : >>> SELFDEVICES <<< : : :")
			fmt.Println("device.PeerId : " + device.PeerId)
			fmt.Println("device.Id : " + device.Id)
			fmt.Println("device.Name : " + device.Name)
		}
		// // // // // // // //
		// serve everything in the css folder, the img folder and mp3 folder as a file
		pwd, _ := os.Getwd()
		fmt.Println("Current working dir : ", pwd)
		http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
		//http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("img"))))
		//http.Handle("/mp3/", http.StripPrefix("/mp3/", http.FileServer(http.Dir("mp3"))))

	} else {
		port = "6686"
		// // // // //
		data.SetNodeId(ip, port, TA_SERVER, edgeXAddress)
		// // // // //
	}

	go startProcess()

	// listen and serve at ip and port
	fmt.Println("Device Manager is listening on : " + data.GetNodeId().Address + ":" + data.GetNodeId().Port)
	log.Fatal(http.ListenAndServe(data.GetNodeId().Address+":"+data.GetNodeId().Port, router))
}

func startProcess() {
	time.Sleep(4 * time.Second)
	http.Get("http://" + SELF_ADDR + "/start/")
}

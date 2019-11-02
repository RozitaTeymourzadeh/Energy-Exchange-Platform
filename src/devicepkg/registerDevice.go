package devicepkg

import (
	"bytes"
	"fmt"
	"github.com/edgexfoundry/device-simple/src/data"
	"io/ioutil"
	"net/http"
)

func RegisterDevice(ip string, port string, connectingIp string, connectingPort string) {

	// creating a client
	client := http.Client{}

	// creating request
	uri := "http://" + connectingIp + ":" + connectingPort + "/registry"
	ri := data.RegistryInfo{
		IpAdd: ip,
		Port:  port,
	}

	req, _ := http.NewRequest(http.MethodPost, uri, bytes.NewBufferString(string(ri.RegistryToJSON())))

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error in reading response body in RegisterDevice")
	}
	defer resp.Body.Close()

	bytesRead, _ := ioutil.ReadAll(resp.Body)

	fmt.Println("Device registered by task manager, " + string(bytesRead))

}

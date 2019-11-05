package data

import (
	"net"
	"os"
)

func SystemIp() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		_, _ = os.Stderr.WriteString("Oops: " + err.Error() + "\n")
		os.Exit(1)
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				//os.Stdout.WriteString(ipnet.IP.String() + "\n")
				return ipnet.IP.String()
			}
		}
	}

	return "localhost"
}

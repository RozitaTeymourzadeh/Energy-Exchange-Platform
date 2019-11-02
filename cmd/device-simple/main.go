// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2017-2018 Canonical Ltd
// Copyright (C) 2018-2019 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

// This package provides a simple example of a device service.
package main

import (
	"fmt"
	"github.com/edgexfoundry/device-sdk-go/pkg/startup"
	"github.com/edgexfoundry/device-simple"
	"github.com/edgexfoundry/device-simple/driver"
	"github.com/edgexfoundry/device-simple/src/devicePkg"
	"os"
)

const (
	serviceName string = "device-simple"
)

func main() {

	//// server
	//router := uri_router.NewRouter()
	//if len(os.Args) > 1 {
	//	log.Fatal(http.ListenAndServe(":"+os.Args[1], router))
	//} else {
	//	log.Fatal(http.ListenAndServe(":6686", router))
	//}
	// server end
	/////////////////////////
	/////////////////////////
	go devicePkg.RunDeviceManager() // device manager

	fmt.Fprintf(os.Stdout, "HERE.......\n")
	sd := driver.SimpleDriver{}
	startup.Bootstrap(serviceName, device.Version, &sd)
}

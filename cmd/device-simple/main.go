// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2017-2018 Canonical Ltd
// Copyright (C) 2018-2019 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

// This package provides a simple example of a device service.
package main

import (
<<<<<<< Updated upstream:edgexfoundry/device-simple/cmd/device-simple/main.go
	"github.com/edgexfoundry/device-sdk-go"
	"github.com/edgexfoundry/device-sdk-go/pkg/startup"
	//"github.com/edgexfoundry/device-sdk-go/example/driver"
	"github.com/edgexfoundry/device-simple/driver"
=======
	"fmt"
	"github.com/edgexfoundry/device-simple"
	"github.com/edgexfoundry/device-simple/driver"
	"github.com/edgexfoundry/device-sdk-go/pkg/startup"
	"os"
>>>>>>> Stashed changes:cmd/device-simple/main.go
)

const (
	serviceName string = "device-simple"
)

func main() {
	fmt.Fprintf(os.Stdout, "HERE.......\n")
	sd := driver.SimpleDriver{}
	startup.Bootstrap(serviceName, device.Version, &sd)
}

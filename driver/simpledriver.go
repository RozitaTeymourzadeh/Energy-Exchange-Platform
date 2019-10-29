// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2018 Canonical Ltd
// Copyright (C) 2018-2019 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

// This package provides a simple example implementation of
// ProtocolDriver interface.
//
package driver

import (
	"bytes"
	"errors"
	"fmt"
	dsModels "github.com/edgexfoundry/device-sdk-go/pkg/models"
	"github.com/edgexfoundry/device-simple/src/parser"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
	contract "github.com/edgexfoundry/go-mod-core-contracts/models"
	"image"
	"image/jpeg"
	"image/png"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type SimpleDriver struct {
	lc           logger.LoggingClient
	asyncCh      chan<- *dsModels.AsyncValues
	switchButton bool
}

// counters
type counters struct {
	randomsuppliernumber int32
	randomsupplierrate   int32
	randomconsumenumber  int32
}

var Counters = counters{}

func getImageBytes(imgFile string, buf *bytes.Buffer) error {
	// Read existing image from file
	img, err := os.Open(imgFile)
	if err != nil {
		return err
	}
	defer img.Close()

	// TODO: Attach MediaType property, determine if decoding
	//  early is required (to optimize edge processing)

	// Expect "png" or "jpeg" image type
	imageData, imageType, err := image.Decode(img)
	if err != nil {
		return err
	}
	// Finished with file. Reset file pointer
	img.Seek(0, 0)
	if imageType == "jpeg" {
		err = jpeg.Encode(buf, imageData, nil)
		if err != nil {
			return err
		}
	} else if imageType == "png" {
		err = png.Encode(buf, imageData)
		if err != nil {
			return err
		}
	}
	return nil
}

// Initialize performs protocol-specific initialization for the device
// service.
func (s *SimpleDriver) Initialize(lc logger.LoggingClient, asyncCh chan<- *dsModels.AsyncValues) error {
	s.lc = lc
	s.asyncCh = asyncCh
	return nil
}

// HandleReadCommands triggers a protocol Read operation for the specified device.
func (s *SimpleDriver) HandleReadCommands(deviceName string, protocols map[string]contract.ProtocolProperties, reqs []dsModels.CommandRequest) (res []*dsModels.CommandValue, err error) {
	//fmt.Fprintf(os.Stdout,  "....... %s .......\n", reqs[0].DeviceResourceName)
	if len(reqs) != 1 {
		err = fmt.Errorf("SimpleDriver.HandleReadCommands; too many command requests; only one supported")
		return
	}
	s.lc.Debug(fmt.Sprintf("SimpleDriver.HandleReadCommands: protocols: %v resource: %v attributes: %v", protocols, reqs[0].DeviceResourceName, reqs[0].Attributes))

	res = make([]*dsModels.CommandValue, 1)
	now := time.Now().UnixNano()

	if reqs[0].DeviceResourceName == "randomsuppliernumber" { // supply device
		Counters.randomsuppliernumber++
		reading := generateOnceAndReadFromFileAfter(Counters.randomsuppliernumber, 100, "randomsuppliernumberValue.txt")
		cv, _ := dsModels.NewInt32Value(reqs[0].DeviceResourceName, now, int32(reading))
		res[0] = cv
	}
	if reqs[0].DeviceResourceName == "randomsupplierrate" {
		cv, _ := dsModels.NewInt32Value(reqs[0].DeviceResourceName, now, int32(rand.Intn(10)))
		res[0] = cv
	}
	if reqs[0].DeviceResourceName == "randomconsumenumber" { // consume device
		Counters.randomconsumenumber++
		reading := generateOnceAndReadFromFileAfter(Counters.randomconsumenumber, 50, "randomconsumenumberValue.txt")
		cv, _ := dsModels.NewInt32Value(reqs[0].DeviceResourceName, now, int32(reading))
		res[0] = cv
	}
	if reqs[0].DeviceResourceName == "SwitchButton" {
		cv, _ := dsModels.NewBoolValue(reqs[0].DeviceResourceName, now, s.switchButton)
		res[0] = cv
	} else if reqs[0].DeviceResourceName == "Image" {
		// Show a binary/image representation of the switch's on/off value
		buf := new(bytes.Buffer)
		if s.switchButton == true {
			err = getImageBytes("./res/on.png", buf)
		} else {
			err = getImageBytes("./res/off.jpg", buf)
		}
		cvb, _ := dsModels.NewBinaryValue(reqs[0].DeviceResourceName, now, buf.Bytes())
		res[0] = cvb
	}
	return
}

// HandleWriteCommands passes a slice of CommandRequest struct each representing
// a ResourceOperation for a specific device resource.
// Since the commands are actuation commands, params provide parameters for the individual
// command.
func (s *SimpleDriver) HandleWriteCommands(deviceName string, protocols map[string]contract.ProtocolProperties, reqs []dsModels.CommandRequest,
	params []*dsModels.CommandValue) error {

	if len(reqs) != 1 {
		err := fmt.Errorf("SimpleDriver.HandleWriteCommands; too many command requests; only one supported")
		return err
	}
	if len(params) != 1 {
		err := fmt.Errorf("SimpleDriver.HandleWriteCommands; the number of parameter is not correct; only one supported")
		return err
	}

	s.lc.Debug(fmt.Sprintf("SimpleDriver.HandleWriteCommands: protocols: %v, resource: %v, parameters: %v", protocols, reqs[0].DeviceResourceName, params))
	var err error
	if s.switchButton, err = params[0].BoolValue(); err != nil {
		err := fmt.Errorf("SimpleDriver.HandleWriteCommands; the data type of parameter should be Boolean, parameter: %s", params[0].String())
		return err
	}

	return nil
}

// Stop the protocol-specific DS code to shutdown gracefully, or
// if the force parameter is 'true', immediately. The driver is responsible
// for closing any in-use channels, including the channel used to send async
// readings (if supported).
func (s *SimpleDriver) Stop(force bool) error {
	// Then Logging Client might not be initialized
	if s.lc != nil {
		s.lc.Debug(fmt.Sprintf("SimpleDriver.Stop called: force=%v", force))
	}
	return nil
}

// AddDevice is a callback function that is invoked
// when a new Device associated with this Device Service is added
func (s *SimpleDriver) AddDevice(deviceName string, protocols map[string]contract.ProtocolProperties, adminState contract.AdminState) error {
	s.lc.Debug(fmt.Sprintf("a new Device is added: %s", deviceName))
	return nil
}

// UpdateDevice is a callback function that is invoked
// when a Device associated with this Device Service is updated
func (s *SimpleDriver) UpdateDevice(deviceName string, protocols map[string]contract.ProtocolProperties, adminState contract.AdminState) error {
	s.lc.Debug(fmt.Sprintf("Device %s is updated", deviceName))
	return nil
}

// RemoveDevice is a callback function that is invoked
// when a Device associated with this Device Service is removed
func (s *SimpleDriver) RemoveDevice(deviceName string, protocols map[string]contract.ProtocolProperties) error {
	s.lc.Debug(fmt.Sprintf("Device %s is removed", deviceName))
	return nil
}

// GenerateOnceAndReadFromFileAfter
func generateOnceAndReadFromFileAfter(count int32, maxVal int, filename string) int {
	fmt.Println("in generateOnceAndReadFromFileAfter, count is : ", count)
	var reading int
	if count <= 1 {
		parser.DeleteFile(filename)
		reading = rand.Intn(maxVal)
		fmt.Println("Writing to file generated value : ", strconv.Itoa(reading))
		parser.WriteFile(filename, strconv.Itoa(reading))

	} else {

		fileValue := parser.ReadFile(filename)
		fmt.Println("fileValue : ", fileValue)
		reading, err := strconv.Atoi(fileValue)
		if err != nil {
			fmt.Println(errors.New("cannot convert file reading to int"))
		}
		if reading > 0 {
			reading = reading - 1
		}

		fmt.Println("Writing to added value : ", reading)
		parser.OverWriteFile(filename, strconv.Itoa(reading)) //(fmt.Sprint(reading)/*strconv.Itoa(reading)*/)
	}
	return reading
}

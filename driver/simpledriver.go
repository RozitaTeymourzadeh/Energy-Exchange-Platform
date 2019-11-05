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
	"log"
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
	supplierChargeNumber int32
	supplierRateNumber   int32
	consumerChargeNumber int32
	consumerRateNumber   int32
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

	//supply device
	if reqs[0].DeviceResourceName == "supplierCharge" { // supply device charge
		Counters.supplierChargeNumber++
		reading := generateOnceAndReadFromFileAfter(Counters.supplierChargeNumber, 100, "supplierChargeValue.txt", 1)
		log.Println("supplierCharge value: ", reading)
		cv, _ := dsModels.NewInt32Value(reqs[0].DeviceResourceName, now, int32(reading))
		res[0] = cv
	}
	if reqs[0].DeviceResourceName == "supplierRate" { // supply device rate
		Counters.supplierRateNumber++
		reading := generateOnceAndReadFromFileAfter(Counters.supplierRateNumber, 10, "supplierRateValue.txt", 0)
		log.Println("supplierRate value: ", reading)
		cv, _ := dsModels.NewInt32Value(reqs[0].DeviceResourceName, now, int32(reading))
		res[0] = cv
	}
	//consume device
	if reqs[0].DeviceResourceName == "consumerCharge" { // consume device charge
		Counters.consumerChargeNumber++
		reading := generateOnceAndReadFromFileAfter(Counters.consumerChargeNumber, 50, "consumerChargeValue.txt", -1)
		log.Println("consumerCharge value: ", int32(reading))
		cv, _ := dsModels.NewInt32Value(reqs[0].DeviceResourceName, now, int32(reading))
		res[0] = cv
	}
	if reqs[0].DeviceResourceName == "consumerRate" { // consume device rate
		Counters.consumerRateNumber++
		reading := generateOnceAndReadFromFileAfter(Counters.consumerRateNumber, 10, "consumerRateValue.txt", 0)
		log.Println("consumerRate value: ", int32(reading))
		cv, _ := dsModels.NewInt32Value(reqs[0].DeviceResourceName, now, int32(reading))
		res[0] = cv
	}

	// supply switch
	if reqs[0].DeviceResourceName == "SupplierSwitchButton" { //supplier Switch device
		cv, _ := dsModels.NewBoolValue(reqs[0].DeviceResourceName, now, s.switchButton)
		res[0] = cv
	} else if reqs[0].DeviceResourceName == "SupplierSwitchButtonImage" {
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
	// consume switch
	if reqs[0].DeviceResourceName == "ConsumerSwitchButton" { //consumer Switch device
		cv, _ := dsModels.NewBoolValue(reqs[0].DeviceResourceName, now, s.switchButton)
		res[0] = cv
	} else if reqs[0].DeviceResourceName == "ConsumerSwitchButtonImage" { //consumer Switch device Image
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
func generateOnceAndReadFromFileAfter(count int32, maxVal int, filename string, change int) int {
	fmt.Println("Event count is : ", count)
	var val int
	if count <= 1 {
		parser.DeleteFile(filename)
		val = rand.Intn(maxVal)
		//fmt.Println("Writing to file generated value : ", strconv.Itoa(val))
		parser.WriteFile(filename, strconv.Itoa(val))

	} else {

		fileValue := parser.ReadFile(filename)
		fmt.Println("fileValue : ", fileValue)
		fileVal, err := strconv.Atoi(fileValue)
		if err != nil {
			fmt.Println(errors.New("cannot convert file reading to int"))
		}
		if fileVal+change >= 0 {
			val = fileVal + change
		}

		//fmt.Println("Writing to added value : ", int32(val))
		parser.OverWriteFile(filename, strconv.Itoa(val)) //(fmt.Sprint(reading)/*strconv.Itoa(reading)*/)

		//return reading
	}
	return val
}

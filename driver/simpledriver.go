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
	"fmt"
	dsModels "github.com/edgexfoundry/device-sdk-go/pkg/models"
	"github.com/edgexfoundry/device-simple/src/data"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
	contract "github.com/edgexfoundry/go-mod-core-contracts/models"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"time"
)

type SimpleDriver struct {
	lc           logger.LoggingClient
	asyncCh      chan<- *dsModels.AsyncValues
	switchButton bool
}

//// counters
//type counters struct {
//	supplierCharge     int32
//	supplierRate       int32
//	supplierChargeRate int32
//	isSupplying        int32
//	toSupply           int32
//	sellRate           int32
//	consumerCharge     int32
//	consumerRate       int32
//	require            int32
//	isReceiving        int32
//	toReceive          int32
//	buyRate            int32
//}
//
//var Counters = counters{}

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
	//data.GetSupplyDevice()
	if reqs[0].DeviceResourceName == "supplierCharge" { // supply charge
		//Counters.supplierCharge++
		////reading := generateOnceAndReadFromFileAfter(Counters.supplierCharge, 100, "supplierChargeValue.txt", 1)
		//data.GetSupplyDevice().SupplierCharge =
		//	generateOnceAndReadFromFileAfter(Counters.supplierCharge, 100, "supplierChargeValue.txt", data.GetSupplyDevice().SupplierChargeRate)
		driverSupplierChargeUpdate()
		log.Println("supplierCharge value: ", data.GetSupplierCharge())
		cv, _ := dsModels.NewInt32Value(reqs[0].DeviceResourceName, now, int32(data.GetSupplierCharge()))
		res[0] = cv
	} //dynamic
	if reqs[0].DeviceResourceName == "supplyRate" { // supply rate
		//Counters.supplierRate++
		//data.GetSupplyDevice().SupplyRate =
		//	generateOnceAndReadFromFileAfter(Counters.supplierRate, 10, "supplierRateValue.txt", 0)
		log.Println("supplyRate value: ", data.GetSupplyRate())
		cv, _ := dsModels.NewInt32Value(reqs[0].DeviceResourceName, now, int32(data.GetSupplyRate()))
		res[0] = cv
	}
	if reqs[0].DeviceResourceName == "supplierChargeRate" { // supply charge rate
		//Counters.supplierChargeRate++
		//data.GetSupplyDevice().SupplierChargeRate =
		//	generateOnceAndReadFromFileAfter(Counters.supplierChargeRate, 10, "supplierChargeRateValue.txt", 0)
		log.Println("supplierChargeRate value: ", data.GetSupplierChargeRate())
		cv, _ := dsModels.NewInt32Value(reqs[0].DeviceResourceName, now, int32(data.GetSupplierChargeRate()))
		res[0] = cv
	}
	if reqs[0].DeviceResourceName == "isSupplying" { // is Supplying
		//Counters.isSupplying++
		//data.GetSupplyDevice().IsSupplying =
		//	generateOnceAndReadFromFileAfter(Counters.isSupplying, 0, "isSupplyingValue.txt", 0)
		if data.GetIsSupplying() == 0 && data.GetHasOffered() == false {
			driverSupplierSurplusUpdate()
		}
		log.Println("isSupplying value: ", data.GetIsSupplying())
		cv, _ := dsModels.NewInt32Value(reqs[0].DeviceResourceName, now, int32(data.GetIsSupplying()))
		res[0] = cv
	}
	if reqs[0].DeviceResourceName == "toSupply" { // to supply
		//Counters.toSupply++
		//data.GetSupplyDevice().ToSupply =
		//	generateOnceAndReadFromFileAfter(Counters.toSupply, 0, "toSupplyValue.txt", 0)
		log.Println("toSupply value: ", data.GetToSupply())
		cv, _ := dsModels.NewInt32Value(reqs[0].DeviceResourceName, now, int32(data.GetToSupply()))
		res[0] = cv
	}
	if reqs[0].DeviceResourceName == "sellRate" { // sell rate
		//Counters.sellRate++
		//data.GetSupplyDevice().SellRate =
		//	generateOnceAndReadFromFileAfter(Counters.sellRate, 20, "sellRateValue.txt", 0)
		driverSellRateUpdate()
		log.Println("sellRate value: ", data.GetSellRate())
		cv, _ := dsModels.NewInt32Value(reqs[0].DeviceResourceName, now, int32(data.GetSellRate()))
		res[0] = cv
	} // dynamic

	//consume device
	//data.GetConsumeDevice()
	if reqs[0].DeviceResourceName == "consumerCharge" { // consumer charge
		//Counters.consumerCharge++
		//data.GetConsumeDevice().ConsumerCharge =
		//	generateOnceAndReadFromFileAfter(Counters.consumerCharge, 50, "consumerChargeValue.txt", -data.GetConsumeDevice().ConsumerDischargeRate)
		driverConsumerChargeUpdate()
		log.Println("consumerCharge value: ", int32(data.GetConsumerCharge()))
		cv, _ := dsModels.NewInt32Value(reqs[0].DeviceResourceName, now, int32(data.GetConsumerCharge()))
		res[0] = cv
	} //dynamic
	if reqs[0].DeviceResourceName == "consumerDischargeRate" { // consumer rate
		//Counters.consumerRate++
		//data.GetConsumeDevice().ConsumerDischargeRate =
		//	generateOnceAndReadFromFileAfter(Counters.consumerRate, 10, "consumerRateValue.txt", 0)
		log.Println("consumerDischargeRate value: ", int32(data.GetConsumerDischargeRate()))
		cv, _ := dsModels.NewInt32Value(reqs[0].DeviceResourceName, now, int32(data.GetConsumerDischargeRate()))
		res[0] = cv
	}
	if reqs[0].DeviceResourceName == "require" { // consumer require units
		//Counters.require++
		//if data.GetConsumeDevice().ConsumerCharge < 40 && data.GetConsumeDevice().HasAsked == false {
		//	data.GetConsumeDevice().Require = 100
		//	data.GetConsumeDevice().HasAsked = true
		//} else {
		//	data.GetConsumeDevice().Require = 0
		//}
		//data.GetConsumeDevice().Require =
		//	writeAndReadFromFileAfter(Counters.require, data.GetConsumeDevice().Require, "requireValue.txt", 0)
		//generateOnceAndReadFromFileAfter(Counters.require, 10, "requireValue.txt", 0)
		driverConsumerRequireUpdate()
		log.Println("require value: ", int32(data.GetRequire()))
		cv, _ := dsModels.NewInt32Value(reqs[0].DeviceResourceName, now, int32(data.GetRequire()))
		res[0] = cv
	} //todo : bc tx
	if reqs[0].DeviceResourceName == "isReceiving" { // consumer require units
		//Counters.isReceiving++
		//data.GetConsumeDevice().IsReceiving =
		//	generateOnceAndReadFromFileAfter(Counters.isReceiving, 0, "isReceivingValue.txt", 0)
		log.Println("isReceiving value: ", int32(data.GetIsReceiving()))
		cv, _ := dsModels.NewInt32Value(reqs[0].DeviceResourceName, now, int32(data.GetIsReceiving()))
		res[0] = cv
	}
	if reqs[0].DeviceResourceName == "toReceive" { // consumer require units
		//Counters.toReceive++
		//data.GetConsumeDevice().ToReceive =
		//	generateOnceAndReadFromFileAfter(Counters.toReceive, 0, "toReceiveValue.txt", 0)
		log.Println("toReceive value: ", int32(data.GetToReceive()))
		cv, _ := dsModels.NewInt32Value(reqs[0].DeviceResourceName, now, int32(data.GetToReceive()))
		res[0] = cv
	}
	if reqs[0].DeviceResourceName == "buyRate" { // buy rate
		//Counters.buyRate++
		//data.GetConsumeDevice().BuyRate =
		//	generateOnceAndReadFromFileAfter(Counters.buyRate, 30, "buyRateValue.txt", 0)
		driverBuyRateUpdate()
		log.Println("buyRate value: ", data.GetBuyRate())
		cv, _ := dsModels.NewInt32Value(reqs[0].DeviceResourceName, now, int32(data.GetBuyRate()))
		res[0] = cv
	} // dynamic

	// supply switch
	if reqs[0].DeviceResourceName == "SSwitchButton" { //supplier Switch device
		cv, _ := dsModels.NewBoolValue(reqs[0].DeviceResourceName, now, s.switchButton)
		res[0] = cv
	} else if reqs[0].DeviceResourceName == "supplierSwitchButtonImage" {
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
	if reqs[0].DeviceResourceName == "CSwitchButton" { //consumer Switch device
		cv, _ := dsModels.NewBoolValue(reqs[0].DeviceResourceName, now, s.switchButton)
		res[0] = cv
	} else if reqs[0].DeviceResourceName == "consumerSwitchButtonImage" { //consumer Switch device Image
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

//// GenerateOnceAndReadFromFileAfter
//func generateOnceAndReadFromFileAfter(count int32, maxVal int, filename string, change int) int {
//	fmt.Println("Event count is : ", count)
//	var val int
//	if count <= 1 {
//		parser.DeleteFile(filename)
//		val = 0
//		if maxVal > 0 {
//			val = rand.Intn(maxVal)
//		}
//		//fmt.Println("Writing to file generated value : ", strconv.Itoa(val))
//		parser.WriteFile(filename, strconv.Itoa(val))
//
//	} else {
//
//		fileValue := parser.ReadFile(filename)
//		fmt.Println("fileValue : ", fileValue)
//		fileVal, err := strconv.Atoi(fileValue)
//		if err != nil {
//			fmt.Println(errors.New("cannot convert file reading to int"))
//		}
//		if fileVal+change >= 0 {
//			val = fileVal + change
//		} else if fileVal+change < 0 {
//			val = 0
//		}
//
//		//fmt.Println("Writing to added value : ", int32(val))
//		parser.OverWriteFile(filename, strconv.Itoa(val)) //(fmt.Sprint(reading)/*strconv.Itoa(reading)*/)
//
//		//return reading
//	}
//	return val
//}
//
//// GenerateOnceAndReadFromFileAfter
//func writeAndReadFromFileAfter(count int32, val int, filename string, change int) int {
//	fmt.Println("Event count is : ", count)
//	//var val int
//	if count <= 1 {
//		parser.DeleteFile(filename)
//		//val = 0
//		//if maxVal > 0 {
//		//	val = rand.Intn(maxVal)
//		//}
//		//fmt.Println("Writing to file generated value : ", strconv.Itoa(val))
//		parser.WriteFile(filename, strconv.Itoa(val))
//
//	} else {
//
//		fileValue := parser.ReadFile(filename)
//		fmt.Println("fileValue : ", fileValue)
//		fileVal, err := strconv.Atoi(fileValue)
//		if err != nil {
//			fmt.Println(errors.New("cannot convert file reading to int"))
//		}
//		if fileVal+change >= 0 {
//			val = fileVal + change
//		} else if fileVal+change < 0 {
//			val = 0
//		}
//
//		//fmt.Println("Writing to added value : ", int32(val))
//		parser.OverWriteFile(filename, strconv.Itoa(val)) //(fmt.Sprint(reading)/*strconv.Itoa(reading)*/)
//
//		//return reading
//	}
//	return val
//}

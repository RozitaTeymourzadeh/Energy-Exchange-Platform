package driver

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
)

func driverSupplierChargeUpdate() {
	supplierCharge := GetSupplierCharge()
	supplierChargeRate := GetSupplierChargeRate()
	isSupplying := GetIsSupplying()
	toSupply := GetToSupply()
	supplyRate := GetSupplyRate()

	supplierCharge = supplierCharge + supplierChargeRate
	if isSupplying > 0 {
		if toSupply > supplyRate {
			supplierCharge -= supplyRate
			SetToSupply(toSupply - supplyRate)
		} else {
			supplierCharge -= toSupply
			SupplyCompleteCleanup()
		}
	}

	if supplierCharge < 0 {
		supplierCharge = 0
	}
	if supplierCharge > 1000 {
		supplierCharge = 1000
	}

	SetSupplierCharge(supplierCharge)
	AddToLast100SDReadings(supplierCharge)
}

func driverConsumerChargeUpdate() {
	consumerCharge := GetConsumerCharge()
	consumerDischargeRate := GetConsumerDischargeRate()
	isReceiving := GetIsReceiving()
	toReceive := GetToReceive()
	toReceiveRate := GetToReceiveRate()

	consumerCharge = consumerCharge - consumerDischargeRate
	if isReceiving > 0 {
		if toReceive > 0 {
			if toReceive > toReceiveRate {
				consumerCharge += toReceiveRate
				SetToReceive(toReceive - toReceiveRate)
			} else {
				consumerCharge += toReceive
				ConsumeCompleteCleanup()
			}
		}
	}

	if consumerCharge < 0 {
		consumerCharge = 0
	}
	if consumerCharge > 1000 {
		consumerCharge = 1000
	}

	SetConsumerCharge(consumerCharge)
	AddToLast100CDReadings(consumerCharge)
}

func driverSupplierSurplusUpdate() {
	supplierCharge := GetSupplierCharge()
	//max := data.GetSupplierMaxCharge()

	threshold := GetSellThreshold()
	if supplierCharge > threshold {
		SetSurplus(supplierCharge - threshold)
	}
	fmt.Println("Surplus value : ", GetSurplus())

}

func driverConsumerRequireUpdate() {
	consumerCharge := GetConsumerCharge()
	//max := data.GetConsumerMaxCharge()

	threshold := GetBuyThreshold()
	if consumerCharge < threshold {
		//data.SetRequire(max - threshold)
		SetRequire(threshold - consumerCharge)
		fmt.Println("Require value : ", GetRequire())
		fmt.Println("HasAsked value : ", GetHasAsked())
		if GetHasAsked() == false && GetRequire() > 0 {
			newTx := NewTransaction("require", GetConsumeDeviceName(), GetConsumeDeviceId(), GetConsumeDeviceAddress(),
				strconv.Itoa(GetRequire()), strconv.Itoa(GetConsumerCharge()), strconv.Itoa(GetConsumerDischargeRate()), strconv.Itoa(GetBuyRate()),
				"", "", "", "", "", "", Balance)
			go sendCnTxToAll(newTx) //todo : send requirement tx to blockchain, check isReceiving

		}

	}
}

func sendCnTxToAll(newTx Transaction) {
	body, err := newTx.TransactionToJSON()
	if err == nil {
		uri := "http://" + GetNodeId().Address + ":" + GetNodeId().Port + "/postevent"
		fmt.Println("require tx to : " + uri)
		http.Post(uri, "application/json", bytes.NewBuffer(body))
		SetHasAsked(true)
		for peer, _ := range Peers.Copy() {
			uri := "http//:" + peer + "/postevent"
			fmt.Println("require tx to : " + uri)
			http.Post(uri, "application/json", bytes.NewBuffer(body))
		}
	}

}

func driverSellRateUpdate() {
	sellRate := GetSellRate()
	sellBaseRate := float64(GetSellBaseRate())
	supplierCharge := float64(GetSupplierCharge())
	max := float64(GetSupplierMaxCharge())
	//threshold := float64(data.GetSellThreshold())

	change := sellBaseRate * (max / supplierCharge)
	sellRate = int(change)
	if sellRate < 1 {
		sellRate = 1
	}
	SetSellRate(sellRate)

}

func driverBuyRateUpdate() {
	buyRate := GetBuyRate()
	buyBaseRate := float64(GetBuyBaseRate())
	consumerCharge := float64(GetConsumerCharge())
	max := float64(GetConsumerMaxCharge())
	//threshold := float64(data.GetBuyThreshold())

	change := buyBaseRate * (max / consumerCharge)
	buyRate = int(change)
	if buyRate < 1 {
		buyRate = 1
	}
	if buyRate > 100 {
		buyRate = 100
	}
	SetBuyRate(buyRate)
}

package driver

import (
	"bytes"
	"fmt"
	"net/http"
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
			SetToSupply(-supplyRate)
		} else {
			supplierCharge -= toSupply
			SupplyCompleteCleanup()
		}
	}

	SetSupplierCharge(supplierCharge)
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
				SetToReceive(-toReceiveRate)
			} else {
				consumerCharge += toReceive
				ConsumeCompleteCleanup()
			}
		}
	}

	SetConsumerCharge(consumerCharge)
}

func driverSupplierSurplusUpdate() {
	supplierCharge := GetSupplierCharge()
	//max := data.GetSupplierMaxCharge()

	threshold := GetSellThreshold()
	if supplierCharge > threshold {
		SetSurplus(supplierCharge - threshold)
		fmt.Println("Supplier value : ", GetSurplus())
		if GetHasOffered() == false {
			//todo : send surplus tx to blockchain, check isSupplying
		}
		SetHasOffered(true)
	}

}

func driverConsumerRequireUpdate() {
	consumerCharge := GetConsumerCharge()
	//max := data.GetConsumerMaxCharge()

	threshold := GetBuyThreshold()
	if consumerCharge < threshold {
		//data.SetRequire(max - threshold)
		SetRequire(threshold + 250 - consumerCharge)
		fmt.Println("Require value : ", GetRequire())
		if GetHasAsked() == false {
			//todo : send requirement tx to blockchain, check isReceiving
			uri := "http//:" + GetNodeId().Address + ":" + GetNodeId().Port + "/postevent"
			body := "ok" //todo : create consumer transaction
			fmt.Println("driverConsumerRequireUpdate : " + uri)
			http.Post(uri, "application/json", bytes.NewBuffer([]byte(body)))

		}
		SetHasAsked(true)

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

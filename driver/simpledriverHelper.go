package driver

import (
	"fmt"
	"github.com/edgexfoundry/device-simple/src/data"
)

func driverSupplierChargeUpdate() {
	supplierCharge := data.GetSupplierCharge()
	supplierChargeRate := data.GetSupplierChargeRate()
	isSupplying := data.GetIsSupplying()
	toSupply := data.GetToSupply()
	supplyRate := data.GetSupplyRate()

	supplierCharge = supplierCharge + supplierChargeRate
	if isSupplying > 0 {
		if toSupply > supplyRate {
			supplierCharge -= supplyRate
			data.SetToSupply(-supplyRate)
		} else {
			supplierCharge -= toSupply
			data.SupplyCompleteCleanup()
		}
	}

	data.SetSupplierCharge(supplierCharge)
}

func driverConsumerChargeUpdate() {
	consumerCharge := data.GetConsumerCharge()
	consumerDischargeRate := data.GetConsumerDischargeRate()
	isReceiving := data.GetIsReceiving()
	toReceive := data.GetToReceive()
	toReceiveRate := data.GetToReceiveRate()

	consumerCharge = consumerCharge - consumerDischargeRate
	if isReceiving > 0 {
		if toReceive > 0 {
			if toReceive > toReceiveRate {
				consumerCharge += toReceiveRate
				data.SetToReceive(-toReceiveRate)
			} else {
				consumerCharge += toReceive
				data.ConsumeCompleteCleanup()
			}
		}
	}

	data.SetConsumerCharge(consumerCharge)
}

func driverSupplierSurplusUpdate() {
	supplierCharge := data.GetSupplierCharge()
	//max := data.GetSupplierMaxCharge()

	threshold := data.GetSellThreshold()
	if supplierCharge > threshold {
		data.SetSurplus(supplierCharge - threshold)
		fmt.Println("Supplier value : ", data.GetSurplus())
		if data.GetHasOffered() == false {
			//todo : send surplus tx to blockchain, check isSupplying
		}
		data.SetHasOffered(true)
	}

}

func driverConsumerRequireUpdate() {
	consumerCharge := data.GetConsumerCharge()
	max := data.GetConsumerMaxCharge()

	threshold := data.GetBuyThreshold()
	if consumerCharge < threshold {
		data.SetRequire(max - threshold)
		fmt.Println("Require value : ", data.GetRequire())
		if data.GetHasAsked() == false {
			//todo : send requirement tx to blockchain, check isReceiving

		}
		data.SetHasAsked(true)

	}
}

func driverSellRateUpdate() {
	sellRate := data.GetSellRate()
	sellBaseRate := float64(data.GetSellBaseRate())
	supplierCharge := float64(data.GetSupplierCharge())
	max := float64(data.GetSupplierMaxCharge())
	//threshold := float64(data.GetSellThreshold())

	change := sellBaseRate * (max / supplierCharge)
	sellRate = int(change)
	if sellRate < 1 {
		sellRate = 1
	}
	data.SetSellRate(sellRate)

}

func driverBuyRateUpdate() {
	buyRate := data.GetBuyRate()
	buyBaseRate := float64(data.GetBuyBaseRate())
	consumerCharge := float64(data.GetConsumerCharge())
	max := float64(data.GetConsumerMaxCharge())
	//threshold := float64(data.GetBuyThreshold())

	change := buyBaseRate * (max / consumerCharge)
	buyRate = int(change)
	if buyRate < 1 {
		buyRate = 1
	}
	if buyRate > 100 {
		buyRate = 100
	}
	data.SetBuyRate(buyRate)
}

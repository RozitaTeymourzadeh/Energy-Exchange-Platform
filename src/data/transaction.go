package data

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Transaction struct {
	SupplierName    string `json:"supplierName"`
	SupplierId      string `json:"supplierId"`
	SupplierAddress string `json:"supplierAddress"`
	ConsumerName    string `json:"consumerName"`
	ConsumerId      string `json:"consumerId"`
	ConsumerAddress string `json:"consumerAddress"`
	PowerUnits      string `json:"powerUnits"`
}

func (tx *Transaction) TransactionToJSON() ([]byte, error) {
	jsonBytes, err := json.Marshal(tx)
	if err != nil {
		fmt.Println(errors.New("Error in Transaction to json" + err.Error()))
	}
	return jsonBytes, err
}

func TransactionFromJSON(jsonBytes []byte) Transaction {
	tx := Transaction{}
	err := json.Unmarshal(jsonBytes, &tx)
	if err != nil {
		fmt.Println(errors.New("Error in Transaction to json" + err.Error()))
	}
	return tx
}

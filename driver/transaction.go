package driver

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/crypto/sha3"
	"log"
	"sync"
	"time"
)

/* Transaction Struct
*
* Data structure for transaction data
*
 */
//type Transaction struct {
//	PublicKey   		*rsa.PublicKey `json:"publicKey"`
//	EventId     		string			`json:"eventId"`
//	EventName     		string			`json:"eventName"`
//	Timestamp  			int64			`json:"eventDate"`
//	EventDescription    string			`json:"eventDescription"`
//	TransactionFee    	int				`json:"transactionFee"`
//	Balance				int				`json:"balance"`
//}

type Transaction struct {
	EventType string `json:"eventType"` // require , supply
	EventId   string `json:"eventId"`
	Timestamp int64  `json:"eventDate"`
	// only for requirement tx
	ConsumerName          string `json:"consumerName"`
	ConsumerId            string `json:"consumerId"`
	ConsumerAddress       string `json:"consumerAddress"`
	ConsumerRequire       string `json:"consumerRequire"`
	ConsumerCharge        string `json:"consumerCharge"`
	ConsumerDischargeRate string `json:"consumerDischargeCharge"`
	BuyRate               string `json:"buyRate"`
	// for supply
	RequireEventId     string `json:"requirementEventId"`
	SupplierName       string `json:"supplierName"`
	SupplierId         string `json:"supplierId"`
	SupplierAddress    string `json:"supplierAddress"`
	SupplierToSupply   string `json:"supplierToSupply"`
	SupplierSupplyRate string `json:"supplierSupplyRate"`
	//PowerUnits      string `json:"powerUnits"`

	EventFee int `json:"eventFee"`
	Balance  int `json:"balance"`
}

/* TransactionPool Struct
*
* Data structure for transactionPool
*
 */
type TransactionPool struct {
	Pool      map[string]Transaction `json:"pool"`
	Confirmed map[string]bool        `json:"confirmed"`
	mux       sync.Mutex
}

/* NewTransaction()
*
* To return new transaction data
*
 */
func NewTransaction(eventType string, consumerName string, consumerId string, consumerAddress string, consumerRequire string,
	consumerCharge string, consumerDischargeRate string, buyRate string, requireEventId string, supplierName string, supplierId string,
	supplierAddress string, supplierToSupply string, supplierSupplyRate string, balance int) Transaction {
	tx := Transaction{
		EventType:             eventType,
		EventId:               "",
		Timestamp:             time.Now().Unix(),
		ConsumerName:          consumerName,
		ConsumerId:            consumerId,
		ConsumerAddress:       consumerAddress,
		ConsumerRequire:       consumerRequire,
		ConsumerCharge:        consumerCharge,
		ConsumerDischargeRate: consumerDischargeRate,
		BuyRate:               buyRate,
		RequireEventId:        requireEventId,
		SupplierName:          supplierName,
		SupplierId:            supplierId,
		SupplierAddress:       supplierAddress,
		SupplierToSupply:      supplierToSupply,
		SupplierSupplyRate:    supplierSupplyRate,
		EventFee:              1,
		Balance:               balance,
	}
	eventId := tx.GetEventId()
	tx.EventId = eventId

	return tx
}

func (tx *Transaction) GetEventId() string {
	jsonbarr, err := json.Marshal(&tx)
	if err != nil {
		fmt.Println("Cannot create Event id, " + err.Error())
		return ""
	}
	eventId := sha3.Sum256(jsonbarr)
	return hex.EncodeToString(eventId[:])
}

/* TransactionFeeCalculation()
*
* To calculate transaction fee for generating block
*
 */
func PowerFeeCalculation(Json string) int {

	//err := transaction.DecodeFromJson(Json)
	//if err == nil {
	//	//TODO take time stamp and calculate power fee
	//}
	PowerFee := (len(Json) * 2) / 10
	return PowerFee
}

/* EncodeToJson()
*
* To Encode Transaction from json format
*
 */
func (transaction *Transaction) EncodeToJson() (string, error) {
	jsonBytes, error := json.Marshal(transaction)
	return string(jsonBytes), error
}

/* DecodeFromJson()
*
* To Decode HeartBeatData from json format
*
 */
func (transaction *Transaction) DecodeFromJson(jsonString string) error {
	return json.Unmarshal([]byte(jsonString), transaction)
}

//////

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

////////

/* AddToTransactionPool()
*
* To add to transaction pool
*
 */
func (txp *TransactionPool) AddToTransactionPool(tx Transaction) { //duplicates in transactinon pool
	txp.mux.Lock()
	defer txp.mux.Unlock()
	if _, ok := txp.Pool[tx.EventId]; !ok {
		log.Println("In AddToTransactionPool : Adding new TX:", tx.EventId)
		txp.Pool[tx.EventId] = tx
	}
}

/* DeleteFromTransactionPool()
*
* To delete from transaction after processing
*
 */
func (txp *TransactionPool) DeleteFromTransactionPool(transactionId string) {
	txp.mux.Lock()
	defer txp.mux.Unlock()
	delete(txp.Pool, transactionId)
	log.Println("In DeleteFromTransactionPool : Deleting  TX:", transactionId)
}

/* GetTransactionPoolMap()
*
* To Get thansaction pool map
*
 */
func (txp *TransactionPool) GetTransactionPoolMap() map[string]Transaction {
	return txp.Pool
}

/* DecodeFromJson()
*
* To Decode HeartBeatData from json format
*
 */
func (txp *TransactionPool) GetOneTxFromPool(TxPool TransactionPool, userBalance int) *Transaction {
	txp.mux.Lock()
	defer txp.mux.Unlock()

	if len(TxPool.GetTransactionPoolMap()) > 0 {
		for _, tx := range TxPool.GetTransactionPoolMap() {
			if tx.EventType == "require" {
				return &tx
			}
			if tx.EventType == "supply" {
				//if tx.SupplierToSupply >= tx.ConsumerRequire {
				return &tx
				//}
			}
			//if userBalance >= transactionObject.EventFee {
			//	transactionObject.Balance = transactionObject.Balance - transactionObject.EventFee
			//	//TODO check how to add
			//	//fmt.Println("transactionObject.Balance:",transactionObject.Balance)
			//	return &transactionObject
			//}
		}
	}
	return nil
}

/* AddToConfirmedPool
*
* To add transaction into the confirmed pool
*
 */
func (txp *TransactionPool) AddToConfirmedPool(tx Transaction) { //duplicates in transactinon pool
	txp.mux.Lock()
	defer txp.mux.Unlock()

	//TODO:BUG. Transaction ID's coming "" (NULL) we should return false in that case.
	if tx.EventId == "" {
		fmt.Println("Tx ID is NULL. Do not add to CheckConfirmedPool,TX:", tx.EventId)
		return
	}
	if _, ok := txp.Confirmed[tx.EventId]; !ok {
		log.Println("In AddToConfirmedPool, TX:", tx.EventId)
		txp.Confirmed[tx.EventId] = true
	}
}

/* CheckConfirmedPool()
*
* To check for confirmed pool
*
 */
func (txp *TransactionPool) CheckConfirmedPool(tx Transaction) bool {
	txp.mux.Lock()
	defer txp.mux.Unlock()
	if tx.EventId == "" {
		fmt.Println("Tx ID is NULL. Returning false for CheckConfirmedPool,TX:", tx.EventId)
		return false
	}
	if _, ok := txp.Confirmed[tx.EventId]; ok {
		fmt.Println("Tx is in ConfirmedPool,TX:", tx.EventId)
		return true
	} else {
		fmt.Println("Tx is NOT in ConfirmedPool,TX:", tx.EventId)
		return false
	}
}

/* NewTransactionPool()
*
* Create new Transaction pool
*
 */
func NewTransactionPool() TransactionPool {
	Pool := make(map[string]Transaction)
	Confirmed := make(map[string]bool)
	mutex := sync.Mutex{}
	return TransactionPool{Pool, Confirmed, mutex}
}

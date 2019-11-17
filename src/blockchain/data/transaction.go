package data

import (

	"crypto/rsa"
	"encoding/json"
	"fmt"
	"log"
	"sync"
)

/* Transaction Struct
*
* Data structure for transaction data
*
 */
type Transaction struct {
	PublicKey   		*rsa.PublicKey `json:"publicKey"`
	EventId     		string			`json:"eventId"`
	EventName     		string			`json:"eventName"`
	Timestamp  			int64			`json:"eventDate"`
	EventDescription    string			`json:"eventDescription"`
	TransactionFee    	int				`json:"transactionFee"`
	Balance				int				`json:"balance"`
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
func NewTransaction(eventId string, publicKey *rsa.PublicKey,eventName string, timestamp int64, eventDescription string, transactionFee int, balance int) Transaction {
	return Transaction{
		EventId: eventId,
		PublicKey:  publicKey,
		EventName: eventName,
		Timestamp: timestamp,
		EventDescription: eventDescription,
		TransactionFee: transactionFee,
		Balance: balance,
	}
}

/* TransactionFeeCalculation()
*
* To calculate transaction fee for generating block
*
 */
func TransactionFeeCalculation(Json string) int{
	transactionFee := (len(Json)* 2)/10
	return transactionFee
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



/* AddToTransactionPool()
*
* To add to transaction pool
*
 */
func (txp *TransactionPool) AddToTransactionPool(tx Transaction) { //duplicates in transactinon pool
	txp.mux.Lock()
	defer txp.mux.Unlock()
	if _, ok := txp.Pool[tx.EventId]; !ok {
		log.Println("In AddToTransactionPool : Adding new TX:",tx.EventId)
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
	log.Println("In DeleteFromTransactionPool : Deleting  TX:",transactionId)
}

/* GetTransactionPoolMap()
*
* To Get thansaction pool map
*
 */
func (txp *TransactionPool) GetTransactionPoolMap() map[string]Transaction{
	return txp.Pool
}

/* DecodeFromJson()
*
* To Decode HeartBeatData from json format
*
 */
func (txp *TransactionPool) GetOneTxFromPool(TxPool TransactionPool, userBalance int) *Transaction{

	if len(TxPool.GetTransactionPoolMap()) > 0 {
		for _, transactionObject := range TxPool.GetTransactionPoolMap() {
			if userBalance >= transactionObject.TransactionFee {
				transactionObject.Balance = transactionObject.Balance - transactionObject.TransactionFee
				//TODO check how to add
				//fmt.Println("transactionObject.Balance:",transactionObject.Balance)
				return &transactionObject
			}
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
	if(tx.EventId == ""){
		fmt.Println("Tx ID is NULL. Do not add to CheckConfirmedPool,TX:",tx.EventId)
		return
	}
	if _, ok := txp.Confirmed[tx.EventId]; !ok {
		log.Println("In AddToConfirmedPool, TX:",tx.EventId)
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
	if(tx.EventId ==""){
		fmt.Println("Tx ID is NULL. Returning false for CheckConfirmedPool,TX:",tx.EventId)
		return false
	}
	if _, ok := txp.Confirmed[tx.EventId]; ok {
		fmt.Println("Tx is in ConfirmedPool,TX:",tx.EventId)
		return true
	}else{
		fmt.Println("Tx is NOT in ConfirmedPool,TX:",tx.EventId)
		return false
	}
}

/* NewTransactionPool()
*
* Create new Transaction pool
*
 */
func NewTransactionPool() TransactionPool {
	Pool :=  make(map[string]Transaction)
	Confirmed:=make(map[string]bool)
	mutex:=sync.Mutex{}
	return TransactionPool{Pool, Confirmed,mutex}
}
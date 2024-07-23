package analytics

import (
	"fmt"
	"log"

	"waas/internal/database"
)

type TransactionLog struct {
	TxnID         string
	WalletAddress string
	TargetAddress string
	TokenName     string
	Amount        string
	Status        string
	ErrorMessage  string
	Timestamp     string
}

// LogEvent logs an event to the database
func LogEvent(event string, data string) error {
	err := database.LogEvent(event, data)
	if err != nil {
		log.Printf("Error logging event: %v", err)
		return fmt.Errorf("error logging event: %v", err)
	}
	return nil
}

func StoreTransaction(txn TransactionLog) error {
	err := database.StoreTxnInDb(txn.WalletAddress, txn.TxnID, txn.TargetAddress, txn.TokenName, txn.Amount, txn.Status, txn.ErrorMessage, txn.Timestamp)
	if err != nil {
		log.Printf("Error storing transaction: %v", err)
		return fmt.Errorf("error storing transaction: %v", err)
	}
	return nil
}

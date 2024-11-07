package analytics

import (
	"fmt"
	"log"

	"waas/internal/database"
)

type TransactionLog struct {
	TxnID         string
	TxnHash       string
	WalletAddress string
	TargetAddress string
	TokenName     string
	Amount        string
	AmountInUSD   string
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
	err := database.StoreTxnInDb(txn.TxnID, txn.WalletAddress, txn.TargetAddress, txn.TxnHash, txn.Amount, txn.AmountInUSD, txn.TokenName, txn.Status, txn.ErrorMessage, txn.Timestamp)
	if err != nil {
		log.Printf("Error storing transaction: %v", err)
		return fmt.Errorf("error storing transaction: %v", err)
	}
	return nil
}

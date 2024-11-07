package database

import (
	"database/sql"
	"fmt"

	"waas/config"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

var DB *sql.DB

func Init() (*sql.DB, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Error("failed to load configuration: " + err.Error())
		return nil, err
	}

	var sqlUsername = cfg.SqlUsername
	var sqlPassword = cfg.SqlPassword
	var sqlUrl = cfg.SqlUrl
	var dbName = cfg.DbName

	var sqlparams = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", sqlUsername, sqlPassword, sqlUrl, dbName)

	DB, err = sql.Open("mysql", sqlparams)
	if err != nil {
		log.Error("Error connecting to database", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Error("Error pinging database ", err)
	}
	return DB, err
}

func StoreWalletDetails(nonce []byte, address string, hashedpin []byte, ciphertext []byte) error {
	stmt, err := DB.Prepare("INSERT INTO wallets (nonce, address, hashedpin, ciphertext) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(nonce, address, hashedpin, ciphertext)
	if err != nil {
		return err
	}

	return nil
}

func GetWalletDetails(address string) ([]byte, []byte, error) {
	var nonce []byte
	var hashedpin []byte
	var ciphertext []byte

	err := DB.QueryRow("SELECT nonce, hashedpin, ciphertext FROM wallets WHERE address = ?", address).Scan(&nonce, &hashedpin, &ciphertext)
	if err != nil {
		return nil, nil, err
	}

	return nonce, ciphertext, nil
}

func ValidateAPIKey(APIKey string) (bool, error) {
	var email string
	row := DB.QueryRow("SELECT email from wallet_api_keys WHERE key = ?", APIKey)
	err := row.Scan(&email)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	} else {
		// check if email is in db
		var emailExists bool
		err = DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)", email).Scan(&emailExists)
		if err != nil {
			log.Error("Error getting user from database", err)
			return false, err
		}
	}
	return true, nil
}

func GetFromDb(whatToSelect, uniqueID, tableName string) (string, error) {
	var userID string
	err := DB.QueryRow(fmt.Sprintf("SELECT %s FROM %s WHERE walletID = ?", whatToSelect, tableName), uniqueID).Scan(&userID)
	if err != nil {
		return "", err
	}
	return userID, nil
}

// func StoreInDb(whatToInsert, uniqueID, tableName string) error {
// 	stmt, err := DB.Prepare(fmt.Sprintf("INSERT INTO %s (%s) VALUES (?)", tableName, whatToInsert))
// 	if err != nil {
// 		return err
// 	}
// 	defer stmt.Close()

// 	_, err = stmt.Exec(uniqueID)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

func StoreTxnInDb(transactionID, walletID, targetAddress, transactionHash, amount, amountinUSD, tokenName, status, errorMessage, timestamp string) error {
	stmt, err := DB.Prepare("INSERT INTO wallet_transactions (transactionID, walletID, targetAddress, txHash, amount, amount_usd, tokenName, status, errorMessage, timestamp) VALUES (?, ?, ?, ?, ?, ?, ?, ?,?,?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(transactionID, walletID, targetAddress, transactionHash, amount, amountinUSD, tokenName, status, errorMessage, timestamp)
	if err != nil {
		return err
	}

	return nil
}

func GetKeyHash(id int) (string, error) {
	var keyHash string
	err := DB.QueryRow("SELECT hash FROM wallet_keys WHERE user = ?", id).Scan(&keyHash)
	if err != nil {
		return "", err
	}
	return keyHash, nil
}

func LogEvent(event, data string) error {
	_, err := DB.Exec("INSERT INTO events (event, data) VALUES (?, ?)", event, data)
	if err != nil {
		log.Printf("Error logging event: %v", err)
		return err
	}
	return nil
}

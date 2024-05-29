package tools

import (
	"database/sql"
	// "errors"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

var DB *sql.DB
var sqlUsername = os.Getenv("SQL_USERNAME")
var sqlPassword = os.Getenv("SQL_PASSWORD")
var sqlUrl = os.Getenv("SQL_URL")
var dbName = os.Getenv("DB_NAME")
var sqlparams = fmt.Sprintf("root:12345678@tcp(127.0.0.1)/waas?charset=utf8mb4&parseTime=True&loc=Local")

func init() {
	var err error
	DB, err = sql.Open("mysql", sqlparams)
	if err != nil {
		log.Error("Error connecting to database", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Error("Error pinging database ", err)
	}
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

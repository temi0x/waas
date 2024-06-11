package database

import (
	"database/sql"
	// "errors"
	"fmt"
	// "os"

	"waas/config"

	_ "github.com/go-sql-driver/mysql"
	// "github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
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
	// godotenv.Load()
	var sqlparams = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", sqlUsername, sqlPassword, sqlUrl, dbName)
	// var sqlparams = fmt.Sprintf("root:12345678@tcp(127.0.0.1)/waas?charset=utf8mb4&parseTime=True&loc=Local")
	// var sqlUsername = os.Getenv("SQL_USERNAME")
	// var sqlPassword = os.Getenv("SQL_PASSWORD")
	// var sqlUrl = os.Getenv("SQL_URL")
	// var dbName = os.Getenv("DB_NAME")

	// godotenv.Load()
	// var sqlparams = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", sqlUsername, sqlPassword, sqlUrl, dbName)

	// var sqlparams = fmt.Sprintf("root:12345678@tcp(127.0.0.1)/waas?charset=utf8mb4&parseTime=True&loc=Local")

	// var err error
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

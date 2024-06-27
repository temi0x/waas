package wallet

import (
	"math/big"
	"strconv"

	etherParams "github.com/ethereum/go-ethereum/params"
	log "github.com/sirupsen/logrus"
)

func ConvertToWei(amount float64) *big.Int {
	// Convert the amount to a string
	amountStr := strconv.FormatFloat(amount, 'f', -1, 64)

	// Convert the string to a big.Int
	amountBigFloat, ok := new(big.Float).SetString(amountStr)
	if !ok {
		log.Fatal("Invalid amount")
	}

	// Multiply the amount by 1 Ether (in Wei) to convert it to Wei
	wei := new(big.Float).Mul(amountBigFloat, new(big.Float).SetFloat64(etherParams.Ether))

	// Convert the big.Float to a big.Int
	weiInt := new(big.Int)
	wei.Int(weiInt)

	return weiInt
}

package wallet

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"waas/api"
	"waas/config"
	"waas/internal/database"
	"waas/internal/tools/transaction"
)

func SendTokenFVM(request *api.SendCustomTokenParams) (txhash string, err error) {
	var privateKey []byte
	var RPC_URL string

	cfg, err := config.LoadConfig()
	if err != nil {
		return "", fmt.Errorf("failed to load configuration: %v", err)
	}
	platformPIN := cfg.CrypteaKey

	ChainID := transaction.GetChainID(request.Chain)

	RPC_URL = GetRPC(ChainID)

	if RPC_URL == "" {
		return "", fmt.Errorf("unsupported chain: %v", request.Chain)
	}

	if !common.IsHexAddress(request.TargetAddress) {
		return "", fmt.Errorf("invalid target address")
	}

	eNonce, ciphertext, err := database.GetWalletDetails(request.UserAddress)
	if err != nil {
		return "", fmt.Errorf("error getting wallet details: %v", err)
	}

	derivedKey, err := Decrypt(eNonce, ciphertext, request.PIN, platformPIN)
	if err != nil {
		return "", fmt.Errorf("error decrypting user wallet: %v", err)
	}
	derivedKeyWithPrefix := fmt.Sprintf("0x%s", derivedKey)

	privateKey, err = hexutil.Decode(derivedKeyWithPrefix)
	if err != nil {
		return "", fmt.Errorf("error decoding private key: %v", err)
	}

	privKey, err := crypto.ToECDSA(privateKey)
	if err != nil {
		return "", fmt.Errorf("error converting to ECDSA private key: %v", err)
	}

	client, err := ethclient.Dial(RPC_URL)
	if err != nil {
		return "", fmt.Errorf("error connecting to Ethereum client: %v", err)
	}

	chainID := big.NewInt(int64(ChainID))
	auth, err := bind.NewKeyedTransactorWithChainID(privKey, chainID)
	if err != nil {
		return "", fmt.Errorf("error creating transactor: %v", err)
	}

	nonce, err := client.PendingNonceAt(context.Background(), auth.From)
	if err != nil {
		return "", fmt.Errorf("error getting nonce: %v", err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return "", fmt.Errorf("error getting gas price: %v", err)
	}

	amount := transaction.ConvertToWei(request.Amount)
	value := new(big.Int).Set(amount)
	gasLimit := uint64(780863) // in units

	toAddress := common.HexToAddress(request.TargetAddress)
	if toAddress == (common.Address{}) {
		return "", fmt.Errorf("invalid target address")
	}

	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privKey)
	if err != nil {
		return "", fmt.Errorf("error signing transaction: %v", err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return "", fmt.Errorf("error sending transaction: %v", err)
	}

	txHash := transaction.GetBlockExplorerURL(ChainID, signedTx.Hash().Hex())

	return txHash, nil
}

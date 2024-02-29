package wallet

import (
	"crypto/ecdsa"
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tuneinsight/lattigo/v5/schemes/ckks"
)

var Ckey = os.Getenv("CrypteaKey")

func CreateWallet() (string, string) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)

	fmt.Println("Private Key: ", hexutil.Encode(privateKeyBytes)[2:]) // removes 0x
	pKey := hexutil.Encode(privateKeyBytes)[2:]

	publicKey := privateKey.Public()

	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	fmt.Println(hexutil.Encode(publicKeyBytes)[4:]) // removes 0x04 which is added by default

	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	fmt.Println(address) // public address

	return pKey, address
}

func EncryptPkey(privateKey string, pin string) string {
// Define parameters
paramsDef := ckks.PN12QP109 // Parameter set for CKKS scheme
params, err := ckks.NewParametersFromLiteral(paramsDef)
if err != nil {
		panic(err)
}

// Load the keys from secure storage instead of generating them
sk, pk := pin + Ckey, "pin"+ Ckey // Example keys

// Create a new CKKS context
ckksContext := ckks.NewContext(params)

// Encoder and Encryptor using loaded keys
encoder := ckksContext.NewEncoder()
encryptor := ckksContext.NewEncryptorFromPk(pk)
decryptor := ckksContext.NewDecryptor(sk)

// Proceed with encryption and decryption as before
privateKeySimulated := []complex128{complex(123.456, 0)} // Example private key
plaintext := ckksContext.NewPlaintext(params.MaxLevel(), params.DefaultScale())
encoder.Encode(plaintext, privateKeySimulated, params.LogSlots())
ciphertext := encryptor.EncryptNew(plaintext)

decryptedPlaintext := decryptor.DecryptNew(ciphertext)
decryptedValues := encoder.Decode(decryptedPlaintext, params.LogSlots())

// Display results
fmt.Printf("Original: %f\n", privateKeySimulated[0])
fmt.Printf("Decrypted: %f\n", decryptedValues[0])
error := cmplx.Abs(privateKeySimulated[0] - decryptedValues[0])
fmt.Printf("Decryption error: %f\n", error)
}

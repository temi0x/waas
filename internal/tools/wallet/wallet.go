package wallet

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
	"log"

	// "math/cmplx"
	"os"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	// "github.com/tuneinsight/lattigo/v5/schemes/ckks"
)

var platformPIN = os.Getenv("CrypteaKey")

func deriveKey(userPIN, platformPIN string) []byte {
	pinCombo := userPIN + platformPIN
	hash := sha256.Sum256([]byte(pinCombo))
	return hash[:]
}

func generateNonce(size int) ([]byte, error) {
	nonce := make([]byte, size)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return nonce, nil
}

// encrypt encrypts the plaintext using AES-GCM with the derived key
func Encrypt(pKey, userPIN, platformPIN string) ([]byte, []byte, []byte, error) {
	key := deriveKey(userPIN, platformPIN)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, nil, err
	}

	nonce, err := generateNonce(aesGCM.NonceSize())
	if err != nil {
		return nil, nil, nil, err
	}

	ciphertext := aesGCM.Seal(nil, nonce, []byte(pKey), nil)
	return nonce, ciphertext, key, nil
}

// decrypt decrypts the ciphertext using AES-GCM with the derived key
func Decrypt(nonce, ciphertext []byte, userPIN, platformPIN string) (string, error) {
	key := deriveKey(userPIN, platformPIN)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// func main() {
// 	userPIN := "user1234"
// 	platformPIN := "platform5678"
// 	message := "Secret Wallet Data"

// 	nonce, ciphertext, err := Encrypt(message, userPIN, platformPIN)
// 	if err != nil {
// 		fmt.Println("Error encrypting:", err)
// 		return
// 	}

// 	fmt.Printf("Encrypted: %x\n", ciphertext)

// 	plaintext, err := Decrypt(nonce, ciphertext, userPIN, platformPIN)
// 	if err != nil {
// 		fmt.Println("Error decrypting:", err)
// 		return
// 	}

// 	fmt.Printf("Decrypted: %s\n", plaintext)
// }

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

package main

import (
	"crypto"
	"crypto/ecdh"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
)

type Transaction struct {
	FromAddress string
	ToAddress   string
	Amount      float64
	Signature   []byte
}

func NewTransaction(from, to string, amount float64) Transaction {
	return Transaction{
		FromAddress: from,
		ToAddress:   to,
		Amount:      amount,
	}
}

func (t *Transaction) calculateHash() string {
	data := t.FromAddress + t.ToAddress + fmt.Sprintf("%.6f", t.Amount)
	return calculateSHA256(data)
}

func (t *Transaction) signTransaction(signingKey *ecdsa.PrivateKey) {
	// Convert ECDSA public key to ECDH to get the encoded bytes (non-deprecated)
	ecdhKey, err := signingKey.PublicKey.ECDH()
	if err != nil {
		fmt.Println("Error converting to ECDH key:", err)
		return
	}
	pubKeyHex := hex.EncodeToString(ecdhKey.Bytes())

	if pubKeyHex != t.FromAddress {
		fmt.Println("You cannot sign transactions for other wallets!")
		return
	}

	hashTx := t.calculateHash()

	hashBytes, err := hex.DecodeString(hashTx)
	if err != nil {
		fmt.Println("Error decoding hash:", err)
		return
	}

	sig, err := signingKey.Sign(rand.Reader, hashBytes, crypto.SHA256)
	if err != nil {
		fmt.Println("Error signing transaction:", err)
		return
	}
	t.Signature = sig
}

func (t *Transaction) isValid() (bool, error) {
	if t.FromAddress == "" {
		return true, nil // Mining reward
	}

	if len(t.Signature) == 0 {
		return false, fmt.Errorf("no signature in this transaction")
	}

	publicKeyBytes, err := hex.DecodeString(t.FromAddress)
	if err != nil {
		return false, fmt.Errorf("error decoding public key: %v", err)
	}

	ecdhPubKey, err := ecdh.P256().NewPublicKey(publicKeyBytes)
	if err != nil {
		return false, fmt.Errorf("invalid public key: %v", err)
	}

	// Extract X and Y coordinates from the ECDH public key bytes
	// Format: [0x04 || X (32 bytes) || Y (32 bytes)] for uncompressed P256
	keyBytes := ecdhPubKey.Bytes()
	x := new(big.Int).SetBytes(keyBytes[1:33])
	y := new(big.Int).SetBytes(keyBytes[33:65])

	publicKey := ecdsa.PublicKey{Curve: elliptic.P256(), X: x, Y: y}

	hashTx := t.calculateHash()
	hashBytes, err := hex.DecodeString(hashTx)
	if err != nil {
		return false, fmt.Errorf("error decoding hash: %v", err)
	}

	valid := ecdsa.VerifyASN1(&publicKey, hashBytes, t.Signature)
	if !valid {
		return false, fmt.Errorf("invalid transaction signature")
	}
	return true, nil
}

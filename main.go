package main

import (
	"crypto/elliptic"
	"fmt"
)

func main() {
	blockchainDifficulty := 2
	miningReward := 100.0
	bc := NewBlockchain(blockchainDifficulty, miningReward)

	privateKey, publicKey, err := GenerateKeyPair()
	if err != nil {
		fmt.Println("Error generating key pair:", err)
		return
	}

	pubKeyBytes := elliptic.Marshal(elliptic.P256(), publicKey.X, publicKey.Y)
	myWalletAddress := fmt.Sprintf("%x", pubKeyBytes)
	toAddress := "someone_else's_public_key"

	tx1 := NewTransaction(myWalletAddress, toAddress, 10)
	tx1.signTransaction(privateKey)
	if err := bc.AddTransaction(tx1); err != nil {
		fmt.Println("Transaction failed:", err)
		return
	}

	fmt.Println("Starting the miner...")
	bc.MinePendingTransactions(myWalletAddress)

	fmt.Println("Mining again to collect reward...")
	bc.MinePendingTransactions(myWalletAddress)

	fmt.Printf("Balance of miner is: %.2f\n", bc.GetBalanceOfAddress(myWalletAddress))
}

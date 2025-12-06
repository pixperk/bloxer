package main

import "fmt"

func main() {
	blockchainDifficulty := 2
	miningReward := 100.0
	bc := NewBlockchain(blockchainDifficulty, miningReward)

	bc.CreateTransaction(NewTransaction("addr1", "addr2", 24.56))
	bc.CreateTransaction(NewTransaction("addr2", "addr1", 10.0))

	fmt.Println("Starting the miner...")
	bc.MinePendingTransactions("miner-address")

	fmt.Printf("Balance of miner is: %.2f\n", bc.GetBalanceOfAddress("miner-address"))
	fmt.Printf("Addr1 balance is: %.2f\n", bc.GetBalanceOfAddress("addr1"))
	fmt.Printf("Addr2 balance is: %.2f\n", bc.GetBalanceOfAddress("addr2"))

	fmt.Println("Starting the miner...")
	bc.MinePendingTransactions("miner-address")

	fmt.Printf("Balance of miner is: %.2f\n", bc.GetBalanceOfAddress("miner-address"))
}

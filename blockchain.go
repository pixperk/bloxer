package main

import (
	"fmt"
	"time"
)

type Blockchain struct {
	Chain               []Block
	Difficulty          int
	PendingTransactions []Transaction
	MiningReward        float64
}

func NewBlockchain(difficulty int, miningReward float64) *Blockchain {
	bc := &Blockchain{
		Chain:               []Block{},
		Difficulty:          difficulty,
		PendingTransactions: []Transaction{},
		MiningReward:        miningReward,
	}
	bc.Chain = append(bc.Chain, NewGenesisBlock())
	return bc
}

func (bc *Blockchain) GetLatestBlock() Block {
	if len(bc.Chain) == 0 {
		return Block{}
	}
	return bc.Chain[len(bc.Chain)-1]
}

func (bc *Blockchain) MinePendingTransactions(miningRewardAddress string) {
	currentTimeStamp := time.Now().Unix()
	pendingTx := bc.PendingTransactions
	block := NewBlock(currentTimeStamp, map[string]interface{}{"transactions": pendingTx})
	block.PrevHash = bc.GetLatestBlock().Hash
	block.Hash = block.calculateHash()

	block.MineBlock(bc.Difficulty)

	fmt.Println("Block successfully mined!")

	bc.Chain = append(bc.Chain, block)

	bc.PendingTransactions = []Transaction{
		NewTransaction("", miningRewardAddress, bc.MiningReward),
	}
}

func (bc *Blockchain) IsChainValid() bool {
	for i := 1; i < len(bc.Chain); i++ {
		currentBlock := bc.Chain[i]
		prevBlock := bc.Chain[i-1]

		if currentBlock.Hash != currentBlock.calculateHash() {
			return false
		}

		if currentBlock.PrevHash != prevBlock.Hash {
			return false
		}
	}
	return true
}

func (bc *Blockchain) CreateTransaction(transaction Transaction) {
	bc.PendingTransactions = append(bc.PendingTransactions, transaction)
}

func (bc *Blockchain) GetBalanceOfAddress(address string) float64 {
	balance := 0.0

	transactions := []Transaction{}

	for _, block := range bc.Chain {
		transactionsData, ok := block.Data["transactions"]
		if !ok {
			continue
		}
		txs, ok := transactionsData.([]Transaction)
		if ok {
			transactions = append(transactions, txs...)
		}
	}

	for _, tx := range transactions {
		if tx.FromAddress == address {
			balance -= tx.Amount
		}
		if tx.ToAddress == address {
			balance += tx.Amount
		}
	}

	return balance
}

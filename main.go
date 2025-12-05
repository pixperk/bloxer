package main

import (
	"crypto/sha256"
	"fmt"
	"strings"
	"time"
)

type Block struct {
	Data      map[string]interface{}
	PrevHash  string
	TimeStamp int64
	Hash      string
	Nonce     int
}

type Transaction struct {
	FromAddress string
	ToAddress   string
	Amount      float64
}

type Blockchain struct {
	Chain               []Block
	Difficulty          int
	PendingTransactions []Transaction
	MiningReward        float64
}

func NewBlock(timestamp int64, data map[string]interface{}) Block {
	b := Block{
		TimeStamp: timestamp,
		Data:      data,
		Nonce:     0,
	}
	b.Hash = b.calculateHash()
	return b
}

func NewBlockchain(difficulty int, miningReward float64) *Blockchain {
	return &Blockchain{
		Chain:               []Block{},
		Difficulty:          difficulty,
		PendingTransactions: []Transaction{},
		MiningReward:        miningReward,
	}
}

func NewTransaction(from, to string, amount float64) Transaction {
	return Transaction{
		FromAddress: from,
		ToAddress:   to,
		Amount:      amount,
	}
}

func (b *Block) calculateHash() string {
	record := fmt.Sprintf("%d%v%s%d", b.TimeStamp, b.Data, b.PrevHash, b.Nonce)
	hash := sha256.Sum256([]byte(record))
	return fmt.Sprintf("%x", hash)
}

func (b *Block) MineBlock(difficulty int) {
	for b.Hash[:difficulty] != strings.Repeat("0", difficulty) {
		b.Nonce++
		b.Hash = b.calculateHash()
	}

	fmt.Printf("Block mined: %s\n", b.Hash)
}

func (bc *Blockchain) CreateGenesisBlock() Block {

	currentTimeStamp := time.Now().Unix()

	genesisBlock := NewBlock(currentTimeStamp, map[string]interface{}{"message": "Genesis Block"})

	genesisBlock.PrevHash = "0"

	return genesisBlock
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

func main() {
	blockchain_difficulty := 2
	mining_reward := 100.0
	bc := NewBlockchain(blockchain_difficulty, mining_reward)
	genesisBlock := bc.CreateGenesisBlock()
	bc.Chain = append(bc.Chain, genesisBlock)

	bc.CreateTransaction(NewTransaction("addr1", "addr2", 24.56))
	bc.CreateTransaction(NewTransaction("addr2", "addr1", 10.0))

	fmt.Println("Starting the miner...")
	bc.MinePendingTransactions("miner-address")

	fmt.Printf("Balance of miner is: %.2f\n", bc.GetBalanceOfAddress("miner-address"))

	fmt.Println("Starting the miner...")
	bc.MinePendingTransactions("miner-address")

	fmt.Printf("Balance of miner is: %.2f\n", bc.GetBalanceOfAddress("miner-address"))

}

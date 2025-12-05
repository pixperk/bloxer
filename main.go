package main

import (
	"crypto/sha256"
	"fmt"
	"strings"
)

type Block struct {
	Index    int
	Data     map[string]interface{}
	PrevHash string
	Hash     string
	Nonce    int
}

type Blockchain struct {
	Chain      []Block
	Difficulty int
}

func NewBlock(index int, data map[string]interface{}, prevHash string) Block {
	b := Block{
		Index:    index,
		Data:     data,
		PrevHash: prevHash,
		Nonce:    0,
	}
	b.Hash = b.calculateHash()
	return b
}

func (b *Block) calculateHash() string {
	record := fmt.Sprintf("%d%v%s%d", b.Index, b.Data, b.PrevHash, b.Nonce)
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

	genesisBlock := NewBlock(0, map[string]interface{}{"message": "Genesis Block"}, "0")
	return genesisBlock
}

func (bc *Blockchain) GetLatestBlock() Block {
	if len(bc.Chain) == 0 {
		return Block{}
	}
	return bc.Chain[len(bc.Chain)-1]
}

func (bc *Blockchain) AddBlock(newBlock Block, data map[string]interface{}) {
	prevBlock := bc.GetLatestBlock()

	newBlock.PrevHash = prevBlock.Hash
	newBlock.MineBlock(bc.Difficulty)
	bc.Chain = append(bc.Chain, newBlock)
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

func main() {
	bc := &Blockchain{Difficulty: 6}
	genesisBlock := bc.CreateGenesisBlock()
	bc.Chain = append(bc.Chain, genesisBlock)

	newData := map[string]interface{}{
		"sender":   "Alice",
		"receiver": "Bob",
		"amount":   50,
	}
	newBlock := NewBlock(1, newData, "")
	fmt.Printf("Mining block 1...\n")
	bc.AddBlock(newBlock, newData)

	anotherNewData := map[string]interface{}{
		"sender":   "Bob",
		"receiver": "Charlie",
		"amount":   30,
	}
	anotherNewBlock := NewBlock(2, anotherNewData, "")
	fmt.Printf("Mining block 2...\n")
	bc.AddBlock(anotherNewBlock, anotherNewData)
}

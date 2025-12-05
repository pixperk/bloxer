package main

import (
	"crypto/sha256"
	"fmt"
)

type Block struct {
	Index    int
	Data     map[string]interface{}
	PrevHash string
	Hash     string
}

type Blockchain struct {
	Chain []Block
}

func NewBlock(index int, data map[string]interface{}, prevHash string) Block {
	b := Block{
		Index:    index,
		Data:     data,
		PrevHash: prevHash,
	}
	b.Hash = b.calculateHash()
	return b
}

func (b *Block) calculateHash() string {
	record := fmt.Sprintf("%d%v%s", b.Index, b.Data, b.PrevHash)
	hash := sha256.Sum256([]byte(record))
	return fmt.Sprintf("%x", hash)
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
	newBlock.Hash = newBlock.calculateHash()
	bc.Chain = append(bc.Chain, newBlock)
}

func main() {
	bc := &Blockchain{}
	genesisBlock := bc.CreateGenesisBlock()
	bc.Chain = append(bc.Chain, genesisBlock)

	newData := map[string]interface{}{
		"sender":   "Alice",
		"receiver": "Bob",
		"amount":   50,
	}
	newBlock := NewBlock(1, newData, "")
	bc.AddBlock(newBlock, newData)

	anotherNewData := map[string]interface{}{
		"sender":   "Bob",
		"receiver": "Charlie",
		"amount":   30,
	}
	anotherNewBlock := NewBlock(2, anotherNewData, "")
	bc.AddBlock(anotherNewBlock, anotherNewData)

	for _, block := range bc.Chain {
		fmt.Printf("Index: %d\n", block.Index)
		fmt.Printf("Data: %v\n", block.Data)
		fmt.Printf("PrevHash: %s\n", block.PrevHash)
		fmt.Printf("Hash: %s\n", block.Hash)
		fmt.Println()
	}
}

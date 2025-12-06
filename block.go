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

func NewBlock(timestamp int64, data map[string]interface{}) Block {
	b := Block{
		TimeStamp: timestamp,
		Data:      data,
		Nonce:     0,
	}
	b.Hash = b.calculateHash()
	return b
}

func NewGenesisBlock() Block {
	currentTimeStamp := time.Now().Unix()
	genesisBlock := NewBlock(currentTimeStamp, map[string]interface{}{"message": "Genesis Block"})
	genesisBlock.PrevHash = "0"
	return genesisBlock
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

package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"
)

// Block is an individual block
type Block struct {
	Index     int
	Timestamp string
	Hash      string
	PrevHash  string

	// Data to keep record of
	Data []byte
}

// Blockchain type
type Blockchain struct {
	Blocks []Block
	mu     sync.Mutex
}

var bcServer chan []Block

// New blockchain
func New(genesis Block) *Blockchain {
	var blocks []Block
	blocks = append(blocks, genesis)
	return &Blockchain{
		Blocks: blocks,
		mu:     sync.Mutex{},
	}
}

func calculateHash(block Block) string {
	record := string(block.Index) + block.Timestamp + string(block.Data) + block.PrevHash
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// IsBlockValid validates a single block
func (bc *Blockchain) IsBlockValid(newBlock, oldBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}

	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}

	if calculateHash(newBlock) != newBlock.Hash {
		return false
	}

	return true
}

// GenerateBlock generates a new block
func (bc *Blockchain) GenerateBlock(oldBlock Block, data []byte) (Block, error) {
	var newBlock Block
	t := time.Now()
	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.Data = data
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Hash = calculateHash(newBlock)

	return newBlock, nil
}

// ReplaceChain replaces the blockchain with the new blocks
// if the new blocks are longer than the current, meaning that
// the new blocks are more up-to-date.
func (bc *Blockchain) ReplaceChain(newBlocks []Block) {
	if len(newBlocks) > len(bc.Blocks) {
		bc.mu.Lock()
		bc.Blocks = newBlocks
		bc.mu.Unlock()
	}
}

// Append adds a new block to the blockchain
func (bc *Blockchain) Append(data []byte) (Block, error) {
	newBlock, err := bc.GenerateBlock(bc.Blocks[len(bc.Blocks)-1], data)
	if err != nil {
		return newBlock, err
	}

	if bc.IsBlockValid(newBlock, bc.Blocks[len(bc.Blocks)-1]) {
		newBlockchain := append(bc.Blocks, newBlock)
		bc.ReplaceChain(newBlockchain)
		spew.Dump(bc.Blocks)
	}

	return newBlock, nil
}

// GetBlocks gets all blocks in the blockchain
func (bc *Blockchain) GetBlocks() []Block {
	return bc.Blocks
}

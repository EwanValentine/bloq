package blockchain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func setUp(t *testing.T) *Blockchain {
	bloq := New(Block{
		Index:     0,
		Hash:      "",
		PrevHash:  "",
		Data:      []byte(""),
		Timestamp: time.Now().String(),
	})
	assert.Equal(t, 1, len(bloq.GetBlocks()))
	return bloq
}

func TestBlockchainEntry(t *testing.T) {
	bloq := setUp(t)
	test := []byte("test")
	newBlock, err := bloq.Append(test)
	assert.NoError(t, err)
	assert.Equal(t, newBlock.Data, test)
	assert.Equal(t, 2, len(bloq.GetBlocks()))
}

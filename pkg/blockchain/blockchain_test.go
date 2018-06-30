package blockchain

import (
	"errors"
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
	test := []byte(`{ "test": "world" }`)
	newBlock, err := bloq.Append(test)
	assert.NoError(t, err)
	assert.Equal(t, newBlock.Data, test)
	assert.Equal(t, 2, len(bloq.GetBlocks()))
}

func TestSmartContracts(t *testing.T) {

	golden := map[string][]byte{
		"users:>:5":   []byte(`{ "users": 6 }`),
		"name:=:Ewan": []byte(`{ "name": "Ewan" }`),
	}

	notGolden := map[string][]byte{
		"users:>:5":   []byte(`{ "users": 5 }`),
		"name:=:Test": []byte(`{ "name": "not Test" }`),
	}

	badContracts := map[string]error{
		"invalid-contract":    errors.New("Invalid contract"),
		"invalid:??:operator": errors.New("Invalid operator"),
	}

	for key, val := range golden {
		called := false
		handler := func(block Block) error {
			called = true
			return nil
		}

		bloq := setUp(t)
		bloq.AddContract(key, handler)
		bloq.Append(val)
		assert.Equal(t, true, called)
	}

	for key, val := range notGolden {
		called := false
		handler := func(block Block) error {
			called = true
			return nil
		}

		bloq := setUp(t)
		bloq.AddContract(key, handler)
		bloq.Append(val)
		assert.NotEqual(t, true, called)
	}

	for key, val := range badContracts {
		called := false
		handler := func(block Block) error {
			called = true
			return nil
		}

		bloq := setUp(t)
		err := bloq.AddContract(key, handler)
		assert.NotEqual(t, true, called)
		assert.Error(t, err, val.Error())
	}
}

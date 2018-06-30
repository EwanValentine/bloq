package main

import (
	"log"
	"time"

	"github.com/EwanValentine/bloq/internal/api"
	"github.com/EwanValentine/bloq/pkg/blockchain"
	"github.com/joho/godotenv"
)

func main() {

	bloq := blockchain.New(blockchain.Block{
		Index:     0,
		Timestamp: time.Now().String(),
		Hash:      "",
		PrevHash:  "",
		Data:      []byte(`{ "hello": "world" }`),
	})

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	rest := api.NewHTTPAPI(bloq)
	tcp := api.NewTCPAPI(bloq)

	go func() {
		log.Fatal(tcp.Run())
	}()
	log.Fatal(rest.Run())
}

package api

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"net"
	"os"
	"time"

	"github.com/EwanValentine/bloq/pkg/blockchain"

	"github.com/davecgh/go-spew/spew"
)

type TCPAPI struct {
	blockchain bc
	bcServer   chan []blockchain.Block
}

func NewTCPAPI(bc bc) *TCPAPI {
	return &TCPAPI{
		blockchain: bc,
		bcServer:   make(chan []blockchain.Block),
	}
}

func (api *TCPAPI) handleConn(conn net.Conn) {
	defer conn.Close()
	io.WriteString(conn, "Enter a new BPM:")

	scanner := bufio.NewScanner(conn)

	// take in BPM from stdin and add it to blockchain after conducting necessary validation
	go func() {
		for scanner.Scan() {
			_, err := api.blockchain.Append([]byte(scanner.Text()))
			if err != nil {
				log.Println(err)
				continue
			}
			api.bcServer <- api.blockchain.GetBlocks()
			io.WriteString(conn, "\nEnter a new BPM:")
		}
	}()

	go func() {
		for {
			time.Sleep(30 * time.Second)
			output, err := json.Marshal(api.blockchain.GetBlocks())
			if err != nil {
				log.Fatal(err)
			}
			io.WriteString(conn, string(output))
		}
	}()

	for _ = range api.bcServer {
		spew.Dump(api.blockchain.GetBlocks())
	}
}

// Run TCP server
func (api *TCPAPI) Run() error {

	// Start TCP and serve TCP server
	server, err := net.Listen("tcp", ":"+os.Getenv("TCP_ADDR"))
	if err != nil {
		return err
	}
	defer server.Close()

	for {
		conn, err := server.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go api.handleConn(conn)
	}
}

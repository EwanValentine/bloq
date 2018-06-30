package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/EwanValentine/bloq/pkg/blockchain"
	"github.com/gorilla/mux"
)

type message struct {
	Data []byte
}

type bc interface {
	GetBlocks() []blockchain.Block
	Append([]byte) (blockchain.Block, error)
}

// HTTPAPI is a standard REST API interface to the blockchain
type HTTPAPI struct {
	blockchain bc
}

// NewHTTPAPI returns a new HTTPAPI instance, and takes a blockchain
// implementation as its only argument
func NewHTTPAPI(blockchain bc) *HTTPAPI {
	return &HTTPAPI{
		blockchain: blockchain,
	}
}

func respondWithJSON(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
	response, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: Internal Server Error"))
		return
	}
	w.WriteHeader(code)
	w.Write(response)
}

func (api *HTTPAPI) handleGetBlockchain(w http.ResponseWriter, r *http.Request) {
	bytes, err := json.MarshalIndent(api.blockchain.GetBlocks(), "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
}

func (api *HTTPAPI) handleWriteBlock(w http.ResponseWriter, r *http.Request) {
	var m message

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	newBlock, err := api.blockchain.Append(m.Data)
	if err != nil {
		respondWithJSON(w, r, http.StatusInternalServerError, m)
		return
	}

	respondWithJSON(w, r, http.StatusCreated, newBlock)
}

func (api *HTTPAPI) makeMuxRouter() http.Handler {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/", api.handleGetBlockchain).Methods("GET")
	muxRouter.HandleFunc("/", api.handleWriteBlock).Methods("POST")
	return muxRouter
}

// Run http api
func (api *HTTPAPI) Run() error {
	mux := api.makeMuxRouter()
	httpAddr := os.Getenv("HTTP_ADDR")
	log.Println("Listening on ", os.Getenv("HTTP_ADDR"))
	s := &http.Server{
		Addr:           ":" + httpAddr,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

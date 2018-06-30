# Bloq

Bloq is a Blockchain framework written in Go. Largely for my own understanding, probably not for production use. 

## Use

### Blockchain
You create a new blockchain with a 'genesis block' as an argument to the constructor. Then you load your .env config file, and spin up a server of your choice (you can spin up different types of servers, just be sure to config them on different ports):
```golang
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
log.Fatal(rest.Run())
```

`$ ADDR=8080 TCP_ADDR=9000 go run main.go`

### Smart Contracts

```golang
// Dummy code, but you get the gist
func handlerFunc(block Block) error {
    d := map[string]string{}
    if err := json.Marshal(block.Data, &d); err != nil {
        return err
    }
    return makePayment(d["user_a"])
}

func main() {
    ... 
    bloq.AddContract("name:=:Ewan", handlerFunc)
    bloq.AddContract("users:>:5", handlerFunc)
    ...
}
```

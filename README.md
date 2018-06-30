# Bloq

Bloq is a Blockchain framework written in Go. Largely for my own understanding, probably not for production use. 

## Use

## Smart Contracts

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

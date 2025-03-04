# Wisdom Server/Client 

## For starting code
run command:
```
make start
```
wisdom_server and wisdom_client will be created.

Wisdom client will stop working right after random quote will be received. Wisdom server will continue working till it will be stopped manually.


## POW Explained

Honestly I didn't focus a lot on choosing right POW algo. I took Bitcoin idea of having some number of "0" in the beginning of sha256 hash from (task + rand). Practically it takes around 2 seconds to run it on my computer.

For adjusting POW you can use the next way for configuration:

```go
srv := server.NewServer(func(cfg *server.Cfg) {
    cfg.PowTimeout = time.Second * 20 // Time we give for solving
    cfg.PowComplexity = 8 
})
```

## PS

As it's test task I don't focus a lot on code quality including testing part. Also to simplify reviewer work I didn't apply architectural patterns for code organization like Clean Architecture, DDD, etc.

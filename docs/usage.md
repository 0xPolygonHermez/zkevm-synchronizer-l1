# Usage

### Create a `main.go` with next contents: 
```
package main

import (
	"context"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/synchronizer"
)

func main() {
	ctx := context.Background()

	sync, err := synchronizer.NewSynchronizerFromConfigfile(ctx, "./config.toml")
	if err != nil {
		panic(err)
	}
	sync.Sync(true)
}
```

### Create config file
You need a config file to pass the parameters

- Create a `config.toml` with next contents: 
```
[Log]
    Environment = "development" # "production" or "development"
    Level = "debug"
    Outputs = ["stderr"]
[DB]
	Name = "sync"
	User = "test_user"
	Password = "test_password"
	Host = "localhost"
	Port = "5436"
	MaxConns = 10
[Synchronizer]
	SyncInterval = "10s"
	SyncChunkSize = 50000
	# 0 means find out the block of ETROG upgrade
	GenesisBlockNumber = 0
[Etherman]
	L1URL = "http://your.L1node.url"
	Contracts.GlobalExitRootManagerAddr = "0x580bda1e7A0CFAe92Fa7F6c20A3794F169CE3CFb"
	Contracts.RollupManagerAddr = "0x5132A183E9F3CB7C848b0AAC5Ae0c4f0491B7aB2"
	Contracts.ZkEVMAddr = "0x519E42c24163192Dca44CD3fBDCEBF6be9130987"
```

### Setup a Postgres Database 
- Be sure that the config file `DB` section match your database


### Run it
It's better don't store condidential values in the config file, to override it on runtime you could use a environment variable that match `ZKEVM_SYNCL1_`\<section>`_`\<variable>. All of them in uppercase
```
ZKEVM_SYNCL1_ETHERMAN_L1URL="https://mainnet.infura.io/v3/your_api_key" go run main.go
```

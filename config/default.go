package config

// DefaultValues is the default configuration
const DefaultValues = `
[Log]
	Environment = "development" # "production" or "development"
	Level = "info"
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
	SyncChunkSize = 1000
	GenesisBlockNumber = 5157839
	[Synchronizer.L1BlockCheck]
		Enable = true
		L1SafeBlockPoint = "finalized"
		L1SafeBlockOffset = 0
		ForceCheckBeforeStart = true
		PreCheckEnable = true
		L1PreSafeBlockPoint = "safe"
		L1PreSafeBlockOffset = 0
[Etherman]
	L1URL = "http://localhost:8545"
	[Etherman.Contracts]
	GlobalExitRootManagerAddr = "0x2968D6d736178f8FE7393CC33C87f29D9C287e78"
	RollupManagerAddr = "0xE2EF6215aDc132Df6913C8DD16487aBF118d1764"
	ZkEVMAddr = "0x89BA0Ed947a88fe43c22Ae305C0713eC8a7Eb361"
`

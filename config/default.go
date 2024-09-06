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
[SQLDB]
	DriverName = "set_driver. example: sqlite"
	DataSource = "example: file::memory:?cache=shared"
[Synchronizer]
	SyncInterval = "10s"
	SyncChunkSize = 500
	GenesisBlockNumber = 0
	SyncUpToBlock = "latest"
	BlockFinality = "finalized"
	OverrideStorageCheck = false
[Etherman]
	L1URL = "http://localhost:8545"
	ForkIDChunkSize = 100
	L1ChainID = 0
	PararellBlockRequest = false
	[Etherman.Contracts]
		GlobalExitRootManagerAddr = "0x2968D6d736178f8FE7393CC33C87f29D9C287e78"
		RollupManagerAddr = "0xE2EF6215aDc132Df6913C8DD16487aBF118d1764"
		ZkEVMAddr = "0x89BA0Ed947a88fe43c22Ae305C0713eC8a7Eb361"
	[Etherman.Validium]
		Enabled = false
		TrustedSequencerURL = ""
		RetryOnDACErrorInterval = "1m"
		DataSourcePriority = ["trusted", "external"]
		[Etherman.Validium.Translator]
			FullMatchRules = []
		[Etherman.Validium.RateLimit]
			NumRequests = 900
			Interval = "1s"
`

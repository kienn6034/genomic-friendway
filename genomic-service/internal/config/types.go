package config

// Store environment variables
type EnvVariable struct {
	PrivateKey string // env: PRIVATE_KEY
}

type StorageSettings struct {
	// Setting for storage service
}

type TEESettings struct {
	// Setting for TEE service
}

type BlockchainSettings struct {
	RPCURL            string
	GeneNFTAddress    string
	PCSPTokenAddress  string
	ControllerAddress string
}

type WalletSettings struct {
	PrivateKey string
}

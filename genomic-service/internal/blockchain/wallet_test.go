package blockchain

import (
	"genomic-service/internal/config"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWalletSetup(t *testing.T) {
	// Load configuration
	config.LoadEnv("../../.env")
	cfg := config.NewConfig("../config/app.ini")

	// Initialize blockchain service
	service, err := NewBlockchainService(
		cfg.BlockchainSettings.RPCURL,
		cfg.WalletSettings.PrivateKey,
		cfg.BlockchainSettings.ControllerAddress,
	)
	assert.NoError(t, err)
	assert.NotNil(t, service)

	// Test wallet setup
	assert.NotNil(t, service.wallet)
	assert.NotEmpty(t, service.wallet.Address.Hex())

	// Test contract loading
	assert.NotNil(t, service.controller)
	assert.NotNil(t, service.nft)
	assert.NotNil(t, service.token)

	// Test getting wallet balance
	balance, err := service.wallet.GetBalance(service.client)
	assert.NoError(t, err)
	assert.NotNil(t, balance)
	// make sure wallet balance existed
	assert.Greater(t, balance.Cmp(big.NewInt(0)), 0)
}

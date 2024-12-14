package blockchain

import (
	"genomic-service/internal/config"
	"genomic-service/internal/types"
	"math/big"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func setupTestConfig() *config.Config {
	config.LoadEnv("../../.env")
	return config.NewConfig("../config/app.ini")
}

func TestBlockchainService(t *testing.T) {
	// Load configuration
	cfg := setupTestConfig()

	// Initialize blockchain service
	service, _ := NewBlockchainService(
		cfg.BlockchainSettings.RPCURL,
		cfg.WalletSettings.PrivateKey,
		cfg.BlockchainSettings.ControllerAddress,
	)

	// Fetch NFTs
	tokenBalance, err := service.token.BalanceOf(nil, service.wallet.Address)
	assert.NoError(t, err)

	// make sure contract got interacted
	assert.Greater(t, tokenBalance.Cmp(big.NewInt(0)), 0)
}

func TestInitiateDataUpload(t *testing.T) {
	cfg := setupTestConfig()
	service, _ := NewBlockchainService(
		cfg.BlockchainSettings.RPCURL,
		cfg.WalletSettings.PrivateKey,
		cfg.BlockchainSettings.ControllerAddress,
	)

	// random docID
	docID := uuid.New().String()
	sessionID, err := service.InitiateDataUpload(docID)
	assert.NoError(t, err)
	assert.NotEmpty(t, sessionID)

	t.Logf("Session ID: %s", sessionID)
}

func TestProcessAndMint(t *testing.T) {
	cfg := setupTestConfig()
	service, _ := NewBlockchainService(
		cfg.BlockchainSettings.RPCURL,
		cfg.WalletSettings.PrivateKey,
		cfg.BlockchainSettings.ControllerAddress,
	)

	// random docID
	docID := uuid.New().String()
	// upload docs
	sessionID, err := service.InitiateDataUpload(docID)
	assert.NoError(t, err)

	result := &types.ProcessResult{
		SessionID:   sessionID,
		RiskScore:   2,
		ContentHash: "0x1234",
		DocID:       docID,
	}

	err = service.ProcessAndMint(result)
	assert.NoError(t, err)
}

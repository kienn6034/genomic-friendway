package config_test

import (
	"genomic-service/internal/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Make sure correctly load env variable and setup config
func TestConfig(t *testing.T) {
	config.LoadEnv("../../.env")
	cfg := config.NewConfig("app.ini")

	assert.NotNil(t, cfg)

	assert.Greater(t, len(cfg.BlockchainSettings.RPCURL), 0)
	assert.Greater(t, len(cfg.WalletSettings.PrivateKey), 0)
}

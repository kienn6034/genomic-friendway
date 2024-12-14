package config

import (
	"log"

	"github.com/go-ini/ini"
)

// Update Config struct to be more specific
type Config struct {
	StorageSettings    *StorageSettings
	TEESettings        *TEESettings
	BlockchainSettings *BlockchainSettings
	WalletSettings     *WalletSettings
}

// Setup initializes the configuration instance and returns it
func SetupConfigSettings(path string) *Config {
	var err error
	cfg, err := ini.Load(path)
	if err != nil {
		log.Fatalf("setting.Setup, fail to parse config file: %v", err)
	}

	storageSetting := &StorageSettings{}
	teeSetting := &TEESettings{}
	blockchainSetting := &BlockchainSettings{}
	walletSetting := &WalletSettings{}

	mapTo(cfg, "storage", storageSetting)
	mapTo(cfg, "tee", teeSetting)
	mapTo(cfg, "blockchain", blockchainSetting)

	return &Config{
		StorageSettings:    storageSetting,
		TEESettings:        teeSetting,
		BlockchainSettings: blockchainSetting,
		WalletSettings:     walletSetting,
	}
}

// mapTo maps section
func mapTo(cfg *ini.File, section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("Cfg.MapTo %s err: %v", section, err)
	}
}

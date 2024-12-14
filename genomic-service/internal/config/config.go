package config

func NewConfig(path string) *Config {
	cfg := SetupConfigSettings(path)

	cfg.SetupEnvVariable()
	return cfg
}

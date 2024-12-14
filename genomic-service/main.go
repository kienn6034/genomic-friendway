package main

import (
	"genomic-service/internal/config"
	"genomic-service/internal/server"
	"log"
)

var cfg *config.Config

func init() {
	config.LoadEnv(".env")
	cfg = config.SetupConfigSettings("internal/config/app.ini")
	cfg.SetupEnvVariable()
}

func main() {
	srv, err := server.NewServer(cfg)
	if err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}

	if err := srv.Run(); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

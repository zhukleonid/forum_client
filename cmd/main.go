package main

import (
	"log"
	"lzhuk/clients/internal/server"
	"lzhuk/clients/pkg/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	server.StartServer(cfg)
}

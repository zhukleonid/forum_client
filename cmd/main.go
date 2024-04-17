package main

import (
	"log"
	"lzhuk/clients/internal/controller"
	"lzhuk/clients/internal/server"
	"lzhuk/clients/pkg/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	controller.InitPointApi(cfg)
	server.StartServer(cfg)
}

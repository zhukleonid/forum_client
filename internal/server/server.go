package server

import (
	"log"
	"lzhuk/clients/internal/controller"
	"lzhuk/clients/pkg/config"
	"net/http"
)

func StartServer(cfg config.Config) {
	server := controller.Router(cfg)
	err := server.ListenAndServe()
	if err != nil {
		if err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}
}

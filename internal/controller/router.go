package controller

import (
	"lzhuk/clients/pkg/config"
	"net/http"
	"time"
)

func Router(cfg config.Config) *http.Server {
	mux := http.NewServeMux()
	
	mux.HandleFunc("/", startPage)
	mux.HandleFunc("/userd3", homePage)
	mux.HandleFunc("/register", registerPage)
	
	
	fileServer := http.FileServer(http.Dir("ui/css"))
	mux.Handle("/ui/css/", http.StripPrefix("/ui/css/", fileServer))
	s := &http.Server{
		Addr:         cfg.Port,
		ReadTimeout:  time.Duration(cfg.ReadTimeout) * time.Second,  // время ожидания для чтения данных
		WriteTimeout: time.Duration(cfg.WriteTimeout) * time.Second,  // время ожидания для записи данных
		IdleTimeout:  time.Duration(cfg.IdleTimeout) * time.Second, // время простоя
		Handler:      mux,
	}
	return s
}

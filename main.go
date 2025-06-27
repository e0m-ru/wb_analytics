package main

import (
	"log"
	"net/http"

	"github.com/e0m-ru/wb_analitics/api"
	"github.com/e0m-ru/wb_analitics/config"
	"github.com/e0m-ru/wb_analitics/frontend"
)

var (
	cfg = config.Load()
)

func main() {

	go func(cfg *config.Config) {
		router := api.NewRouter()
		log.Printf("API server starting on %s", cfg.APIAddress)
		log.Fatal(http.ListenAndServe(cfg.APIAddress, router))
	}(cfg)

	server := frontend.NewServer(cfg)
	log.Printf("Frontend server starting on http://localhost%s", cfg.FrontendAddress)
	log.Fatal(http.ListenAndServe(cfg.FrontendAddress, server))
}

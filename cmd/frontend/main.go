package main

import (
	"log"
	"net/http"

	"github.com/e0m-ru/wb_analitics/config"
	"github.com/e0m-ru/wb_analitics/frontend"
)

func main() {
	cfg := config.Load()

	server := frontend.NewServer(cfg)

	log.Printf("Frontend server starting on %s", cfg.FrontendAddress)
	log.Fatal(http.ListenAndServe(cfg.FrontendAddress, server))
}

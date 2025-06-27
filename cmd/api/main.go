package main

import (
	"log"
	"net/http"

	"github.com/e0m-ru/wb_analitics/api"
	"github.com/e0m-ru/wb_analitics/config"
)

func main() {
	cfg := config.Load()

	router := api.NewRouter()

	log.Printf("API server starting on %s", cfg.APIAddress)
	log.Fatal(http.ListenAndServe(cfg.APIAddress, router))
}

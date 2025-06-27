package main

import (
	"log"

	"github.com/e0m-ru/wb_analitics/internal/parser"
)

func main() {
	err := parser.ParseProducts("iphone")
	if err != nil {
		log.Fatal(err)
	}
}

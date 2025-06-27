package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/e0m-ru/wb_analitics/config"
	"github.com/e0m-ru/wb_analitics/internal/parser"
	"github.com/e0m-ru/wb_analitics/internal/storage"
)

type ProductResponse struct {
	Products []config.Product `json:"products"`
	Count    int              `json:"count"`
}

func ParseProducts(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Query string `json:"query"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := parser.ParseProducts(request.Query)
	if err != nil {
		log.Printf("Parser error: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "started",
		"query":  request.Query,
	})
}

func GetProducts(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	filters := storage.ProductFilters{
		MinPrice:     parseFloatParam(query, "min_price"),
		MaxPrice:     parseFloatParam(query, "max_price"),
		MinRating:    parseIntParam(query, "min_rating"),
		MaxRating:    parseIntParam(query, "max_rating"),
		MinFeedbacks: parseIntParam(query, "min_feedbacks"),
		MaxFeedbacks: parseIntParam(query, "max_feedbacks"),
	}

	products, err := storage.GetFilteredProducts(filters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := ProductResponse{
		Products: products,
		Count:    len(products),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

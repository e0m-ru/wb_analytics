package api

import (
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter()

	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/parse", ParseProducts).Methods("POST")
	productsRouter := apiRouter.PathPrefix("/products").Subrouter()
	productsRouter.HandleFunc("", GetProducts).Methods("GET")

	return router
}

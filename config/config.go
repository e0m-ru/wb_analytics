package config

import "time"

// Структура для хранения данных о товаре
type Product struct {
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	SalePrice float64 `json:"sale_price"`
	Rating    float64 `json:"rating"`
	Feedbacks int     `json:"feedbacks"`
}

type Config struct {
	APIAddress      string
	FrontendAddress string
	APIBaseURL      string
	DBPath          string
	WBURL           string
	Pages           int
	TimeOut         time.Duration
}

func Load() *Config {
	config := &Config{
		APIAddress:      ":8081",
		FrontendAddress: ":8080",
		APIBaseURL:      "http://localhost:8081",
		WBURL:           "https://search.wb.ru/exactmatch/ru/common/v4/search?",
		Pages:           10, // количество страниц парсинга
		DBPath:          "./data/products.db",
		TimeOut:         time.Millisecond * 100, // чтобы не нервировать wb
	}
	return config
}

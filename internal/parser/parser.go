package parser

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/e0m-ru/wb_analitics/config"
	"github.com/e0m-ru/wb_analitics/internal/storage"
)

var (
	cfg = config.Load()
)

func ParseProducts(query string) error {
	var parsedProducts []config.Product
	queryParams := initQueryParams(query)

	log.Print("save data in database")
	for i := range cfg.Pages {
		queryParams.Set("page", fmt.Sprintf("%d", i))
		fullURL := cfg.WBURL + queryParams.Encode()

		resp, err := http.Get(fullURL)
		if err != nil {
			return fmt.Errorf("error sending request: %v", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("responce reading error: %v", err)
		}

		var result map[string]any
		if err := json.Unmarshal(body, &result); err != nil {
			return fmt.Errorf("parsing JSON error: %v", err)
		}

		data, ok := result["data"].(map[string]any)
		if !ok {
			log.Print("failed to find data on goods in response")
			return nil
		}

		products, ok := data["products"].([]any)
		if !ok || len(products) == 0 {
			log.Printf("no goods found or the pages of issuance ended")
		}

		// Собираем информацию о товарах
		for _, p := range products {
			product := p.(map[string]any)
			name := product["name"].(string)
			price := product["priceU"].(float64) / 100 // Цена в рублях (WB возвращает в копейках)
			salePrice := product["salePriceU"].(float64) / 100
			rating := product["rating"].(float64)
			feedbacks := int(product["feedbacks"].(float64))

			parsedProducts = append(parsedProducts, config.Product{
				Name:      name,
				Price:     price,
				SalePrice: salePrice,
				Rating:    rating,
				Feedbacks: feedbacks,
			})
		}
		time.Sleep(cfg.TimeOut)
	}
	storage.WriteToStorage(parsedProducts, query)
	return nil
}

func initQueryParams(query string) url.Values {
	// URL для запроса к Wildberries API
	params := url.Values{}
	params.Add("query", query)
	params.Add("resultset", "catalog")

	// params.Add("sort", "popular")
	params.Add("locale", "ru")
	params.Add("curr", "rub")
	params.Add("lang", "ru")

	params.Add("appType", "1")
	params.Add("dest", "-1257786")

	return params
}

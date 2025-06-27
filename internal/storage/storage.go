package storage

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/e0m-ru/wb_analitics/config" // сомнительно

	_ "github.com/mattn/go-sqlite3"
)

var (
	db     *sql.DB
	DBPath = config.Load().DBPath
)

type ProductFilters struct {
	MinPrice     float64
	MaxPrice     float64
	MinRating    int
	MaxRating    int
	MinFeedbacks int
	MaxFeedbacks int
}

func InitDB(dataSourceName string) error {
	var err error
	db, err = sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return err
	}
	return db.Ping()
}

func WriteToStorage(parsedProducts []config.Product, query string) error {
	InitDB(DBPath)
	defer db.Close()

	_, err := db.Exec(`DROP TABLE products;
CREATE TABLE products (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    price REAL,
    sale_price REAL,
    rating REAL,
    feedbacks INTEGER,
    search_query TEXT,
    added_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)`)
	if err != nil {
		return fmt.Errorf("error creating table: %v", err)
	}

	stmt, err := db.Prepare(`
INSERT INTO products(name, price, sale_price, rating, feedbacks, search_query)
VALUES(?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return fmt.Errorf("prepearre request error: %v", err)
	}
	defer stmt.Close()

	// Вставляем данные
	for _, p := range parsedProducts {
		_, err := stmt.Exec(p.Name, p.Price, p.SalePrice, p.Rating, p.Feedbacks, query)
		if err != nil {
			log.Printf("Error inserting %s: %v\n", p.Name, err)
			continue
		}
	}
	log.Printf("Успешно сохранено %d товаров в базу\n", len(parsedProducts))
	return nil
}

func GetFilteredProducts(filters ProductFilters) ([]config.Product, error) {
	InitDB(DBPath)
	defer db.Close()
	query := `SELECT name, price, sale_price, rating, feedbacks 
	          FROM products WHERE 1=1`
	args := []any{}

	if filters.MinPrice > 0 {
		query += " AND price >= ?"
		args = append(args, filters.MinPrice)
	}
	if filters.MaxPrice > 0 {
		query += " AND price <= ?"
		args = append(args, filters.MaxPrice)
	}
	if filters.MinFeedbacks > 0 {
		query += " AND feedbacks >= ?"
		args = append(args, filters.MinFeedbacks)
	}
	if filters.MaxFeedbacks > 0 {
		query += " AND feedbacks <= ?"
		args = append(args, filters.MaxFeedbacks)
	}
	if filters.MinRating > 0 {
		query += " AND rating >= ?"
		args = append(args, filters.MinRating)
	}
	if filters.MaxRating > 0 {
		query += " AND rating <= ?"
		args = append(args, filters.MaxRating)
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []config.Product
	for rows.Next() {
		var p config.Product
		err := rows.Scan(&p.Name, &p.Price, &p.SalePrice, &p.Rating, &p.Feedbacks)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

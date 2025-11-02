package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Product struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
	Price    int    `json:"price"`
}

func getProductsHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	query := `
		SELECT p.id, p.name, c.name AS category, p.price
		FROM products p
		JOIN categories c ON p.category_id = c.id
	`
	conditions := []string{}
	args := []interface{}{}
	argIdx := 1

	if category := r.URL.Query().Get("category"); category != "" {
		conditions = append(conditions, fmt.Sprintf("c.name = $%d", argIdx))
		args = append(args, category)
		argIdx++
	}
	if minPrice := r.URL.Query().Get("min_price"); minPrice != "" {
		conditions = append(conditions, fmt.Sprintf("p.price >= $%d", argIdx))
		args = append(args, minPrice)
		argIdx++
	}
	if maxPrice := r.URL.Query().Get("max_price"); maxPrice != "" {
		conditions = append(conditions, fmt.Sprintf("p.price <= $%d", argIdx))
		args = append(args, maxPrice)
		argIdx++
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	sort := r.URL.Query().Get("sort")
	switch sort {
	case "price_asc":
		query += " ORDER BY p.price ASC"
	case "price_desc":
		query += " ORDER BY p.price DESC"
	default:
		query += " ORDER BY p.id ASC"
	}

	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	if limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			query += fmt.Sprintf(" LIMIT %d", limit)
		}
	}
	if offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset > 0 {
			query += fmt.Sprintf(" OFFSET %d", offset)
		}
	}

	startQuery := time.Now()
	rows, err := db.Query(query, args...)
	duration := time.Since(startQuery)
	if err != nil {
		http.Error(w, fmt.Sprintf("DB query error: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Category, &p.Price); err != nil {
			http.Error(w, fmt.Sprintf("Scan error: %v", err), http.StatusInternalServerError)
			return
		}
		products = append(products, p)
	}

	w.Header().Set("X-Query-Time", duration.String())
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)

	log.Printf("Handled /products in %v (query took %v)", time.Since(start), duration)
}

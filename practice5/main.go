package main

import (
	"log"
	"net/http"
)

func main() {
	initDB()
	defer db.Close()

	http.HandleFunc("/products", getProductsHandler)

	log.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

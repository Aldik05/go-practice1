package middleware

import (
	"log"
	"net/http"

	"example.com/practice2/internal/httpx"
)

const apiKeyHeader = "X-API-Key"
const expectedKey = "secret123"

func APIKey(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)

		key := r.Header.Get(apiKeyHeader)
		if key != expectedKey {
			httpx.ErrorJSON(w, http.StatusUnauthorized, "unauthorized")
			return
		}
		next.ServeHTTP(w, r)
	})
}

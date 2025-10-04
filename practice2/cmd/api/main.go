package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"example.com/practice2/internal/httpx"
	"example.com/practice2/internal/middleware"
)

func main() {
	mux := http.NewServeMux()

	mux.Handle("/user", middleware.APIKey(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetUser(w, r)
		case http.MethodPost:
			handlePostUser(w, r)
		default:
			httpx.ErrorJSON(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	})))

	addr := ":8080"
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}

func handleGetUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || idStr == "" {
		httpx.ErrorJSON(w, http.StatusBadRequest, "invalid id")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]int{"user_id": id})
}

type createUserReq struct {
	Name string `json:"name"`
}

func handlePostUser(w http.ResponseWriter, r *http.Request) {
	var req createUserReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrorJSON(w, http.StatusBadRequest, "invalid name")
		return
	}
	if req.Name == "" {
		httpx.ErrorJSON(w, http.StatusBadRequest, "invalid name")
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, map[string]string{"created": req.Name})
}

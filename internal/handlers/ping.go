package handlers

import (
	"encoding/json"
	"net/http"
)

// Ping is simple keep-alive/ping handler
func Ping() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode("All Green")
	}
}

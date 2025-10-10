package csvqlserver

import (
	"encoding/json"
	"net/http"
)

func Error(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(code)
	response := map[string]string{"error": message}
	_ = json.NewEncoder(w).Encode(response)
}

func Send(w http.ResponseWriter, code int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(data)
}

type H = map[string]any

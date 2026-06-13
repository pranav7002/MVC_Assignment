package controller

import (
	"encoding/json"
	"log"
	"net/http"
)

type APIResponse struct {
	Data  any    `json:"data:omitempty"`
	Error string `json:"error:omitempty"`
}

func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("content-type", "application/json")

	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(APIResponse{Data: data}); err != nil {
		log.Println("Error encoding response!!")
	}
}	

func WriteError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("content-type", "application/json")

	w.WriteHeader(status)
	json.NewEncoder(w).Encode(APIResponse{Error: message})
}
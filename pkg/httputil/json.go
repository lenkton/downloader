package httputil

import (
	"encoding/json"
	"log"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, object any, code int) {
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(object)
	if err != nil {
		log.Printf("ERROR: encoding json: %v\n", err)
	}
}

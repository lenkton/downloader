package httputil

import (
	"context"
	"encoding/json"
	"io"
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

func ParseJSON[T any](r io.Reader) (*T, error) {
	var t T
	err := json.NewDecoder(r).Decode(&t)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func WithJSONBody[T any](next http.Handler, key any) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dto, err := ParseJSON[T](r.Body)
		if err != nil {
			log.Printf("ERROR: decoding json: %v\n", err)
			// TODO: add a file with default error messages
			http.Error(w, `{"error":"malformed request body"}`, http.StatusUnprocessableEntity)
			return
		}
		ctx := context.WithValue(r.Context(), key, dto)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

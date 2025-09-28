package middleware

import (
	"log"
	"net/http"
)

type SpyWriter struct {
	http.ResponseWriter

	code int
}

func (spy *SpyWriter) WriteHeader(code int) {
	spy.code = code
	spy.ResponseWriter.WriteHeader(code)
}

func WithLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("INFO: start processing %v %v\n", r.Method, r.URL)
		spyWriter := &SpyWriter{ResponseWriter: w}
		next.ServeHTTP(spyWriter, r)
		log.Printf("INFO: responded to %v %v with %v\n", r.Method, r.URL, spyWriter.code)
	})
}

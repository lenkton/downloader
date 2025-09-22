package main

import (
	"log"
	"net/http"
)

func main() {
	server := &http.Server{}
	log.Fatalf("%v\n", server.ListenAndServe())
}

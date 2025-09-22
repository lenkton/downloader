package main

import (
	"downloader/pkg/services/task"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	service := task.NewService()
	mux.HandleFunc("POST /tasks", service.HandleCreateTask)

	handler := mux

	server := &http.Server{Addr: ":8080", Handler: handler}

	log.Println("INFO: starting server on :8080")
	log.Fatalf("%v\n", server.ListenAndServe())
}

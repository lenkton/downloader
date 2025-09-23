package main

import (
	"downloader/pkg/httputils/middlewares"
	"downloader/pkg/services/task"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	service := task.NewService()
	mux.HandleFunc("POST /tasks", service.HandleCreateTask)
	mux.HandleFunc("GET /tasks/{task_id}", service.HandleGetTask)

	handler := middlewares.WithLogger(mux)

	server := &http.Server{Addr: ":8080", Handler: handler}

	log.Println("INFO: starting server on :8080")
	log.Fatalf("%v\n", server.ListenAndServe())
}

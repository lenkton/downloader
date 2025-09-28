package main

import (
	"log"
	"net/http"

	"github.com/lenkton/downloader/pkg/httputil/middleware"
	"github.com/lenkton/downloader/pkg/services/task"
)

func main() {
	mux := http.NewServeMux()

	service := task.NewService()
	service.EnsureDownloadDir()
	mux.HandleFunc("POST /tasks", service.HandleCreateTask)
	mux.HandleFunc("GET /tasks/{task_id}", service.HandleGetTask)

	handler := middleware.WithLogger(mux)

	server := &http.Server{Addr: ":8080", Handler: handler}

	log.Println("INFO: starting server on :8080")
	// TODO: add graceful shutdown
	log.Fatalf("%v\n", server.ListenAndServe())
}

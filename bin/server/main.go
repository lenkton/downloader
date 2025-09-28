package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

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
	go runServer(server)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("INFO: stopping")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)
	if err != nil {
		log.Printf("ERROR: shutting down server: %v\n", err)
	}
	log.Println("INFO: server shut down")
}

func runServer(server *http.Server) {
	err := server.ListenAndServe()
	if err == http.ErrServerClosed {
		log.Println("INFO: server closed")
		return
	}
	// err is always non-nil
	log.Printf("ERROR: server stopped: %v\n", err)
}

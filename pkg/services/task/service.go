package task

import (
	"fmt"
	"log"
	"net/http"

	"github.com/lenkton/downloader/pkg/fileutil"
	"github.com/lenkton/downloader/pkg/httputil"
)

type Service struct {
	tasksStorage *storage
	processor    *processor
}

func NewService() *Service {
	return &Service{tasksStorage: NewStorage(), processor: newProcessor()}
}

const DownloadDir = "./downloads/"

func (s *Service) EnsureDownloadDir() error {
	err := fileutil.EnsureDir(DownloadDir)
	if err != nil {
		return fmt.Errorf("ensure download dir: %v", err)
	}
	return nil
}

type taskRequestDTO struct {
	Links []string `json:"links"`
}

type taskRequestDTOContextKey struct{}

func (s *Service) HandleCreateTask(w http.ResponseWriter, r *http.Request) {
	var handler http.Handler = http.HandlerFunc(s.handleCreateTask)
	handler = httputil.WithJSONBody[taskRequestDTO](handler, taskRequestDTOContextKey{})

	handler.ServeHTTP(w, r)
}

func (s *Service) handleCreateTask(w http.ResponseWriter, r *http.Request) {
	tr := r.Context().Value(taskRequestDTOContextKey{}).(*taskRequestDTO)

	task, err := s.createTask(tr.Links)
	if err != nil {
		log.Printf("ERROR: creating task: %v\n", err)
		// TODO: differentiate the errors
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	httputil.WriteJSON(w, task.AsJSON(), http.StatusAccepted)
}

type taskIDContextKey struct{}

func (s *Service) HandleGetTask(w http.ResponseWriter, r *http.Request) {
	var handler http.Handler = http.HandlerFunc(s.handleGetTask)
	handler = httputil.WithPathIntID(handler, "task_id", taskIDContextKey{})

	handler.ServeHTTP(w, r)
}

func (s *Service) handleGetTask(w http.ResponseWriter, r *http.Request) {
	taskID := r.Context().Value(taskIDContextKey{}).(int)

	task, found := s.tasksStorage.Find(taskID)
	if !found {
		log.Printf("WARN: can't find task with id %v\n", taskID)
		http.Error(w, `{"error":"task not found"}`, http.StatusNotFound)
		return
	}

	httputil.WriteJSON(w, task.AsJSON(), http.StatusOK)
}

// TODO: add validation
// TODO: add the actual logic
func (s *Service) createTask(links []string) (*task, error) {
	t := newTask(links)

	s.tasksStorage.Save(t)
	s.processor.Start(t)

	return t, nil
}

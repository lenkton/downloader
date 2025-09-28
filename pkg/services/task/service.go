package task

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

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

func (s *Service) HandleCreateTask(w http.ResponseWriter, r *http.Request) {
	tr := &taskRequestDTO{}
	err := json.NewDecoder(r.Body).Decode(tr)
	if err != nil {
		log.Printf("ERROR: decoding json: %v\n", err)
		http.Error(w, `{"error":"malformed request body"}`, http.StatusUnprocessableEntity)
		return
	}

	task, err := s.createTask(tr.Links)
	if err != nil {
		log.Printf("ERROR: creating task: %v\n", err)
		// TODO: differentiate the errors
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	httputil.WriteJSON(w, task.AsJSON(), http.StatusAccepted)
}

func (s *Service) HandleGetTask(w http.ResponseWriter, r *http.Request) {
	taskIDStr := r.PathValue("task_id")
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		log.Printf("WARN: decoding task_id: %v\n", err)
		http.Error(w, `{"error":"task not found"}`, http.StatusNotFound)
		return
	}
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

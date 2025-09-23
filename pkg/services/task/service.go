package task

import (
	"encoding/json"
	"log"
	"net/http"
)

type Service struct {
	tasksStorage *storage
}

func NewService() *Service {
	return &Service{tasksStorage: NewStorage()}
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

	w.WriteHeader(http.StatusAccepted)
	err = json.NewEncoder(w).Encode(task.AsJSON())
	if err != nil {
		log.Printf("ERROR: encoding json: %v\n", err)
	}
}

// TODO: add validation
// TODO: add the actual logic
func (s *Service) createTask(links []string) (*task, error) {
	t := &task{links: links}

	s.tasksStorage.Save(t)

	return t, nil
}

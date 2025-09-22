package task

import "net/http"

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) HandleCreateTask(http.ResponseWriter, *http.Request) {
}

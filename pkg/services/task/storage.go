package task

type storage struct {
	// TODO: lock this by the mutex
	// NOTE: maybe we should make this UUID
	nextTaskID int
	tasks      map[int]*task
}

func NewStorage() *storage {
	return &storage{tasks: make(map[int]*task), nextTaskID: 1}
}

// also sets the task's ID
// WARN: not thread safe
func (s *storage) Save(t *task) {
	id := s.nextTaskID

	t.id = id
	s.tasks[id] = t

	s.nextTaskID++
}

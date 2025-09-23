package task

type status string

const (
	created  status = "Created"
	started  status = "Started"
	finished status = "Finished"
)

// TODO: add status
type task struct {
	id     int
	links  []string
	status status
}

type taskDTO struct {
	ID     int    `json:"id"`
	Status status `json:"status"`
}

func newTask(links []string) *task {
	return &task{links: links, status: created}
}

func (t *task) AsJSON() taskDTO {
	return taskDTO{
		ID:     t.id,
		Status: t.status,
	}
}

package task

// TODO: add status
type task struct {
	id    int
	links []string
}

type taskDTO struct {
	ID int `json:"id"`
}

func (t *task) AsJSON() taskDTO {
	return taskDTO{
		ID: t.id,
	}
}

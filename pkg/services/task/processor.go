package task

type processor struct {
}

func newProcessor() *processor {
	return &processor{}
}

// TODO: wait for goroutines to finish
func (p *processor) Start(t *task) {
	t.status = started
	go func() {
		t.status = finished
	}()
}

package bus

type task struct {
	Queue       string
	ContentType string
	Body        []byte
	success     chan bool
}

func makeTask(queue string, contentType string, body []byte) task {
	return task{
		Queue:       queue,
		ContentType: contentType,
		Body:        body,
		success:     make(chan bool),
	}
}

func (t *task) Success() {
	t.success <- true
}

func (t *task) Failure() {
	t.success <- false
}

func (t *task) IsSuccess() bool {
	return <-t.success
}

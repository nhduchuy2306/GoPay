package webhook

type Worker interface {
}

type worker struct {
}

func NewWorker() Worker {
	return &worker{}
}

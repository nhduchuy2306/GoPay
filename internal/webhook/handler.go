package webhook

type Handler struct {
	worker Worker
}

func NewHandler(worker Worker) *Handler {
	return &Handler{worker: worker}
}

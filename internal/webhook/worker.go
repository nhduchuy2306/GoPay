package webhook

import (
	"bytes"
	"io"
	"net/http"
	"time"
)

type Worker interface {
	Deliver(endpoint Endpoint, payload []byte) (int, []byte, error)
}

type worker struct {
	client *http.Client
}

func NewWorker() Worker {
	return &worker{client: &http.Client{Timeout: 10 * time.Second}}
}

func (w *worker) Deliver(endpoint Endpoint, payload []byte) (int, []byte, error) {
	req, err := http.NewRequest(http.MethodPost, endpoint.URL, bytes.NewReader(payload))
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Webhook-Secret", endpoint.Secret)

	resp, err := w.client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		return resp.StatusCode, body, http.ErrAbortHandler
	}
	return resp.StatusCode, body, nil
}

var _ = NewWorker

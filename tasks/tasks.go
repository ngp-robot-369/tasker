package tasks

import (
	"net/http"

	"github.com/google/uuid"
)

type (
	TaskRequest struct {
		Method  string            `json:"method"`
		URL     string            `json:"url"`
		Headers map[string]string `json:"headers"`
	}
	TaskStatus struct {
		ID              taskID      `json:"id"`
		Status          taskStatus  `json:"status"`
		HTTPStatusCode  int         `json:"httpStatusCode"`
		ResponseHeaders http.Header `json:"responseHeaders"`
		ContentLength   int64       `json:"contentLength"`
	}

	taskID     = string
	taskStatus = string
)

const (
	statusInProgress = taskStatus("in_process")
	statusDone       = taskStatus("done")
	statusError      = taskStatus("error")
	statusNew        = taskStatus("new") // TODO: what the hell is new status?
)

func NewTask(t TaskRequest) taskID {
	id := scheduleTask(t)
	return id
}
func GetStatus(id taskID) *TaskStatus {
	return Storage().GetStatus(id)
}

func scheduleTask(t TaskRequest) string {
	uuid, _ := uuid.NewUUID()
	status := TaskStatus{
		ID:     uuid.String(),
		Status: statusInProgress,
	}
	Storage().PutStatus(status)

	// No scheduling, just perform each task asynchronously
	go executeTask(t, status.ID)
	return status.ID
}

func executeTask(t TaskRequest, id taskID) {
	var (
		new = TaskStatus{
			ID:     id,
			Status: statusDone,
		}
		req, _ = http.NewRequest(t.Method, t.URL, nil)
		client = http.DefaultClient
	)
	for key, val := range t.Headers {
		req.Header.Add(key, val)
	}

	resp, err := client.Do(req)
	if err != nil {
		new.Status = statusError
	}
	if resp != nil {
		defer resp.Body.Close()
		new.HTTPStatusCode = resp.StatusCode
		new.ResponseHeaders = resp.Header
		new.ContentLength = resp.ContentLength
	}
	Storage().PutStatus(new)
}

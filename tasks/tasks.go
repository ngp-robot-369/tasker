package tasks

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
)

type (
	taskID      = string
	status      = string
	httpHeaders http.Header // New type created for converters adding

	TaskRequest struct {
		Method  string            `json:"method"`
		URL     string            `json:"url"`
		Headers map[string]string `json:"headers"`
	}
	TaskStatus struct {
		ID              taskID      `json:"id"`
		Status          status      `json:"status"`
		HTTPStatusCode  int         `json:"httpStatusCode"`
		ResponseHeaders httpHeaders `json:"responseHeaders"`
		ContentLength   int64       `json:"contentLength"`
	}
	IStatusStorage interface {
		PutStatus(ctx context.Context, status TaskStatus) error
		GetStatus(ctx context.Context, id string) (*TaskStatus, error)
	}

	job struct {
		taskID
		TaskRequest
	}
	taskerService struct {
		storage IStatusStorage
		client  http.Client

		poolJobs chan job
		poolSize int
		poolWG   sync.WaitGroup
		poolStop chan struct{}
	}
)

const (
	WORKERS_POOL_SIZE = 4

	statusInProgress = status("in_process")
	statusDone       = status("done")
	statusError      = status("error")
	statusNew        = status("new") // TODO: what the hell is new status?
)

func NewTaskerService(storage IStatusStorage) *taskerService {
	ts := &taskerService{
		storage:  storage,
		client:   http.Client{Timeout: time.Second * 5},
		poolJobs: make(chan job, 100),
		poolSize: WORKERS_POOL_SIZE,
		poolWG:   sync.WaitGroup{},
		poolStop: make(chan struct{}),
	}
	for range ts.poolSize {
		go ts.startWorker()
	}
	return ts
}

func (ts *taskerService) Shutdown(ctx context.Context) error {
	close(ts.poolStop)
	ts.poolWG.Wait()
	return nil
}

func (ts *taskerService) NewTask(request TaskRequest) (taskID, error) {
	var (
		uuid, _ = uuid.NewUUID()
		id      = uuid.String()
		status  = TaskStatus{
			ID:              id,
			Status:          statusInProgress,
			HTTPStatusCode:  0,
			ResponseHeaders: map[string][]string{},
			ContentLength:   0,
		}
	)
	err := ts.storage.PutStatus(context.TODO(), status)
	if err != nil {
		return "", err
	}

	ts.poolJobs <- job{
		taskID:      id,
		TaskRequest: request,
	}
	return id, nil
}
func (ts *taskerService) GetStatus(id taskID) (*TaskStatus, error) {
	return ts.storage.GetStatus(context.TODO(), id)
}

func (ts *taskerService) startWorker() {
	ts.poolWG.Add(1)
	defer ts.poolWG.Done()

	for {
		select {
		case job := <-ts.poolJobs:
			err := ts.executeJob(job)
			if err != nil {
				log.Printf("Failed to execute job: %v", err)
			}

		case <-ts.poolStop:
			return
		}
	}
}
func (ts *taskerService) executeJob(w job) error {
	var (
		st = TaskStatus{
			ID:              w.taskID,
			Status:          statusDone,
			HTTPStatusCode:  0,
			ResponseHeaders: map[string][]string{},
			ContentLength:   0,
		}
		req, _ = http.NewRequest(w.Method, w.URL, nil)
	)
	for key, val := range w.Headers {
		req.Header.Add(key, val)
	}

	resp, err := ts.client.Do(req)
	if err != nil {
		st.Status = statusError
	}
	if resp != nil {
		defer resp.Body.Close()
		st.HTTPStatusCode = resp.StatusCode
		st.ResponseHeaders = httpHeaders(resp.Header)
		st.ContentLength = resp.ContentLength
	}
	err = ts.storage.PutStatus(context.TODO(), st)
	if err != nil {
		return err
	}
	return nil
}

package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"tasker/tasks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func MakeHandlersMock() Handlers {
	return Handlers{
		Tasker: tasks.NewTaskerService(tasks.MakeStorageRam()),
	}
}

func TestTaskerAPI(t *testing.T) {

	handler := MakeHandlersMock()

	t.Run("Invalid create task request", func(t *testing.T) {
		var (
			body   = strings.NewReader("invalid json")
			req, _ = http.NewRequest("POST", "/task", body)
			rr     = httptest.NewRecorder()
		)
		handler.serveCreateTask(rr, req)
		assert.Equal(t, rr.Code, http.StatusBadRequest)
		assert.Equal(t, rr.Body.String(), "Invalid request body\n")
	})

	t.Run("Request task with bad id", func(t *testing.T) {
		var (
			req, _ = http.NewRequest("GET", "/task/bad_id", nil)
			rr     = httptest.NewRecorder()
		)
		handler.serveGetTaskStatus(rr, req)
		assert.NotEqual(t, rr.Code, http.StatusOK)
	})

	t.Run("Valid create and get task", func(t *testing.T) {
		var (
			reqBody = strings.NewReader(`{"method": "GET", "url": "http://google.com"}`)
			req, _  = http.NewRequest("POST", "/task", reqBody)
			id      = ""
		)

		// Create task
		{
			rr := httptest.NewRecorder()
			handler.serveCreateTask(rr, req)
			assert.Equal(t, http.StatusOK, rr.Code)

			var response map[string]string
			assert.NoError(t, json.Unmarshal(rr.Body.Bytes(), &response))
			id = response["id"]
			assert.NotEmpty(t, id)
		}

		// Get initial status
		req, _ = http.NewRequest("GET", "/task/"+id, nil)
		req.SetPathValue("id", id)
		{
			rr := httptest.NewRecorder()
			handler.serveGetTaskStatus(rr, req)
			assert.Equal(t, http.StatusOK, rr.Code)

			var resp tasks.TaskStatus
			assert.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
			assert.Equal(t, "in_process", resp.Status)
		}

		// Give some time and get final status
		time.Sleep(time.Second * 5)
		{
			rr := httptest.NewRecorder()
			handler.serveGetTaskStatus(rr, req)
			assert.Equal(t, http.StatusOK, rr.Code)

			var resp tasks.TaskStatus
			assert.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
			assert.Equal(t, "done", resp.Status)
		}
	})
}

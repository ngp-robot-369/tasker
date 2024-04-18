package tasks

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatus(t *testing.T) {
	s := Storage()
	status := &TaskStatus{
		ID:             "good_id",
		Status:         "done",
		HTTPStatusCode: 500,
		ResponseHeaders: http.Header{
			"header A": {"123", "456", "789"},
			"header B": {"456"},
			"header C": {"789"},
		},
		ContentLength: 9000000,
	}

	t.Run("Get non-existing status from storage", func(t *testing.T) {
		storedStatus := s.GetStatus("bad_id")
		assert.Nil(t, storedStatus)
	})

	t.Run("Put and get status from storage", func(t *testing.T) {
		s.PutStatus(*status)
		storedStatus := s.GetStatus("good_id")
		assert.Equal(t, storedStatus, status)
	})

	t.Run("Update status in storage", func(t *testing.T) {
		status.ContentLength *= 2
		s.PutStatus(*status)
		updatedStatus := s.GetStatus("good_id")
		assert.Equal(t, updatedStatus, status)
	})
}

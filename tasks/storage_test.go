package tasks

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStorageRam(t *testing.T) {
	ram := MakeStorageRam()
	testStorage(t, ram)
}

func TestStoragePG(t *testing.T) {
	pg := MakeStoragePG()
	testStorage(t, pg)
}

func testStorage(t *testing.T, s IStatusStorage) {
	ctx := context.TODO()
	ref := &TaskStatus{
		ID:             "good_id",
		Status:         "done",
		HTTPStatusCode: 500,
		ResponseHeaders: httpHeaders{
			"header A": {"123", "456", "789"},
			"header B": {"456"},
			"header C": {"789"},
		},
		ContentLength: 9000000,
	}

	t.Run("Get non-existing status from storage", func(t *testing.T) {
		status, err := s.GetStatus(ctx, "bad_id")
		assert.Error(t, err)
		assert.Nil(t, status)
	})

	t.Run("Put and get status from storage", func(t *testing.T) {
		err := s.PutStatus(ctx, *ref)
		assert.NoError(t, err)

		status, err := s.GetStatus(ctx, "good_id")
		assert.NoError(t, err)
		assert.Equal(t, ref, status)
	})

	t.Run("Update status in storage", func(t *testing.T) {
		ref.ContentLength *= 2
		err := s.PutStatus(ctx, *ref)
		assert.NoError(t, err)

		status, err := s.GetStatus(ctx, "good_id")
		assert.NoError(t, err)
		assert.Equal(t, ref, status)
	})
}

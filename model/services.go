package model

import (
	"context"
	"tasker/tasks"
)

type IBaseService interface {
	Shutdown(context.Context) error
}

type ITaskerService interface {
	IBaseService

	NewTask(request tasks.TaskRequest) (string, error)
	GetStatus(id string) (*tasks.TaskStatus, error)
}

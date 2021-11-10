package tasks

import "context"

type TaskInterface interface {
	CreateFunction(parameters map[string]string) (func(ctx context.Context) error, error)
	GetTitle() string
}

type TaskStore struct {
	MapOfTask map[string]TaskInterface
}

func NewTaskStore(tasks []TaskInterface) *TaskStore {
	mapOfTask := map[string]TaskInterface{}
	for _, task := range tasks {
		mapOfTask[task.GetTitle()] = task
	}
	return &TaskStore{
		MapOfTask: mapOfTask,
	}
}

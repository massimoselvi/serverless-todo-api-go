package database

import (
	"github.com/massimoselvi/serverless-todo-api-go/server"
)

// ToDoRepo is an interface for database actions
type ToDoRepo interface {
	Get(id string) (*server.ToDo, error)
	GetAll() ([]server.ToDo, error)
	Save(todo *server.ToDo) error
	Delete(id string) error
}

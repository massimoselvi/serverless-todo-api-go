package handlers_test

import (
	"github.com/massimoselvi/serverless-todo-api-go/server"
)

// ClientMock is used to mock a client that uses makes call to DynamoDBAPI
type RepoMock struct {
	GetFn         func(string) (*server.ToDo, error)
	GetAllFn      func() ([]server.ToDo, error)
	SaveFn        func(todo *server.ToDo) error
	DeleteFn      func(string) error
	GetInvoked    bool
	GetAllInvoked bool
	SaveInvoked   bool
	DeleteInvoked bool
}

// Get returns a ToDo by its ID
func (m *RepoMock) Get(id string) (*server.ToDo, error) {
	m.GetInvoked = true
	return m.GetFn(id)
}

// GetAll returns all ToDos
func (m *RepoMock) GetAll() ([]server.ToDo, error) {
	m.GetAllInvoked = true
	return m.GetAllFn()
}

// Save creates or updates a ToDo
func (m *RepoMock) Save(todo *server.ToDo) error {
	m.SaveInvoked = true
	return m.SaveFn(todo)
}

// Delete permanently removes a ToDo
func (m *RepoMock) Delete(id string) error {
	m.DeleteInvoked = true
	return m.DeleteFn(id)
}

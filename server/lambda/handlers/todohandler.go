package handlers

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/massimoselvi/serverless-todo-api-go/server"
	"github.com/massimoselvi/serverless-todo-api-go/server/database"
	"github.com/pkg/errors"
)

// ToDoHandler provides a handle method to handle incoming AWS API Gateway request
type ToDoHandler struct {
	repo database.ToDoRepo
}

// NewToDoHandler creates a new ToDo handler
func NewToDoHandler(repo database.ToDoRepo) *ToDoHandler {
	return &ToDoHandler{
		repo: repo,
	}
}

// Handle handles a request from AWS API Gateway and returns a response
func (h *ToDoHandler) Handle(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	switch req.HTTPMethod {
	case "GET":
		return h.get(req)
	case "POST":
		return h.post(req)
	case "PUT":
		return h.put(req)
	case "DELETE":
		return h.delete(req)
	default:
		return CreateErrorResponse(ErrMethodNotAllowed)
	}
}

func (h *ToDoHandler) get(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	if id, ok := req.PathParameters["id"]; ok {
		return h.getOne(id)
	}

	return h.getAll()
}

func (h *ToDoHandler) getOne(id string) (events.APIGatewayProxyResponse, error) {

	todo, err := h.repo.Get(id)
	if err != nil {
		return CreateErrorResponse(ErrInternal)
	}

	if todo == nil {
		return CreateErrorResponse(ErrNotFound)
	}

	return CreateOKResponse(todo)

}

func (h *ToDoHandler) getAll() (events.APIGatewayProxyResponse, error) {

	todos, err := h.repo.GetAll()
	if err != nil {
		return CreateErrorResponse(ErrInternal)
	}

	return CreateOKResponse(todos)

}

func (h *ToDoHandler) post(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	todo, err := parseToDo(req.Body)
	if err != nil {
		return CreateErrorResponse(ErrInternal)
	}

	if todo.ID != "" {
		return CreateErrorResponse(errors.Wrap(ErrBadRequest, "ID must be empty"))
	}

	err = h.repo.Save(&todo)
	if err != nil {
		return CreateErrorResponse(ErrInternal)
	}
	return CreateOKResponse(todo)
}

func (h *ToDoHandler) put(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	id, ok := req.PathParameters["id"]
	if !ok {
		return CreateErrorResponse(errors.Wrap(ErrBadRequest, "ID is required"))
	}

	todo, err := parseToDo(req.Body)
	if err != nil {
		return CreateErrorResponse(ErrInternal)
	}

	if id != todo.ID {
		return CreateErrorResponse(errors.Wrap(ErrBadRequest, "ID in body does not match ID in path"))
	}

	if t, err := h.repo.Get(id); err != nil {
		return CreateErrorResponse(ErrInternal)
	} else if t == nil {
		return CreateErrorResponse(ErrNotFound)
	}

	err = h.repo.Save(&todo)
	if err != nil {
		return CreateErrorResponse(ErrInternal)
	}
	return CreateOKResponse(todo)
}

func (h *ToDoHandler) delete(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	id, ok := req.PathParameters["id"]

	if !ok {
		return CreateErrorResponse(errors.Wrap(ErrBadRequest, "ID is required"))
	}

	t, err := h.repo.Get(id)
	if err != nil {
		return CreateErrorResponse(ErrInternal)
	}

	if t == nil {
		return CreateErrorResponse(ErrNotFound)
	}

	if err := h.repo.Delete(id); err != nil {
		return CreateErrorResponse(ErrInternal)
	}

	return CreateOKResponse("")

}

func parseToDo(body string) (server.ToDo, error) {
	var t server.ToDo
	err := json.Unmarshal([]byte(body), &t)
	return t, err
}

package dynamodb

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/massimoselvi/serverless-todo-api-go/server"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

const todosTableName = "todos"

// ToDoRepo represents a boltdb repository for managing todos
type ToDoRepo struct {
	db dynamodbiface.DynamoDBAPI
}

// NewToDoRepo returns a new ToDo repository using the given bolt database. It also creates the ToDos
// bucket if it is not yet created on disk.
func NewToDoRepo(db dynamodbiface.DynamoDBAPI) *ToDoRepo {
	return &ToDoRepo{db}
}

// Get returns a ToDo by its ID
func (r *ToDoRepo) Get(id string) (*server.ToDo, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(todosTableName),
		Key:       mapID(id),
	}

	result, err := r.db.GetItem(input)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not get ToDo %s from database", id)
	}

	t := &server.ToDo{}

	err = dynamodbattribute.UnmarshalMap(result.Item, t)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not unmarshal ToDo %s", id)
	}

	if t.ID == "" {
		return nil, nil
	}

	return t, nil
}

// GetAll returns all ToDos
func (r *ToDoRepo) GetAll() ([]server.ToDo, error) {

	input := &dynamodb.ScanInput{
		TableName: aws.String(todosTableName),
	}

	result, err := r.db.Scan(input)
	if err != nil {
		return nil, errors.Wrap(err, "Could not get ToDos from database")
	}

	t := []server.ToDo{}

	// Unmarshal the Items field in the result value to the Item Go type.
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &t)
	if err != nil {
		return nil, errors.Wrap(err, "Could not unmarshal ToDos")
	}

	return t, nil
}

// Save creates or updates a ToDo
func (r *ToDoRepo) Save(todo *server.ToDo) error {

	if todo.ID == "" {
		todo.ID = uuid.NewV4().String()
	}

	todo.ModTime = time.Now()

	t, err := dynamodbattribute.MarshalMap(todo)
	if err != nil {
		return errors.Wrapf(err, "Could not unmarshal ToDo %s", todo.ID)
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(todosTableName),
		Item:      t,
	}

	if _, err := r.db.PutItem(input); err != nil {
		return errors.Wrapf(err, "Could not save ToDo %s to database", todo.ID)
	}

	return nil
}

// Delete permanently removes a ToDo
func (r *ToDoRepo) Delete(id string) error {

	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(todosTableName),
		Key:       mapID(id),
	}

	if _, err := r.db.DeleteItem(input); err != nil {
		return errors.Wrapf(err, "Could not delete ToDo %s to database", id)
	}

	return nil
}

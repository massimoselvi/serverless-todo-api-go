package dynamodb_test

import (
	"errors"
	"testing"
	"time"

	awsdynamodb "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/massimoselvi/serverless-todo-api-go/server"
	"github.com/massimoselvi/serverless-todo-api-go/server/database/dynamodb"
	uuid "github.com/satori/go.uuid"
)

const testUUID = "a8a43435-20d8-4af2-8f94-f504aff2c6f3"

func TestToDoRepo(t *testing.T) {
	t.Run("GetToDoFound", testGetToDoFound)
	t.Run("GetToDoNotFound", testGetToDoNotFound)
	t.Run("GetToDoError", testGetToDoError)
	t.Run("GetAllToDos", testGetAllToDos)
	t.Run("GetAllToDosError", testGetAllToDosError)
	t.Run("CreateToDo", testCreateToDo)
	t.Run("CreateToDoError", testCreateToDoError)
	t.Run("UpdateToDo", testUpdateToDo)
	t.Run("DeleteToDo", testDeleteToDo)
	t.Run("DeleteToDoError", testDeleteToDoError)
}

func testGetToDoFound(t *testing.T) {

	m := &ClientMock{}

	m.GetItemFn = func(*awsdynamodb.GetItemInput) (*awsdynamodb.GetItemOutput, error) {

		toDo := &server.ToDo{
			ID:      testUUID,
			Title:   "Test ToDo",
			ModTime: time.Now(),
		}

		item, err := dynamodbattribute.MarshalMap(toDo)
		if err != nil {
			t.Fatal(err)
		}

		out := &awsdynamodb.GetItemOutput{
			Item: item,
		}

		return out, nil
	}

	repo := dynamodb.NewToDoRepo(m)

	toDo, err := repo.Get(testUUID)
	if err != nil {
		t.Fatal(err)
	}

	if toDo == nil {
		t.Fatal("Expected ToDo have a value")
	}

	if !m.GetItemInvoked {
		t.Fatal("GetItem not invoked")
	}

}

func testGetToDoNotFound(t *testing.T) {

	m := &ClientMock{}

	m.GetItemFn = func(*awsdynamodb.GetItemInput) (*awsdynamodb.GetItemOutput, error) {
		return &awsdynamodb.GetItemOutput{
			Item: make(map[string]*awsdynamodb.AttributeValue),
		}, nil
	}

	repo := dynamodb.NewToDoRepo(m)

	toDo, err := repo.Get(testUUID)
	if err != nil {
		t.Fatal(err)
	}

	if toDo != nil {
		t.Fatal("Expected ToDo to be nil")
	}

	if !m.GetItemInvoked {
		t.Fatal("GetItem not invoked")
	}

}

func testGetToDoError(t *testing.T) {

	m := &ClientMock{}

	m.GetItemFn = func(*awsdynamodb.GetItemInput) (*awsdynamodb.GetItemOutput, error) {
		return nil, errors.New("DB Error")
	}

	repo := dynamodb.NewToDoRepo(m)

	_, err := repo.Get(testUUID)
	if err == nil {
		t.Fatal("Expected Error")
	}

	if !m.GetItemInvoked {
		t.Fatal("GetItem not invoked")
	}

}

func testGetAllToDos(t *testing.T) {

	m := &ClientMock{}

	m.ScanFn = func(*awsdynamodb.ScanInput) (*awsdynamodb.ScanOutput, error) {

		item1, err := dynamodbattribute.MarshalMap(server.ToDo{
			ID:      "99211782-158f-4ccc-99fc-812a583c7e9d",
			Title:   "Test ToDo 1",
			ModTime: time.Now(),
		})
		if err != nil {
			t.Fatal(err)
		}

		item2, err := dynamodbattribute.MarshalMap(server.ToDo{
			ID:      "0ab230bc-dc6a-419f-8501-62eb629a34d2",
			Title:   "Test ToDo 2",
			ModTime: time.Now(),
		})
		if err != nil {
			t.Fatal(err)
		}

		item3, err := dynamodbattribute.MarshalMap(server.ToDo{
			ID:      "072f3d6e-45c3-4940-aee6-c17002aa7302",
			Title:   "Test ToDo 3",
			ModTime: time.Now(),
		})
		if err != nil {
			t.Fatal(err)
		}

		items := []map[string]*awsdynamodb.AttributeValue{item1, item2, item3}

		out := &awsdynamodb.ScanOutput{
			Items: items,
		}

		return out, nil
	}

	repo := dynamodb.NewToDoRepo(m)

	toDos, err := repo.GetAll()
	if err != nil {
		t.Fatal(err)
	}

	if len(toDos) != 3 {
		t.Fatal("Expected 3 ToDos in result")
	}

	if !m.ScanInvoked {
		t.Fatal("Scan not invoked")
	}
}

func testGetAllToDosError(t *testing.T) {

	m := &ClientMock{}

	m.ScanFn = func(*awsdynamodb.ScanInput) (*awsdynamodb.ScanOutput, error) {
		return nil, errors.New("DB Error")
	}

	repo := dynamodb.NewToDoRepo(m)

	_, err := repo.GetAll()
	if err == nil {
		t.Fatal("Expected Error")
	}

	if !m.ScanInvoked {
		t.Fatal("Scan not invoked")
	}

}

func testCreateToDo(t *testing.T) {

	m := &ClientMock{}

	m.PutItemFn = func(input *awsdynamodb.PutItemInput) (*awsdynamodb.PutItemOutput, error) {

		var toDo server.ToDo
		err := dynamodbattribute.UnmarshalMap(input.Item, &toDo)
		if err != nil {
			t.Fatal(err)
		}

		item, err := dynamodbattribute.MarshalMap(toDo)
		if err != nil {
			t.Fatal(err)
		}

		out := &awsdynamodb.PutItemOutput{
			Attributes: item,
		}

		return out, nil
	}

	repo := dynamodb.NewToDoRepo(m)

	newToDo := &server.ToDo{Title: "New ToDo"}

	err := repo.Save(newToDo)
	if err != nil {
		t.Fatal(err)
	}

	if newToDo.ID == "" {
		t.Fatal("Expected ToDo to have an ID")
	}

	if newToDo.ModTime.IsZero() {
		t.Fatal("Expected ToDo to have a not zero ModTime")
	}

	if !m.PutItemInvoked {
		t.Fatal("PutItem not invoked")
	}
}

func testCreateToDoError(t *testing.T) {

	m := &ClientMock{}

	m.PutItemFn = func(*awsdynamodb.PutItemInput) (*awsdynamodb.PutItemOutput, error) {
		return nil, errors.New("DB Error")
	}

	repo := dynamodb.NewToDoRepo(m)

	newToDo := &server.ToDo{Title: "New ToDo"}

	err := repo.Save(newToDo)
	if err == nil {
		t.Fatal("Expected Error")
	}

	if !m.PutItemInvoked {
		t.Fatal("PuItem not invoked")
	}

}

func testUpdateToDo(t *testing.T) {

	id := uuid.NewV4().String()

	m := &ClientMock{}

	m.PutItemFn = func(input *awsdynamodb.PutItemInput) (*awsdynamodb.PutItemOutput, error) {

		var toDo server.ToDo
		err := dynamodbattribute.UnmarshalMap(input.Item, &toDo)
		if err != nil {
			t.Fatal(err)
		}

		item, err := dynamodbattribute.MarshalMap(toDo)
		if err != nil {
			t.Fatal(err)
		}

		out := &awsdynamodb.PutItemOutput{
			Attributes: item,
		}

		return out, nil
	}

	repo := dynamodb.NewToDoRepo(m)

	toDoToUpdate := &server.ToDo{
		ID:        id,
		Title:     "Updated ToDo",
		Completed: true,
		ModTime:   time.Now(),
	}

	err := repo.Save(toDoToUpdate)
	if err != nil {
		t.Fatal(err)
	}

	if toDoToUpdate.ID != id {
		t.Fatalf("Expected ToDo to ID %s", id)
	}

	if toDoToUpdate.ModTime.IsZero() {
		t.Fatal("Expected ToDo to have a not zero ModTime")
	}

	if !m.PutItemInvoked {
		t.Fatal("PutItem not invoked")
	}
}

func testDeleteToDo(t *testing.T) {

	m := &ClientMock{}

	m.DeleteItemFn = func(input *awsdynamodb.DeleteItemInput) (*awsdynamodb.DeleteItemOutput, error) {

		toDo := &server.ToDo{
			ID:      testUUID,
			Title:   "Test ToDo",
			ModTime: time.Now(),
		}

		item, err := dynamodbattribute.MarshalMap(toDo)
		if err != nil {
			t.Fatal(err)
		}

		out := &awsdynamodb.DeleteItemOutput{
			Attributes: item,
		}

		return out, nil
	}

	repo := dynamodb.NewToDoRepo(m)

	err := repo.Delete(testUUID)
	if err != nil {
		t.Fatal(err)
	}

	if !m.DeleteItemInvoked {
		t.Fatal("DeleteItem not invoked")
	}
}

func testDeleteToDoError(t *testing.T) {

	m := &ClientMock{}

	m.DeleteItemFn = func(input *awsdynamodb.DeleteItemInput) (*awsdynamodb.DeleteItemOutput, error) {
		return nil, errors.New("DB Error")
	}

	repo := dynamodb.NewToDoRepo(m)

	err := repo.Delete(testUUID)
	if err == nil {
		t.Fatal("Expected Error")
	}

	if !m.DeleteItemInvoked {
		t.Fatal("DeleteItem not invoked")
	}
}

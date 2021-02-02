package dynamodb_test

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

// ClientMock is used to mock a client that uses makes call to DynamoDBAPI
type ClientMock struct {
	dynamodbiface.DynamoDBAPI
	GetItemFn         func(*dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error)
	ScanFn            func(*dynamodb.ScanInput) (*dynamodb.ScanOutput, error)
	PutItemFn         func(*dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error)
	DeleteItemFn      func(*dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error)
	GetItemInvoked    bool
	ScanInvoked       bool
	PutItemInvoked    bool
	DeleteItemInvoked bool
}

// GetItem returns a set of attributes for the item with the given primary key
func (m *ClientMock) GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	m.GetItemInvoked = true
	return m.GetItemFn(input)
}

// Scan returns one or more items and item attributes by accessing every item in a table or a secondary index
func (m *ClientMock) Scan(input *dynamodb.ScanInput) (*dynamodb.ScanOutput, error) {
	m.ScanInvoked = true
	return m.ScanFn(input)
}

// PutItem creates a new item, or replaces an old item with a new item
func (m *ClientMock) PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	m.PutItemInvoked = true
	return m.PutItemFn(input)
}

// DeleteItem deletes a single item in a table by primary key
func (m *ClientMock) DeleteItem(input *dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error) {
	m.DeleteItemInvoked = true
	return m.DeleteItemFn(input)

}

package main

import (
	awslambda "github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	awsdynamodb "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/massimoselvi/serverless-todo-api-go/server/database/dynamodb"
	"github.com/massimoselvi/serverless-todo-api-go/server/lambda/handlers"
)

func main() {

	s, err := session.NewSession(aws.NewConfig().WithRegion("us-west-2"))
	if err != nil {
		panic(err)
	}

	db := awsdynamodb.New(s)
	repo := dynamodb.NewToDoRepo(db)

	h := handlers.NewToDoHandler(repo)

	awslambda.Start(h.Handle)
}

package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/pkg/errors"
)

// CreateResponse generates an APIGatewayProxyResponse using the provided data and http code
func CreateResponse(data interface{}, code int) (events.APIGatewayProxyResponse, error) {

	r := events.APIGatewayProxyResponse{
		StatusCode: code,
	}

	// Try to marashal data, if it fail return errorResponse
	js, err := json.Marshal(data)
	if err != nil {
		r.StatusCode = http.StatusserverServerError
		js, err = json.Marshal(errorResponse{Err: err.Error()})
		if err != nil {
			return r, err
		}
	}

	r.Headers = make(map[string]string)
	r.Headers["Access-Control-Allow-Origin"] = "*"
	r.Headers["Access-Control-Allow-Credentials"] = "true"

	r.Body = string(js)

	return r, err
}

// CreateOKResponse generates an APIGatewayProxyResponse with a 200 http status code
func CreateOKResponse(data interface{}) (events.APIGatewayProxyResponse, error) {
	return CreateResponse(data, 200)
}

// CreateErrorResponse generates an APIGatewayProxyResponse using the provided error
func CreateErrorResponse(err error) (events.APIGatewayProxyResponse, error) {

	var code int

	switch errors.Cause(err) {
	case ErrNotFound:
		code = http.StatusNotFound
	case ErrBadRequest: //ToDO: what was bad?
		code = http.StatusBadRequest
	case ErrMethodNotAllowed:
		code = http.StatusMethodNotAllowed
	case ErrUnauthorized:
		code = http.StatusUnauthorized
	default:
		code = http.StatusserverServerError
	}

	e := &errorResponse{
		Err: err.Error(),
	}

	return CreateResponse(e, code)

}

var (
	// ErrNotFound is returned when an entity is not found
	ErrNotFound = errors.New("not found")
	// Errserver is returned when an server error has occurred
	Errserver = errors.New("server error")
	// ErrBadRequest is returned when the request is invalid
	ErrBadRequest = errors.New("bad request")
	// ErrMethodNotAllowed is returned when the request method (GET, POST, etc.) is not allowed
	ErrMethodNotAllowed = errors.New("method not allowed")
	// ErrUnauthorized is returned when the request is not authorized
	ErrUnauthorized = errors.New("unauthorized")
)

// errorResponse is the response sent to the client in the event of a error
type errorResponse struct {
	Err string `json:"error,omitempty"`
}

package server

import "time"

// ToDo represents details of a "todo" task to be compelted
type ToDo struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	ModTime   time.Time `json:"modTime"`
}

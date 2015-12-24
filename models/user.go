package models

import (
	"time"

	rethink "github.com/dancannon/gorethink"
)

// User define an object at the API
type User struct {
	ID       string    `gorethink:"id,omitempty"`
	Username string    `json:"username"`
	Password string    `json:"password"`
	Status   string    `json:"status"`
	Created  time.Time `json:"created"`
}

// Users define a collection of objects at the API
type Users struct {
	Conn *rethink.Session
}

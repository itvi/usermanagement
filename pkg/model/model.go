package model

import (
	"errors"
	"time"
)

var (
	// ErrNoRecord shows no record found
	ErrNoRecord = errors.New("models: no matching record found")
	// ErrInvalidCredentials will use this later if a user
	// tries to login with an incorrect email address or password.
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	// ErrDuplicate shows duplicate value in database table's field(see setup.sql)
	ErrDuplicate = errors.New("models: duplicate value")
)

// User ...
type User struct {
	ID             int
	SN             string
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

// Role ...
type Role struct {
	ID          int
	Name        string
	Description string
}

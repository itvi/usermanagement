package model

import (
	"time"
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

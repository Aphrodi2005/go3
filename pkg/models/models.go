package models

import (
	"errors"
	"time"
)

const (
	RoleTeacher = "teacher"
	RoleStudent = "student"
)

var (
	ErrNoRecord           = errors.New("models: no matching record found")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrDuplicateEmail     = errors.New("models: duplicate email")
)

type News struct {
	ID       int
	Title    string
	Content  string
	Details  string
	Created  time.Time
	Category string
}
type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
	Active         bool
	Role           string
}

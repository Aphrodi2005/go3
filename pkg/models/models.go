package models

import (
	"errors"
	"time"
)

const (
	RoleAdmin    = "admin"
	Rolesupplier = "supplier"
)

var (
	ErrNoRecord           = errors.New("models: no matching record found")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrDuplicateEmail     = errors.New("models: duplicate email")
)

var (
	ErrNoMovie   = errors.New("models: no matching movie found")
	ErrDuplicate = errors.New("models: duplicate movie title")
)

type Movie struct {
	ID          int
	Title       string
	Genre       string
	Rating      float64
	SessionTime time.Time
	CSRFToken   string
	IsAdmin     string
}
type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
	Role           string
}
type Ticket struct {
	ID          int
	UserID      int
	MovieID     int
	MovieTitle  string
	SessionTime time.Time
}

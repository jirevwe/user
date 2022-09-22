package models

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

// Model a base Go struct which includes the following fields: ID, CreatedAt, UpdatedAt, DeletedAt
//
//	type User struct {
//	  Model
//	}
type Model struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt null.Time `json:"deleted_at" db:"deleted_at"`
}

type User struct {
	FullName string `json:"full_name" db:"full_name"`
	Email    string `json:"email"`
	Password string `json:"-"`
	Model
}

type CreateUser struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

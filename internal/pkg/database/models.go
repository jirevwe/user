package database

import "gorm.io/gorm"

type User struct {
	gorm.Model
	FullName string `json:"full_name" gorm:"not null"`
	Email    string `json:"email" gorm:"unique;not null"`
	Password string `json:"-"`
}

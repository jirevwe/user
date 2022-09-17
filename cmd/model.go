package main

import "gorm.io/gorm"

type User struct {
	gorm.Model
	FullName string `json:"full_name" gorm:"not null"`
	Email    string `json:"email" gorm:"unique;not null"`
	Password string `json:"-"`
}

type UserResponse struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}

type SignUpRequest struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

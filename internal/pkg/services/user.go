package services

import (
	"errors"

	"github.com/jaevor/go-nanoid"
	"github.com/jirevwe/user/internal/pkg/models"
	"github.com/jmoiron/sqlx"
)

var ErrUserPasswordNotUpdated = errors.New("user password could not be update")

const (
	createUserTableQuery = `
	-- models/user.go:createUserTable
	CREATE TABLE IF NOT EXISTS users (
		full_name char(14) NOT NULL,
		email varchar(255) NOT NULL PRIMARY KEY UNIQUE,
		password varchar(255) NOT NULL
	);
	`

	insertUserRecordQuery = `
	-- models/user.go:Create
	INSERT INTO users (id, full_name, email, password) 
	VALUES ($1, $2, $3, $4);
	`

	findUserQuery = `
	--models/user.go:Authenticate
	SELECT * FROM users WHERE email = $1;
	`

	updateUserPassword = `
	--models/user.go:UpdateUserPassword
	UPDATE users SET password = $1 WHERE email = $2;
	`

	getAllUsers = `
	--models/user.go:getAllUsers
	SELECT * FROM users;
	`
)

// UserService contains all methods and fields for interacting
// with the `users` table in the database.
type UserService struct {
	DB *sqlx.DB
}

func NewUserService(db *sqlx.DB) (*UserService, error) {
	u := &UserService{DB: db}
	err := u.createUserTable()
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (u *UserService) createUserTable() error {
	_, err := u.DB.Exec(createUserTableQuery)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserService) Create(user *models.CreateUser) error {
	generator, err := nanoid.Standard(21)
	if err != nil {
		return err
	}
	id := generator()

	_, err = u.DB.Exec(insertUserRecordQuery, id, user.FullName, user.Email, user.Password)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserService) Authenticate(email string) (*models.User, error) {
	var user models.User
	err := u.DB.QueryRowx(findUserQuery, email).StructScan(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *UserService) UpdateUserPassword(email string, password string) error {
	result, err := u.DB.Exec(updateUserPassword, password, email)

	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()

	if rowsAffected < 1 {
		return ErrUserPasswordNotUpdated
	}

	if err != nil {
		return err
	}

	return nil
}

func (u *UserService) GetAllUsers() (*[]models.User, error) {
	var users []models.User
	rows, err := u.DB.Queryx(getAllUsers)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var user models.User
		_ = rows.StructScan(&user)
		users = append(users, user)
	}

	return &users, nil
}

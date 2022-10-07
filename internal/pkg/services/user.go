package services

import (
	"errors"

	"github.com/jaevor/go-nanoid"
	"github.com/jirevwe/user/internal/pkg/models"
	"github.com/jmoiron/sqlx"
)

var (
	ErrUserPasswordNotUpdated = errors.New("user password could not be update")
	ErrUserNotDeleted         = errors.New("user could not be deleted")
)

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

	findUserByIdQuery = `
	--models/user.go:Authenticate
	SELECT * FROM users WHERE id = $1;
	`

	updateUserPassword = `
	--models/user.go:UpdateUserPassword
	UPDATE users SET password = $1 WHERE id = $2;
	`

	getAllUsers = `
	--models/user.go:GetAllUsers
	SELECT * FROM users;
	`

	deleteUser = `
    --models/user.go:DeleteUser
    DELETE FROM users WHERE id = $1;
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

func (u *UserService) FindUserById(id string) (*models.User, error) {
	var user models.User
	err := u.DB.QueryRowx(findUserByIdQuery, id).StructScan(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *UserService) UpdateUserPassword(userId string, password string) error {
	result, err := u.DB.Exec(updateUserPassword, password, userId)

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

func (u *UserService) GetAllUsers() ([]models.User, error) {
	var users []models.User
	rows, err := u.DB.Queryx(getAllUsers)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var user models.User

		err = rows.StructScan(&user)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (u *UserService) DeleteUser(userId string) error {
	results, err := u.DB.Exec(deleteUser, userId)
	if err != nil {
		return err
	}

	rowsAffected, err := results.RowsAffected()

	if rowsAffected == 0 {
		return ErrUserNotDeleted
	}

	return nil
}
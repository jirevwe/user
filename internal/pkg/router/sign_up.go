package router

import (
	"encoding/json"
	"errors"
	"github.com/jirevwe/user/internal/pkg/database"
	"github.com/jirevwe/user/util"
	bcrypt2 "golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func SignUpRoute(db *gorm.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		var requestBody SignUpRequest
		err := util.DecodeJson(r.Body, &requestBody)

		if err != nil {
			_ = util.EncodeJson(w, err)
			return
		}

		if requestBody.Email == "" {
			_ = util.EncodeJson(w, err)
			return
		}

		if requestBody.FullName == "" {
			_ = util.EncodeJson(w, errors.New("full name cannot be empty"))
			return
		}

		if requestBody.Password == "" || len(requestBody.Password) < 8 {

			_ = util.EncodeJson(w, errors.New("password cannot be empty or less than 8 characters"))
			return
		}

		hashedPassword, bcryptErr := bcrypt2.GenerateFromPassword([]byte(requestBody.Password), 8)

		if bcryptErr != nil {
			log.Fatal(bcryptErr)
		}

		result := db.Create(&database.User{FullName: requestBody.FullName,
			Email:    requestBody.Email,
			Password: string(hashedPassword)})

		if result.Error != nil {
			_ = util.EncodeJson(w, result.Error)
			return
		}

		_ = json.NewEncoder(w).Encode(util.NewServerResponse(
			"User successfully created",
			UserResponse{
				Email:    requestBody.Email,
				FullName: requestBody.FullName,
			}, http.StatusOK))
	}

}

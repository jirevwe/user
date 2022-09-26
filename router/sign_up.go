package router

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/jirevwe/user/internal/pkg/database"
	"github.com/jirevwe/user/internal/pkg/models"
	"github.com/jirevwe/user/util"
	bcrypt2 "golang.org/x/crypto/bcrypt"
)

func SignUp(db database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		var requestBody SignUpRequest
		err := util.DecodeJson(r.Body, &requestBody)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = util.EncodeJson(w, err)
			return
		}

		if requestBody.Email == "" {
			w.WriteHeader(http.StatusBadRequest)
			_ = util.EncodeJson(w, err)
			return
		}

		if requestBody.FullName == "" {
			w.WriteHeader(http.StatusBadRequest)
			_ = util.EncodeJson(w, errors.New("full name cannot be empty"))
			return
		}

		if requestBody.Password == "" || len(requestBody.Password) < 8 {
			w.WriteHeader(http.StatusBadRequest)
			_ = util.EncodeJson(w, errors.New("password cannot be empty or less than 8 characters"))
			return
		}

		hashedPassword, bcryptErr := bcrypt2.GenerateFromPassword([]byte(requestBody.Password), 8)

		if bcryptErr != nil {
			log.Fatal(bcryptErr)
		}

		u := &models.CreateUser{
			FullName: requestBody.FullName,
			Email:    requestBody.Email,
			Password: string(hashedPassword),
		}

		err = db.GetUserService().Create(u)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = util.EncodeJson(w, err)
			return
		}

		_ = json.NewEncoder(w).Encode(util.NewServerResponse(
			"user account created",
			UserResponse{
				Email:    requestBody.Email,
				FullName: requestBody.FullName,
			}, http.StatusOK))
	}

}

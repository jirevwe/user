package router

import (
	"errors"
	"net/http"

	"github.com/jirevwe/user/internal/pkg/database"
	"github.com/jirevwe/user/util"
	bcrypt2 "golang.org/x/crypto/bcrypt"
)

var (
	ErrAccountNotFound    = errors.New("user with email/password does not exist")
	ErrEmailCannotBeEmpty = errors.New("email cannot be empty")
)

func Login(db database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		var requestBody LoginRequest
		err := util.DecodeJson(r.Body, &requestBody)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = util.EncodeJson(w, err)
			return
		}

		if requestBody.Email == "" {
			w.WriteHeader(http.StatusBadRequest)
			_ = util.EncodeJson(w, ErrEmailCannotBeEmpty)
			return
		}

		user, err := db.GetUserService().Authenticate(requestBody.Email)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			_ = util.EncodeJson(w, ErrAccountNotFound)
			return
		}

		err = bcrypt2.CompareHashAndPassword([]byte(user.Password), []byte(requestBody.Password))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = util.EncodeJson(w, ErrAccountNotFound)
			return
		}

		_ = util.EncodeJsonStatus(w, "login successful", http.StatusOK, user)
	}
}

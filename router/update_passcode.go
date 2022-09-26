package router

import (
	"errors"
	"github.com/jirevwe/user/internal/pkg/database"
	"github.com/jirevwe/user/util"
	bcrypt2 "golang.org/x/crypto/bcrypt"
	"net/http"
)

var ErrUserPasswordNotUpdated = errors.New("user password could not be update")

func UpdatePasscode(db database.Database) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		var requestBody UpdatePasscodeRequest

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

		if requestBody.NewPassword == "" {
			w.WriteHeader(http.StatusBadRequest)
			_ = util.EncodeJson(w, ErrEmailCannotBeEmpty)
			return
		}

		if requestBody.CurrentPassword == "" {
			w.WriteHeader(http.StatusBadRequest)
			_ = util.EncodeJson(w, ErrEmailCannotBeEmpty)
			return
		}

		user, err := db.GetUserService().Authenticate(requestBody.Email)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = util.EncodeJson(w, ErrAccountNotFound)
			return
		}

		err = bcrypt2.CompareHashAndPassword([]byte(user.Password), []byte(requestBody.CurrentPassword))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = util.EncodeJson(w, ErrAccountNotFound)
			return
		}

		hashedPassword, err := bcrypt2.GenerateFromPassword([]byte(requestBody.NewPassword), 8)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = util.EncodeJson(w, err)
		}

		err = db.GetUserService().UpdateUserPassword(requestBody.Email, string(hashedPassword))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = util.EncodeJson(w, err)
			return
		}

		_ = util.EncodeJsonStatus(w, "user updated successfully", http.StatusOK, user)
	}

}

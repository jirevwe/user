package router

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/jirevwe/user/internal/pkg/database"
	"github.com/jirevwe/user/util"
	bcrypt2 "golang.org/x/crypto/bcrypt"
)

func UpdatePasscode(db database.Database) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		var requestBody UpdatePasscodeRequest
		trimPath := strings.TrimPrefix(r.URL.Path, "/user/")
		userId := strings.TrimSuffix(trimPath, "/password")

		fmt.Printf("user id is %v", userId)

		err := util.DecodeJson(r.Body, &requestBody)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = util.EncodeJson(w, err)
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

		user, err := db.GetUserService().FindUserById(userId)
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

		err = db.GetUserService().UpdateUserPassword(userId, string(hashedPassword))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = util.EncodeJson(w, err)
			return
		}

		_ = util.EncodeJsonStatus(w, "user updated successfully", http.StatusOK, user)
	}

}

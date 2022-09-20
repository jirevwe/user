package router

import (
	"errors"
	"github.com/jirevwe/user/internal/pkg/database"
	"github.com/jirevwe/user/util"
	bcrypt2 "golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
)

func LoginRoute(db *gorm.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		var requestBody LoginRequest
		var user database.User
		err := util.DecodeJson(r.Body, &requestBody)
		if err != nil {
			_ = util.EncodeJson(w, err)
			return
		}

		if requestBody.Email == "" {
			_ = util.EncodeJson(w, errors.New("email cannot be empty"))
			return
		}

		result := db.First(&user, "email = ?", requestBody.Email)
		if result.Error == gorm.ErrRecordNotFound {
			_ = util.EncodeJson(w, errors.New("User with email/password does not exist"))
			return
		}

		if err = bcrypt2.CompareHashAndPassword([]byte(user.Password), []byte(requestBody.Password)); err != nil {
			_ = util.EncodeJson(w, errors.New("User with email/password does not exist"))
			return
		}

		_ = util.EncodeJsonStatus(w, "Success", http.StatusOK, user)

	}
}

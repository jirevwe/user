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

func UpdatePasscodeRoute(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")

		var requestBody UpdatePasscodeRequest
		var dbUser database.User

		err := util.DecodeJson(r.Body, &requestBody)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = util.EncodeJson(w, err)
			return
		}

		if requestBody.Email == "" {
			w.WriteHeader(http.StatusBadRequest)
			_ = util.EncodeJson(w, errors.New("email cannot be empty"))
			return
		}

		result := db.First(&dbUser, "email = ?", requestBody.Email)
		if result.Error == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			_ = util.EncodeJson(w, errors.New("User with email/password does not exist"))
			return
		}

		if err = bcrypt2.CompareHashAndPassword([]byte(dbUser.Password), []byte(requestBody.CurrentPassword)); err != nil {
			w.WriteHeader(http.StatusNotFound)
			_ = util.EncodeJson(w, errors.New("User with email/password does not exist"))
			return
		}

		hashedPassword, bcryptErr := bcrypt2.GenerateFromPassword([]byte(requestBody.NewPassword), 8)

		if bcryptErr != nil {
			log.Fatal(bcryptErr)
		}

		db.Model(dbUser).Update("password", hashedPassword)

		_ = json.NewEncoder(w).Encode(util.NewServerResponse(
			"successful",
			dbUser, http.StatusOK))
	}
}

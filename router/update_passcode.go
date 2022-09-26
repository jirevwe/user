package router

import (
	"github.com/jirevwe/user/internal/pkg/database"
	"github.com/jirevwe/user/util"
	log "github.com/sirupsen/logrus"
	bcrypt2 "golang.org/x/crypto/bcrypt"
	"net/http"
)

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
			_ = util.EncodeJson(w, ErrAccountNotFound)
			return
		}

		hashedPassword, bcryptErr := bcrypt2.GenerateFromPassword([]byte(requestBody.NewPassword), 8)

		if bcryptErr != nil {
			log.Fatal(bcryptErr)
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

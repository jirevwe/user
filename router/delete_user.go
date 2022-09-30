package router

import (
	"github.com/jirevwe/user/internal/pkg/database"
	"github.com/jirevwe/user/util"
	"net/http"
)

func DeleteUser(db database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		var requestBody DeleteRequest

		err := util.DecodeJson(r.Body, &requestBody)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = util.EncodeJson(w, err)
			return
		}

		err = db.GetUserService().DeleteUser(requestBody.Email)

		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			_ = util.EncodeJson(w, err)
			return
		}

		_ = util.EncodeJsonStatus(w, "user deleted", http.StatusOK, "user successfully deleted")
	}
}

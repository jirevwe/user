package router

import (
	"github.com/jirevwe/user/internal/pkg/database"
	"github.com/jirevwe/user/util"
	"net/http"
)

func FetchAllUsers(db database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")

		users, err := db.GetUserService().GetAllUsers()

		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			_ = util.EncodeJson(w, err)
			return
		}

		_ = util.EncodeJsonStatus(w, "all users", http.StatusOK, users)
	}
}

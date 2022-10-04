package router

import (
	"github.com/jirevwe/user/internal/pkg/database"
	"github.com/jirevwe/user/util"
	"net/http"
	"strings"
)

func DeleteUser(db database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		userId := strings.TrimPrefix(r.URL.Path, "/users/")

		err := db.GetUserService().DeleteUser(userId)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			_ = util.EncodeJson(w, err)
			return
		}

		_ = util.EncodeJsonStatus(w, "user deleted", http.StatusOK, "user successfully deleted")
	}
}

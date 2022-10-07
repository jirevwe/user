package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jirevwe/user/internal/pkg/database"
)

func NewRouter(db database.Database) http.Handler {
	r := chi.NewRouter()
	r.Route("/users", func(userRoute chi.Router) {
		userRoute.Post("/signup", SignUp(db))
		userRoute.Post("/login", Login(db))
		userRoute.Put("/{user-id}/password", UpdatePasscode(db))
		userRoute.Get("/", FetchAllUsers(db))
		userRoute.Delete("/{user-id}", DeleteUser(db))
	})

	return r
}

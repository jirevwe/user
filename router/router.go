package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jirevwe/user/internal/pkg/database"
)

func NewRouter(db database.Database) http.Handler {
	r := chi.NewRouter()
	r.Route("/user", func(userRoute chi.Router) {
		userRoute.Post("/signup", SignUp(db))
		userRoute.Post("/login", Login(db))
		userRoute.Put("/passcode", UpdatePasscode(db))
		userRoute.Get("/all", FetchAllUsers(db))
		userRoute.Delete("/", DeleteUser(db))
	})

	return r
}

package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jirevwe/user/internal/pkg/database"
)

func NewRouter(db database.Database) http.Handler {
	r := chi.NewRouter()
	r.Route("/user", func(userRoute chi.Router) {
		// userRoute.Post("/signup", SignUpRoute(db))
		// userRoute.Post("/login", LoginRoute(db))
	})

	return r
}

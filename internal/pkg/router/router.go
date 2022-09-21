package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func NewRouter(db *gorm.DB) http.Handler {
	r := chi.NewRouter()
	r.Route("/user", func(userRoute chi.Router) {
		userRoute.Post("/signup", SignUpRoute(db))
		userRoute.Post("/login", LoginRoute(db))
		userRoute.Put("/passcode/update", UpdatePasscodeRoute(db))
	})

	return r
}

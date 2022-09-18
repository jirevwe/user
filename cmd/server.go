package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"github.com/jirevwe/user/util"
	bcrypt2 "golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/go-chi/chi/v5"
	"github.com/jirevwe/user/internal/pkg/server"
	"github.com/spf13/cobra"

	log "github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:     "server",
	Aliases: []string{"s"},
	Short:   "Starts the http server",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})

		if err != nil {
			panic("failed to connect database")
		}

		log.SetLevel(log.InfoLevel)

		log.SetFormatter(&prefixed.TextFormatter{
			DisableColors:   false,
			TimestampFormat: "2006-01-02 15:04:05",
			FullTimestamp:   true,
			ForceFormatting: true,
		})
		log.SetReportCaller(true)

		err = os.Setenv("TZ", "") // Use UTC by default :)
		if err != nil {
			log.Fatal("failed to set env - ", err)
		}

		r := chi.NewRouter()
		srv := server.NewServer(9000)
		srv.SetHandler(r.Route("/user", func(userRoute chi.Router) {
			userRoute.Post(
				"/signup",
				func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("content-type", "application/json")
					var requestBody SignUpRequest
					err := util.DecodeJson(r.Body, &requestBody)

					if err != nil {
						_ = util.EncodeJson(w, err)
						return
					}

					if requestBody.Email == "" {
						_ = util.EncodeJson(w, err)
						return
					}

					if requestBody.FullName == "" {
						_ = util.EncodeJson(w, errors.New("full name cannot be empty"))
						return
					}

					if requestBody.Password == "" || len(requestBody.Password) < 8 {

						_ = util.EncodeJson(w, errors.New("password cannot be empty or less than 8 characters"))
						return
					}

					hashedPassword, bcryptErr := bcrypt2.GenerateFromPassword([]byte(requestBody.Password), 8)

					if bcryptErr != nil {
						log.Fatal(bcryptErr)
					}

					result := db.Create(&User{FullName: requestBody.FullName,
						Email:    requestBody.Email,
						Password: string(hashedPassword)})

					if result.Error != nil {
						_ = util.EncodeJson(w, result.Error)
						return
					}

					_ = json.NewEncoder(w).Encode(util.NewServerResponse(
						"User successfully created",
						UserResponse{
							Email:    requestBody.Email,
							FullName: requestBody.FullName,
						}, http.StatusOK))

				},
			)

			userRoute.Post("/login", func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("content-type", "application/json")
				var requestBody LoginRequest
				var user User
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

			})
		}))
		log.Infof("server running on port %v", 9000)
		srv.Listen()

		// Migrate the schema
		dbMigration(db)
	},
}

func dbMigration(db *gorm.DB) {
	err := db.AutoMigrate(&User{})
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	rootCmd.AddCommand(serverCmd)
}

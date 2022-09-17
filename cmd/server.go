package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jirevwe/user/util"
	bcrypt2 "golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/jirevwe/user/internel/pkg/server"
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
					err := json.NewDecoder(r.Body).Decode(&requestBody)

					if err != nil {
						_ = json.NewEncoder(w).Encode(util.NewServiceErrResponse(err))
						return
					}

					if requestBody.Email == "" {
						_ = json.NewEncoder(w).Encode(util.NewServiceErrResponse(errors.New("email cannot be empty")))
						return
					}

					if requestBody.FullName == "" {
						_ = json.NewEncoder(w).Encode(util.NewServiceErrResponse(errors.New("full name cannot be empty")))
						return
					}

					if requestBody.Password == "" || len(requestBody.Password) < 8 {
						_ = json.NewEncoder(w).Encode(util.NewServiceErrResponse(
							errors.New("password cannot be empty or less than 8 characters")))
						return
					}

					hashedPassword, bcryptErr := bcrypt2.GenerateFromPassword([]byte(requestBody.Password), 8)

					if bcryptErr != nil {
						log.Fatal(bcryptErr)
						return
					}

					result := db.Create(&User{FullName: requestBody.FullName,
						Email:    requestBody.Email,
						Password: string(hashedPassword)})

					if result.Error != nil {
						_ = json.NewEncoder(w).Encode(util.NewServiceErrResponse(
							result.Error))
						return
					}

					_ = json.NewEncoder(w).Encode(util.NewServerResponse(
						"User succesffully created",
						UserResponse{
							Email:    requestBody.Email,
							FullName: requestBody.FullName,
						}, http.StatusOK))

				},
			)

			userRoute.Get("/user", func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("content-type", "application/json")

				if err != nil {
					log.Fatal(err)
				}
				_ = json.NewEncoder(w).Encode(util.NewServerResponse(
					"Successful", SignUpRequest{Email: "wkjf"}, http.StatusOK))
			})

			userRoute.Post("/login", Login)
		}))
		log.Infof("server running on port %v", 9000)
		srv.Listen()

		// Migrate the schema
		db.AutoMigrate(&User{})
	},
}

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	var requestBody LoginRequest
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		_ = json.NewEncoder(w).Encode(util.NewServiceErrResponse(
			err))
		return
	}

	fmt.Println(requestBody)
}

func init() {
	rootCmd.AddCommand(serverCmd)
}

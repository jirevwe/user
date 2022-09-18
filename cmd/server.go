package main

import (
	"os"
	"time"

	"github.com/jirevwe/user/internal/pkg/database"
	"github.com/jirevwe/user/internal/pkg/router"
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
		log.SetLevel(log.InfoLevel)
		log.SetFormatter(&prefixed.TextFormatter{
			DisableColors:   false,
			TimestampFormat: time.RFC3339,
			FullTimestamp:   true,
			ForceFormatting: true,
		})
		log.SetReportCaller(true)

		err := os.Setenv("TZ", "") // Use UTC by default :)
		if err != nil {
			log.Fatal("failed to set env - ", err)
		}

		db := database.NewDB()
		r := router.NewRouter(db)

		srv := server.NewServer(9000)
		srv.SetHandler(r)
		log.Infof("server running on port %v", 9000)
		srv.Listen()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}

package main

import (
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
		log.SetLevel(log.InfoLevel)

		log.SetFormatter(&prefixed.TextFormatter{
			DisableColors:   false,
			TimestampFormat: "2006-01-02 15:04:05",
			FullTimestamp:   true,
			ForceFormatting: true,
		})
		log.SetReportCaller(true)

		err := os.Setenv("TZ", "") // Use UTC by default :)
		if err != nil {
			log.Fatal("failed to set env - ", err)
		}

		srv := server.NewServer(9000)
		srv.SetHandler(chi.NewRouter())
		log.Infof("server running on port %v", 9000)
		srv.Listen()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}

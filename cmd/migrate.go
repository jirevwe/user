package main

import (
	"os"
	"time"

	"github.com/jirevwe/user/internal/pkg/database"
	"github.com/jirevwe/user/internal/pkg/migrator"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var direction string

// migrateCmd represents the version command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Performs SQL migrations",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			direction = args[0]
		}

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

		db, err := database.New()
		if err != nil {
			log.Fatal(err)
		}

		m := migrator.New(db)
		switch direction {
		case "up":
			err := m.Up()
			if err != nil {
				logrus.Fatalf("migration up failed with error: %+v", err)
			}
		case "down":
			err := m.Down()
			if err != nil {
				logrus.Fatalf("migration down failed with error: %+v", err)
			}
		default:
			log.Fatal("migration failed, invalid direction specified")
		}
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}

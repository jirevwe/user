package main

import (
	"fmt"

	"github.com/jirevwe/user/version"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print out the cli version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s\n", version.GetVersion())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

//sign
//login
//forgot password
//update password

package cmd

import (
	"portal-blog/internal/app"

	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start",
	Long:  "start",
	Run: func(cmd *cobra.Command, args []string) {
		app.RunServer()
	},
}

// init initializes the start command by adding it to the root command.
// This function is automatically called when the package is imported.
func init() {
	rootCmd.AddCommand(startCmd)
}

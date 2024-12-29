package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "core-api",
	Short: "this api for news portal",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Run(startCmd, nil)
	},
}

// Execute runs the root command of the application.
// It uses cobra's CheckErr function to handle any errors that occur during execution.
// If an error occurs, the program will terminate and print the error message.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

// init initializes the root command for the application.
// It sets up the configuration initialization and defines command-line flags.
//
// This function does the following:
// 1. Registers the initConfig function to be called when cobra initializes.
// 2. Adds a persistent flag for specifying a custom config file.
// 3. Adds a boolean toggle flag for demonstration purposes.
//
// The function takes no parameters and returns nothing.
func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is .env)")

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

// initConfig initializes the configuration for the application.
// It sets up the configuration file and environment variables.
//
// If a config file is specified via the global flag, it will be used.
// Otherwise, it defaults to using a .env file in the current directory.
//
// The function also sets up automatic environment variable reading
// and prints the config file being used if there's an error reading it.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigFile(`.env`)
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Fprintln(os.Stderr, "Using config file: ", viper.ConfigFileUsed())
	}
}

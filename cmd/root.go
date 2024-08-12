package cmd

import (
	"fmt"
	"os"

	"github.com/d4r1us-drk/clido/internal/version"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "clido",
	Short: "Clido is an awesome CLI to-do list management application",
	Long: `Clido is a simple yet powerful CLI tool designed to help you manage 
your projects and tasks effectively from the terminal.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Clido",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version.FullVersion())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

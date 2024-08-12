package cmd

import (
	"fmt"
	"os"

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

func init() {
}

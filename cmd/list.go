package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/d4r1us-drk/clido/pkg/repository"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list [projects|tasks]",
	Short: "List projects or tasks",
	Long:  `List all projects or tasks, optionally filtered by project for tasks.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Insufficient arguments. Use 'list projects' or 'list tasks'.")
			return
		}

		repo, err := repository.NewRepository()
		if err != nil {
			fmt.Printf("Error initializing repository: %v\n", err)
			return
		}
		defer repo.Close()

		switch args[0] {
		default:
			fmt.Println("Invalid option. Use 'list projects' or 'list tasks'.")
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringP("project", "p", "", "Filter tasks by project name or ID")
}

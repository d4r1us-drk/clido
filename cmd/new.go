package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/d4r1us-drk/clido/pkg/repository"
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new [project|task]",
	Short: "Create a new project or task",
	Long:  `Create a new project or task with the specified details.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Insufficient arguments. Use 'new project' or 'new task'.")
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
			fmt.Println("Invalid option. Use 'new project' or 'new task'.")
		}
	},
}


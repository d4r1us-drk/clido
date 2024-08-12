package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/d4r1us-drk/clido/pkg/repository"
	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit [project|task] <id>",
	Short: "Edit an existing project or task",
	Long:  `Edit the details of an existing project or task identified by its ID.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			fmt.Println("Insufficient arguments. Use 'edit project <id>' or 'edit task <id>'.")
			return
		}

		repo, err := repository.NewRepository()
		if err != nil {
			fmt.Printf("Error initializing repository: %v\n", err)
			return
		}
		defer repo.Close()

		id, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("Invalid ID. Please provide a numeric ID.")
			return
		}

		switch args[0] {
		default:
			fmt.Println("Invalid option. Use 'edit project <id>' or 'edit task <id>'.")
		}
	},
}

func init() {
	rootCmd.AddCommand(editCmd)

	editCmd.Flags().StringP("name", "n", "", "New name")
	editCmd.Flags().StringP("description", "d", "", "New description")
	editCmd.Flags().StringP("due", "D", "", "New due date for task (format: YYYY-MM-DD HH:MM)")
	editCmd.Flags().
		IntP("priority", "r", 0, "New priority for task (1: High, 2: Medium, 3: Low, 4: None)")
}


package main

import (
	"log"

	"github.com/d4r1us-drk/clido/cmd"
	"github.com/d4r1us-drk/clido/controllers"
	"github.com/d4r1us-drk/clido/repository"
)

func main() {
	// Initialize the repository
	repo, repoErr := repository.NewRepository()
	if repoErr != nil {
		log.Printf("Error initializing repository: %v", repoErr)
		return
	}
	defer repo.Close()

	// Initialize controllers
	projectController := controllers.NewProjectController(repo)
	taskController := controllers.NewTaskController(repo)

	// Initialize the root command with controllers
	rootCmd := cmd.NewRootCmd(projectController, taskController)

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		log.Printf("Error executing command: %v", err)
		repo.Close()
	}
}

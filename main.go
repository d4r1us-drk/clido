package main

import (
	"os"

	"github.com/fatih/color"
)

func main() {
	if len(os.Args) < 2 {
		handleHelp()
		return
	}

	repo, err := NewRepository()
	if err != nil {
		color.Red("Error initializing repository: %v\n", err)
		return
	}
	defer repo.Close()

	switch os.Args[1] {
	case "new":
		handleNew(repo, os.Args[2:])
	case "edit":
		handleEdit(repo, os.Args[2:])
	case "list":
		handleList(repo, os.Args[2:])
	case "remove":
		handleRemove(repo, os.Args[2:])
	case "toggle":
		handleToggle(repo, os.Args[1:])
	case "help":
		handleHelp()
	default:
		color.Red("Invalid command.")
		handleHelp()
	}
}

package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

func main() {
	if len(os.Args) < 2 {
		handleHelp()
		return
	}

	homePath := os.Getenv("HOME")
	if homePath == "" {
		color.Red("The HOME environment variable is not set.")
		return
	}

	dbPath := fmt.Sprintf("%s/.local/share/clido/data.db", homePath)
	repo, err := NewRepository(dbPath)
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
	case "set":
		handleSet(repo, os.Args[2:])
	case "help":
		handleHelp()
	default:
		color.Red("Invalid command.")
		handleHelp()
	}
}

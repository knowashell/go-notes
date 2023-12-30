package main

import (
	"fmt"
	"os"

	"go-notes/internal/cli"
	"go-notes/internal/storage/sqlite"
)

const storageName = "storage.db" // Name of the SQLite database file

func main() {
	// initialize the sqlite storage using the specified database file name
	storage, err := sqlite.New(storageName)
	if err != nil {
		fmt.Printf("Error initializing storage: %v\n", err)
		os.Exit(1)
	}

	defer func() {
		// close the storage when the main function exits
		if err = storage.Close(); err != nil {
			fmt.Printf("Error closing storage: %v\n", err)
		}
	}()

	// create a new CLI application with the initialized storage
	app := cli.NewCLI(storage)

	// run the CLI application with the command-line arguments passed to the program
	err = app.Run(os.Args)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

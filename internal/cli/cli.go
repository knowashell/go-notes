package cli

import (
	"fmt"
	"strconv"

	"github.com/urfave/cli"

	"go-notes/internal/entities"
)

type Storage interface {
	// NewNote creates a new note with the given title and content and returns its ID
	NewNote(noteTitle, content string) (int, error)

	// DeleteNote deletes a note by its ID
	DeleteNote(id int) (int, error)

	// SetNoteContent updates the content of a note with the specified ID
	SetNoteContent(noteID int, content string) error

	// GetNoteByID retrieves a note by its ID and returns it as an entities.Note
	GetNoteByID(noteID int) (entities.Note, error)

	// GetAllNotes retrieves all notes and returns them as a slice of entities.Note
	GetAllNotes() ([]entities.Note, error)

	// SearchNotesByKeyword searches for notes containing the specified keyword and returns them as a slice of entities.Note
	SearchNotesByKeyword(keyword string) ([]entities.Note, error)
}

const (
	appName  = "Note Storage CLI"            // name of CLI application
	appUsage = "Manage your notes using CLI" // description of application's purpose
)

// NewCLI creates new CLI application with provided storage object
func NewCLI(storage Storage) *cli.App {
	// create a new CLI application
	app := cli.NewApp()
	app.Name = appName   // set application's name
	app.Usage = appUsage // set application's usage description

	// define available commands for CLI application
	app.Commands = []cli.Command{
		newNoteCommand(storage),           // create a new note
		deleteNoteCommand(storage),        // delete a note by ID
		getNoteByIDCommand(storage),       // get a note by ID
		listNotesCommand(storage),         // list all notes
		updateNoteContentCommand(storage), // update content of a note
		searchNotesCommand(storage),       // search notes by keyword in title or content
	}

	return app
}

// updateNoteContentCommand creates new CLI command with provided storage object
func updateNoteContentCommand(storage Storage) cli.Command {
	// constants for command name and usage description
	const (
		commandName  = "update"
		commandUsage = "Update content of a note"
	)

	// create a new CLI command configuration
	updateNoteContent := cli.Command{
		Name:  commandName,  // name of command (e.g., "update")
		Usage: commandUsage, // description of command
		Action: func(c *cli.Context) error {
			// retrieve first argument as note ID
			noteIDStr := c.Args().First()
			if noteIDStr == "" {
				fmt.Println("Please provide ID of note to update.")
				return nil
			}

			// convert note ID string to an integer
			noteID, err := strconv.Atoi(noteIDStr)
			if err != nil {
				return fmt.Errorf("invalid note ID: %w", err)
			}

			// retrieve second argument as new content for note
			content := c.Args().Get(1)
			if content == "" {
				fmt.Println("Please provide content to update note.")
				return nil
			}

			// call a function from 'storage' object to update note's content
			err = storage.SetNoteContent(noteID, content)
			if err != nil {
				return fmt.Errorf("updating note: %w", err)
			}

			fmt.Printf("Updated note with ID %d\n", noteID)

			return nil
		},
	}

	return updateNoteContent
}

// searchNotesCommand creates a new CLI command for searching notes by keyword.
func searchNotesCommand(storage Storage) cli.Command {
	// constants for command name and usage description.
	const (
		commandName  = "search"
		commandUsage = "Search notes by keyword"
	)

	// create a new CLI command configuration
	searchNotes := cli.Command{
		Name:  commandName,  // name of command (e.g., "update")
		Usage: commandUsage, // description of command
		Action: func(c *cli.Context) error {
			// extract the command-line argument as the keyword to search for
			keyword := c.Args().First()
			if keyword == "" {
				fmt.Println("Please provide a keyword to search for notes.")
				return nil
			}

			// call method from the 'storage' object to search for notes
			notes, err := storage.SearchNotesByKeyword(keyword)
			if err != nil {
				fmt.Printf("Error searching notes: %v\n", err)
				return err
			}

			// display search results
			if len(notes) == 0 {
				fmt.Printf("No notes found for keyword: %s\n", keyword)
			} else {
				fmt.Printf("Notes found for keyword '%s':\n", keyword)
				for _, note := range notes {
					fmt.Printf("ID: %d, Title: %s, Content: %s, CreatedAt: %s, LastEditedAt: %s\n",
						note.ID, note.Title, note.Content, note.CreatedAt, note.LastEditedAt)
				}
			}

			return nil
		},
	}

	return searchNotes
}

// getNoteByIDCommand creates new CLI command with provided storage object
func getNoteByIDCommand(storage Storage) cli.Command {
	// constants for command name and usage description
	const (
		commandName  = "get"
		commandUsage = "Get a note by ID"
	)

	// create a new CLI command configuration
	getNoteByID := cli.Command{
		Name:  commandName,  // name of command (e.g., "get")
		Usage: commandUsage, // description of command
		Action: func(c *cli.Context) error {
			// retrieve first argument as note ID
			noteIDStr := c.Args().First()
			if noteIDStr == "" {
				fmt.Println("Please provide ID of note to retrieve.")
				return nil
			}

			// convert note ID string to an integer
			noteID, err := strconv.Atoi(noteIDStr)
			if err != nil {
				return fmt.Errorf("invalid note ID: %w", err)
			}

			// call a function from 'storage' object to retrieve note by its ID
			note, err := storage.GetNoteByID(noteID)
			if err != nil {
				return fmt.Errorf("retrieving note: %w", err)
			}

			// print details of retrieved note
			fmt.Printf("Note ID: %d\nTitle: %s\nContent: %s\nCreatedAt: %s\nLastEditedAt: %s\n",
				note.ID, note.Title, note.Content, note.CreatedAt, note.LastEditedAt)

			return nil
		},
	}

	return getNoteByID
}

// listNotesCommand creates new CLI command with provided storage object
func listNotesCommand(storage Storage) cli.Command {
	// constants for command name and usage description
	const (
		commandName  = "list"
		commandUsage = "List all notes"
	)

	// create a new CLI command configuration
	listNotes := cli.Command{
		Name:  commandName,  // name of command (e.g., "list")
		Usage: commandUsage, // description of command
		Action: func(c *cli.Context) error {
			// call a function from 'storage' object to retrieve all notes
			notes, err := storage.GetAllNotes()
			if err != nil {
				fmt.Printf("Error listing notes: %v\n", err)
				return err
			}

			// print a header for list of notes
			fmt.Println("List of notes:")

			// iterate through retrieved notes and print their details
			for _, note := range notes {
				fmt.Printf("ID: %d, Title: %s, CreatedAt: %s, LastEditedAt: %s\n",
					note.ID, note.Title, note.CreatedAt, note.LastEditedAt)
			}

			return nil
		},
	}

	return listNotes
}

// deleteNoteCommand creates new CLI command for deleting note from storage with provided storage object
func deleteNoteCommand(storage Storage) cli.Command {
	// constants for command name and usage description
	const (
		commandName  = "delete"
		commandUsage = "Delete a note by ID"
	)

	// create a new CLI command configuration.
	deleteNote := cli.Command{
		Name:  commandName,  // name of command (e.g., "delete")
		Usage: commandUsage, // description of command
		Action: func(c *cli.Context) error {
			// retrieve first argument as note ID
			noteIDStr := c.Args().First()
			if noteIDStr == "" {
				return fmt.Errorf("please provide ID of note to delete")
			}

			// convert note ID string to an integer
			noteID, err := strconv.Atoi(noteIDStr)
			if err != nil {
				return fmt.Errorf("invalid note ID: %w", err)
			}

			// call a function from storage object to delete note by its ID
			deletedNoteID, err := storage.DeleteNote(noteID)
			if err != nil {
				return fmt.Errorf("Error deleting note: %v\n", err)
			}

			fmt.Printf("Deleted note with ID %d\n", deletedNoteID)

			return nil
		},
	}

	return deleteNote
}

// newNoteCommand creates new CLI command for creating new notes in storage with provided storage object
func newNoteCommand(storage Storage) cli.Command {
	// constants for command name and usage description
	const (
		commandName  = "new"
		commandUsage = "Create a new note"
	)

	// create a new CLI command configuration
	newNote := cli.Command{
		Name:  commandName,  // name of command (e.g., "new")
		Usage: commandUsage, // description of command
		Action: func(c *cli.Context) error {
			// retrieve first argument as title of new note
			title := c.Args().First()
			if title == "" {
				fmt.Println("Please provide a title for new note.")
				return nil
			}

			// retrieve second argument as content of new note
			content := c.Args().Get(1)
			if content == "" {
				fmt.Println("Please provide content for new note.")
				return nil
			}

			// call a function from 'storage' object to create a new note with provided title
			noteID, err := storage.NewNote(title, content)
			if err != nil {
				return fmt.Errorf("creating new note: %v\n", err)
			}

			fmt.Printf("Created a new note with ID %d\n", noteID)

			return nil
		},
	}

	return newNote
}

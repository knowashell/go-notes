package sqlite

import (
	"database/sql"
	"errors"
	"math"

	_ "github.com/mattn/go-sqlite3"

	"go-notes/internal/entities"
)

type (
	Storage struct {
		// db holds the database connection.
		db *sql.DB
	}
)

var (
	invalidNum         = errors.New("invalid number")
	invalidParamLength = errors.New("invalid param length")
)

// New creates a new Storage instance and establishes a connection to the SQLite database
func New(storagePath string) (*Storage, error) {
	// opening connection to sqlite db
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		// return error if connection fails
		return nil, err
	}

	// prepare statement to create a table
	createTable, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS notes (
    		note_id INTEGER PRIMARY KEY,
    		title TEXT NOT NULL,
    		content TEXT,
    		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    		last_edited_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP);
	`)
	if err != nil {
		// return error if preparing fails
		return nil, err
	}
	// ensure statement are closed when done processing
	defer createTable.Close()

	// createTable execution
	_, err = createTable.Exec()
	if err != nil {
		// return error if exec statement is fails
		return nil, err
	}

	// preparing statement to create a trigger for updating last edit of note
	onUpdateTrigger, err := db.Prepare(`
		CREATE TRIGGER IF NOT EXISTS update_last_edited_at
		AFTER UPDATE ON notes
		FOR EACH ROW
		BEGIN
    		UPDATE notes
    		SET last_edited_at = CURRENT_TIMESTAMP
    		WHERE note_id = OLD.note_id;
		END;
`)
	if err != nil {
		// return error if preparing fails
		return nil, err
	}
	// ensure statement are closed when done processing
	defer onUpdateTrigger.Close()

	// creating trigger execution
	_, err = onUpdateTrigger.Exec()
	if err != nil {
		// return err if creating trigger exec fails
		return nil, err
	}

	// returning new storage with established db connect
	return &Storage{db: db}, nil
}

// Close closes the database connection associated with the Storage instance
func (s *Storage) Close() error {
	err := s.db.Close()

	return err
}

// NewNote creates a new note with the given title and content and returns its ID
func (s *Storage) NewNote(noteTitle, content string) (int, error) {
	err := validateSQLParam(noteTitle, content)
	if err != nil {
		return 0, err
	}
	// preparing statement for creating new note with title and content
	newNote, err := s.db.Prepare("INSERT INTO notes (title, content) VALUES (?, ?)")
	if err != nil {
		// return error if preparing fails
		return 0, err
	}
	// ensure statement are closed when done processing
	defer newNote.Close()

	// creating new note execution with title and content
	res, err := newNote.Exec(noteTitle, content)
	if err != nil {
		// return err if execution fails
		return 0, err
	}

	// getting id of new note
	id, err := res.LastInsertId()

	// return id and error
	return int(id), err
}

// DeleteNote deletes a note by its ID
func (s *Storage) DeleteNote(id int) (int, error) {
	err := validateSQLParam(id)
	if err != nil {
		return 0, err
	}
	// preparing statement for deleting note by id
	deleteNote, err := s.db.Prepare("DELETE FROM notes WHERE note_id = ?")
	if err != nil {
		return 0, err
	}
	// ensure statement are closed when done processing
	defer deleteNote.Close()

	// execute delete statement
	result, err := deleteNote.Exec(id)
	if err != nil {
		return 0, err
	}

	// check number of rows affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	// if no rows were affected - return an error
	if rowsAffected == 0 {
		return 0, sql.ErrNoRows
	}

	// return ID of deleted note
	return id, nil
}

// SetNoteContent updates the content of a note with the specified ID
func (s *Storage) SetNoteContent(noteID int, content string) error {
	err := validateSQLParam(noteID, content)
	if err != nil {
		return err
	}
	// preparing statement for setting note content by id
	setNoteContent, err := s.db.Prepare("UPDATE notes SET content = ? WHERE note_id = ?")
	if err != nil {
		return err
	}
	// ensure statement are closed when done processing
	defer setNoteContent.Close()

	// execute setting note content
	res, err := setNoteContent.Exec(content, noteID)
	if err != nil {
		return err
	}

	// check number of rows affected
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	// if no rows were affected - return an error
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// SearchNotesByKeyword searches for notes containing the specified keyword in titles or content
func (s *Storage) SearchNotesByKeyword(keyword string) ([]entities.Note, error) {
	err := validateSQLParam(keyword)
	if err != nil {
		return nil, err
	}
	// SQL query to search for notes containing the keyword in titles or content
	query := "SELECT * FROM notes WHERE title LIKE ? OR content LIKE ?"

	// create a wildcard pattern for keyword (e.g., "%keyword%") to match partial strings
	keywordPattern := "%" + keyword + "%"

	// execute the query with the keyword pattern twice (for title and content) and retrieve the result rows
	rows, err := s.db.Query(query, keywordPattern, keywordPattern)
	if err != nil {
		return nil, err
	}

	// ensure rows are closed when done processing
	defer rows.Close()

	// create a slice to store matching notes
	var notes []entities.Note

	// iterate through result rows
	for rows.Next() {
		// declare a variable to store a single note
		var note entities.Note

		// scan the values from the row into 'note' struct
		err = rows.Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt, &note.LastEditedAt)
		if err != nil {
			return []entities.Note{}, err
		}

		// append the retrieved note to 'notes' slice
		notes = append(notes, note)
	}

	// return the list of matching notes and any error that occurred
	return notes, nil
}

// GetNoteByID retrieves a note by its ID and returns it as an entities.Note
func (s *Storage) GetNoteByID(noteID int) (entities.Note, error) {
	err := validateSQLParam(noteID)
	if err != nil {
		return entities.Note{}, err
	}
	// SQL query to select a note by its ID
	getNoteQuery := `SELECT * FROM notes WHERE note_id = ?`

	// declare a variable to store the retrieved note
	var note entities.Note

	// execute the query and scan the result into the 'note' struct
	err = s.db.QueryRow(getNoteQuery, noteID).Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt, &note.LastEditedAt)

	// return the retrieved note and any error that occurred
	return note, err
}

// GetAllNotes retrieves all notes and returns them as a slice of entities.Note
func (s *Storage) GetAllNotes() ([]entities.Note, error) {
	// execute an SQL query to retrieve all notes from table
	rows, err := s.db.Query(`SELECT * FROM notes`)
	if err != nil {
		return nil, err
	}
	// ensure rows are closed when done processing
	defer rows.Close()

	// create a slice to store the retrieved notes
	var notes []entities.Note

	// iterate through the result rows
	for rows.Next() {
		// create a variable to store a single note
		var note entities.Note

		// scan the values from the row into the 'note' struct
		err = rows.Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt, &note.LastEditedAt)
		if err != nil {
			return nil, err
		}

		// append the retrieved note to the 'notes' slice
		notes = append(notes, note)
	}

	// return slice of notes
	return notes, nil
}

// validateSQLParam validates parameters based on their type and value
// it checks if integers are within a valid range and if strings have a valid length
func validateSQLParam(params ...interface{}) error {
	// define a constant for the maximum allowed string length
	const maxStringLength = 256000

	// iterate over each parameter in variadic 'params' slice
	for _, param := range params {
		// use a type switch to check type of the parameter
		switch v := param.(type) {
		case int:
			// if the parameter is an integer:
			// check if it's within valid range
			if v < 1 || v > math.MaxInt32 {
				// If not - return an error with a message
				return invalidNum
			}
		case string:
			// if the parameter is a string
			// check if its length is within the valid range
			if len(v) < 1 || len(v) > maxStringLength {
				// if not - return an error with a message
				return invalidParamLength
			}
		}
	}

	// if all parameters pass validation, return nil (no error)
	return nil
}

package sqlite

import (
	"os"
	"testing"
)

func TestNewStorage(t *testing.T) {
	dbPath := "test.db"
	defer func() {
		_ = os.Remove(dbPath)
	}()

	storage, err := New(dbPath)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if storage.db == nil {
		t.Error("Expected database connection, got nil")
	}

	// Закрыть соединение с базой данных
	err = storage.Close()
	if err != nil {
		t.Errorf("Expected no error on Close(), got %v", err)
	}
}

func TestNewNote(t *testing.T) {
	dbPath := "test.db"
	defer func() {
		_ = os.Remove(dbPath)
	}()

	storage, _ := New(dbPath)

	// Создать новую заметку
	noteTitle := "Test Note"
	noteContent := "This is a test note."
	noteID, err := storage.NewNote(noteTitle, noteContent)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if noteID <= 0 {
		t.Errorf("Expected a positive note ID, got %d", noteID)
	}
}

func TestDeleteNote(t *testing.T) {
	dbPath := "test.db"
	defer func() {
		_ = os.Remove(dbPath)
	}()

	storage, _ := New(dbPath)

	noteTitle := "Test Note"
	noteContent := "This is a test note."
	noteID, _ := storage.NewNote(noteTitle, noteContent)

	deletedID, err := storage.DeleteNote(noteID)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if deletedID != noteID {
		t.Errorf("Expected deleted note ID to match the created note ID, got %d", deletedID)
	}
}

func TestSetNoteContent(t *testing.T) {
	dbPath := "test.db"
	defer func() {
		_ = os.Remove(dbPath)
	}()

	storage, _ := New(dbPath)

	noteTitle := "Test Note"
	noteContent := "This is a test note."
	noteID, _ := storage.NewNote(noteTitle, noteContent)

	newContent := "This is the updated content."
	err := storage.SetNoteContent(noteID, newContent)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Проверяем, что содержимое заметки было обновлено
	retrievedNote, err := storage.GetNoteByID(noteID)
	if err != nil {
		t.Errorf("Error retrieving note: %v", err)
	}

	if retrievedNote.GetContent() != newContent {
		t.Errorf("Expected updated content, got %s", retrievedNote.GetContent())
	}
}

func TestSearchNotesByKeyword(t *testing.T) {
	dbPath := "test.db"
	defer func() {
		_ = os.Remove(dbPath)
	}()

	storage, _ := New(dbPath)

	note1Title := "Test Note 1"
	note1Content := "This is the first test note."
	_, _ = storage.NewNote(note1Title, note1Content)

	note2Title := "Test Note 2"
	note2Content := "This is the second test note with a keyword."
	_, _ = storage.NewNote(note2Title, note2Content)

	keyword := "keyword"
	notes, err := storage.SearchNotesByKeyword(keyword)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(notes) != 1 {
		t.Errorf("Expected 1 matching note, got %d", len(notes))
	}

	if notes[0].GetTitle() != note2Title {
		t.Errorf("Expected matching note title, got %s", notes[0].GetTitle())
	}
}

func TestGetAllNotes(t *testing.T) {
	dbPath := "test.db"
	defer func() {
		_ = os.Remove(dbPath)
	}()

	storage, _ := New(dbPath)

	note1Title := "Test Note 1"
	note1Content := "This is the first test note."
	_, _ = storage.NewNote(note1Title, note1Content)

	note2Title := "Test Note 2"
	note2Content := "This is the second test note."
	_, _ = storage.NewNote(note2Title, note2Content)

	notes, err := storage.GetAllNotes()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(notes) != 2 {
		t.Errorf("Expected 2 notes, got %d", len(notes))
	}
}

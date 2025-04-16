package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateNote(t *testing.T) {
	clearNotes()
	var jsonStr = []byte(`{"content":"Test Note"}`)
	req, err := http.NewRequest("POST", "/notes", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/notes" {
			if r.Method == http.MethodPost {
				createNote(w, r)
			}
		}
	})
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	var note Note
	json.Unmarshal(rr.Body.Bytes(), &note)

	if note.Content != "Test Note" {
		t.Errorf("handler returned unexpected body: got %v want %v",
			note.Content, "Test Note")
	}
}

func TestListNotesEmpty(t *testing.T) {
	clearNotes()
	req, err := http.NewRequest("GET", "/notes", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/notes" {
			if r.Method == http.MethodGet {
				listNotes(w, r)
			}
		}
	})
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	var notes []Note
	json.Unmarshal(rr.Body.Bytes(), &notes)
	if len(notes) != 0 {
		t.Errorf("handler returned unexpected body: got %v want %v",
			len(notes), 0)
	}
}
func TestListNotes(t *testing.T) {
	clearNotes()
	createTestNote("Test Note 1")
	createTestNote("Test Note 2")
	req, err := http.NewRequest("GET", "/notes", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/notes" {
			if r.Method == http.MethodGet {
				listNotes(w, r)
			}
		}
	})
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	var notes []Note
	json.Unmarshal(rr.Body.Bytes(), &notes)
	if len(notes) != 2 {
		t.Errorf("handler returned unexpected body: got %v want %v",
			len(notes), 2)
	}
}

func TestGetNote(t *testing.T) {
	clearNotes()
	note := createTestNote("Test Note")
	req, err := http.NewRequest("GET", "/notes/"+note.ID, nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getNote)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var retrievedNote Note
	json.Unmarshal(rr.Body.Bytes(), &retrievedNote)
	if retrievedNote.ID != note.ID {
		t.Errorf("handler returned unexpected body: got %v want %v",
			retrievedNote.ID, note.ID)
	}
}

func TestGetNoteNotFound(t *testing.T) {
	clearNotes()
	req, err := http.NewRequest("GET", "/notes/nonexistent", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getNote)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}

func TestUpdateNote(t *testing.T) {
	clearNotes()
	note := createTestNote("Test Note")
	var jsonStr = []byte(`{"content":"Updated Test Note"}`)
	req, err := http.NewRequest("PUT", "/notes/"+note.ID, bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(updateNote)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var updatedNote Note
	json.Unmarshal(rr.Body.Bytes(), &updatedNote)

	if updatedNote.Content != "Updated Test Note" {
		t.Errorf("handler returned unexpected body: got %v want %v",
			updatedNote.Content, "Updated Test Note")
	}
}
func TestUpdateNoteNotFound(t *testing.T) {
	clearNotes()
	var jsonStr = []byte(`{"content":"Updated Test Note"}`)
	req, err := http.NewRequest("PUT", "/notes/nonexistent", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(updateNote)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}

}

func TestDeleteNote(t *testing.T) {
	clearNotes()
	note := createTestNote("Test Note")
	req, err := http.NewRequest("DELETE", "/notes/"+note.ID, nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(deleteNote)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNoContent)
	}
}

func TestDeleteNoteNotFound(t *testing.T) {
	clearNotes()
	req, err := http.NewRequest("DELETE", "/notes/nonexistent", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(deleteNote)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}
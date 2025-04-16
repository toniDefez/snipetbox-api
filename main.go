package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
)

type Note struct {
	ID      string `json:"id"`
	Content string `json:"content"`
}

func createTestNote(content string) Note {
	id := uuid.NewString()
	note := Note{
		ID:      id,
		Content: content,
	}
	notes[note.ID] = note
	return note
}

func clearNotes() {
	notes = make(map[string]Note)
}

var notes = make(map[string]Note)

func createNote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var newNote Note
	err := json.NewDecoder(r.Body).Decode(&newNote)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	newNote.ID = uuid.New().String()
	notes[newNote.ID] = newNote

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newNote)
}

func listNotes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var allNotes []Note
	for _, note := range notes {
		allNotes = append(allNotes, note)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allNotes)
}

func getNote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	id := parts[2]
	note, ok := notes[id]
	if !ok {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(note)
}

func updateNote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	id := parts[2]
	_, ok := notes[id]
	if !ok {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}

	var updatedNote Note
	err := json.NewDecoder(r.Body).Decode(&updatedNote)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	updatedNote.ID = id
	notes[id] = updatedNote

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedNote)
}

func deleteNote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	id := parts[2]
	_, ok := notes[id]
	if !ok {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}

	delete(notes, id)
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	http.HandleFunc("/notes", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			createNote(w, r)
		} else if r.Method == http.MethodGet {
			listNotes(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/notes/", func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) == 3 && parts[2] != "" {
			if r.Method == http.MethodGet {
				getNote(w, r)
			} else if r.Method == http.MethodPut {
				updateNote(w, r)
			} else if r.Method == http.MethodDelete {
				deleteNote(w, r)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		} else {
			http.Error(w, "Bad request", http.StatusBadRequest)
		}
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3002"
	}

	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
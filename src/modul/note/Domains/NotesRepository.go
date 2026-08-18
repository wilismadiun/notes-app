package domains

import "notes-app/src/modul/note/Domains/entities"

type NotesRepository interface {
	CreateNote(note entities.Note) error
	FindNoteById(id string) (entities.Note, error)
	EditNoteById(note entities.Note, update map[string]any) error
	DeleteNoteById(id string) error
}

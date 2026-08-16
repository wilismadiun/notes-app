package domains

import "notes-app/src/modul/note/Domains/entities"

type NotesRepository interface {
	CreateNote(note entities.Note) (entities.CreateNoteResponse, error)
	FindNoteById(id string) (entities.Note, error)
	DeleteNoteById(id string) (string, error)
	UpdateNoteById(note entities.Note) (entities.CreateNoteResponse, error)
}

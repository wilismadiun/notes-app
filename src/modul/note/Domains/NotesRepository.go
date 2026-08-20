package domains

import "notes-app/src/modul/note/Domains/entities"

type NotesRepository interface {
	CreateNote(note entities.Note) error
	GetAllNote(userId string) []entities.Note
	FindNoteById(note *entities.Note) error
	EditNoteById(note entities.Note, update map[string]any) error
	DeleteNoteById(note entities.Note) error
}

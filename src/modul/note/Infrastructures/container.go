package infrastructures

import (
	"notes-app/src/modul/note/Applications/usecase"
	"notes-app/src/modul/note/Infrastructures/generator"
	"notes-app/src/modul/note/Infrastructures/repositories"
	"notes-app/src/modul/note/Interfaces/http"

	"gorm.io/gorm"
)

func NoteContainer(db *gorm.DB) *http.NoteHandler {
	repoHandler := repositories.NoteRepository{DB: db}
	generatorHandler := generator.GeneratorIdUuid{}

	createNote := usecase.CreateNote{
		Repo:      &repoHandler,
		Generator: &generatorHandler,
	}

	getAllNote := usecase.GetAllNotes{
		Repo: &repoHandler,
	}

	getNoteById := usecase.GetNoteById{
		Repo: &repoHandler,
	}

	deleteNote := usecase.DeleteNote{
		Repo: &repoHandler,
	}

	editNote := usecase.EditNoteById{
		Repo: &repoHandler,
	}

	return &http.NoteHandler{
		CreateHandler:       &createNote,
		GetAllNoteshandler:  &getAllNote,
		GetNoteByIdhandler:  &getNoteById,
		EditNoteByIdHandler: &editNote,
		Deletehandler:       &deleteNote,
	}
}

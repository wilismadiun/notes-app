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

	CreateNote := usecase.CreateNote{
		Repo:      &repoHandler,
		Generator: &generatorHandler,
	}

	return &http.NoteHandler{
		CreateHandler: &CreateNote,
	}
}

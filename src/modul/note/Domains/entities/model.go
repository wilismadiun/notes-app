package entities

import "time"

type Note struct {
	ID       string
	Title    string `json:"title"`
	Content  string `json:"content"`
	CreateAt time.Time
	UpdateAt time.Time
	Owner    string `json:"owner"`
}

type CreateNoteResponse struct {
	ID    string
	Title string
}

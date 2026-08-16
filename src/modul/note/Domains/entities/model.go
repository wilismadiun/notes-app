package entities

import "time"

type Note struct {
	ID       string
	Title    string
	Content  string
	CreateAt time.Time
	UpdateAt time.Time
	Owner    string
}

package entities

import (
	"time"
)

type Note struct {
	ID           int
	Title        string
	Content      string
	CreatedAt    time.Time
	LastEditedAt time.Time
}

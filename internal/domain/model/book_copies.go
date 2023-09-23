package model

import "time"

type BookStatus int

const (
	InLibrary BookStatus = 0
	Borrowed  BookStatus = 1
	Lost      BookStatus = 2
)

type BookCopies struct {
	ID        int
	BookID    int
	Status    BookStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewBookCopies(bookID int, status BookStatus) BookCopies {
	return BookCopies{
		ID:     bookID,
		Status: status,
	}
}

package model

type BookStatus int

const (
	InLibrary BookStatus = 0
	Borrowed  BookStatus = 1
	Lost      BookStatus = 2
)

type BookCopies struct {
	ID     int
	BookID int
	Status BookStatus
}

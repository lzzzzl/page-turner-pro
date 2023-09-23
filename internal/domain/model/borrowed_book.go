package model

import "time"

type BorrowedBook struct {
	ID         int
	UserID     int
	CopyID     int
	BorrowDate time.Time
	DueDate    time.Time
	ReturnDate time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func NewBorrowdBook(userID, copyID int, borrowDate, dueDate, returnDate time.Time) BorrowedBook {
	return BorrowedBook{
		UserID:     userID,
		CopyID:     copyID,
		BorrowDate: borrowDate,
		DueDate:    dueDate,
		ReturnDate: returnDate,
	}
}

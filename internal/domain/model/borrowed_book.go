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

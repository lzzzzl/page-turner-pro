package model

import "time"

type Book struct {
	ID            int
	Title         string
	Author        string
	PublishedYear time.Time
	ISBN          string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func NewBook(title, author, isbn string, publishedYear time.Time) Book {
	return Book{
		Title:         title,
		Author:        author,
		ISBN:          isbn,
		PublishedYear: publishedYear,
	}
}

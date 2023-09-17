package model

import "time"

type User struct {
	ID        int
	UID       string
	Email     string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(uid, email, name string) User {
	return User{
		UID:   uid,
		Email: email,
		Name:  name,
	}
}

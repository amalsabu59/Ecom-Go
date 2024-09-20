package models

type User struct {
	ID       int64  `bun:",pk,autoincrement" json:"id,omitempty"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"` // Use omitempty to avoid exposing the password in responses
}

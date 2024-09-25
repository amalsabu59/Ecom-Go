package models

type User struct {
	ID        int64      `bun:",pk,autoincrement" json:"id,omitempty"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	IsActive  bool       `json:"is_active"`
	Password  string     `json:"password,omitempty"`
	Addresses []*Address `bun:"rel:has-many,join:id=user_id"`
}

type Address struct {
	ID        int64  `bun:",pk,autoincrement" json:"id,omitempty"`
	UserID    int64  `json:"user_id"` // Foreign key reference to User
	Street    string `json:"street"`
	City      string `json:"city"`
	State     string `json:"state"`
	ZipCode   string `json:"zip_code"`
	Country   string `json:"country"`
	IsPrimary bool   `json:"is_primary"` // Indicate if this is the primary address
}

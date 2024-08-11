package model

type User struct {
	ID     string  `json:"id"`
	Points float64 `json:"points"`
}

func NewUser(id string) *User {
	return &User{
		ID:     id,
		Points: 0.0,
	}
}

package model

type User struct {
	ID       string
	Email    string
	Username string
	Password string
	Salt     []byte
	Admin    bool
}

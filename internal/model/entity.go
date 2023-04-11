package model

type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
	Salt     []byte `json:"salt"`
	Admin    bool   `json:"admin"`
}

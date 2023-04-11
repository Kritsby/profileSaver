package model

type UserResponse struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
	Salt     []byte `json:"salt" swaggerignore:"true"`
	Admin    bool   `json:"admin" swaggerignore:"true"`
}

type UserRequest struct {
	ID       string `json:"id" swaggerignore:"true"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
	Salt     []byte `json:"salt" swaggerignore:"true"`
	Admin    bool   `json:"admin"`
}

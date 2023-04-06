package repository

import "dev/profileSaver/internal/model"

//go:generate mockgen -source=interfaces.go -destination=mocks/mock.go

type Repository interface {
	CreateUser(u model.User) error
	GetAllUsers() []model.User
	GetUserByName(name string) (model.User, error)
	GetUserByID(id string) (model.User, error)
	UpdateUser(u model.User) error
	DeleteUser(id string) error
	HashPass(password, salt []byte) ([]byte, []byte)
	CheckAdmin(id string) bool
}

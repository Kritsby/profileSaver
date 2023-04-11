package repository

import "dev/profileSaver/internal/model"

//go:generate mockgen -source=interfaces.go -destination=mocks/mock.go

type Repository interface {
	CreateUser(u model.UserResponse) error
	GetAllUsers() []model.UserResponse
	GetUserByName(name string) (model.UserResponse, error)
	GetUserByID(id string) (model.UserResponse, error)
	UpdateUser(u model.UserResponse) error
	DeleteUser(id string) error
	HashPass(password, salt []byte) ([]byte, []byte)
}

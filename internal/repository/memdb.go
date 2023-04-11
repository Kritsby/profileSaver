package repository

import (
	"crypto/rand"
	"dev/profileSaver/internal/model"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/argon2"
	"sync"
)

const (
	BusysUsername = "username is busy"
	UserNotFound  = "user not found"
)

type DB struct {
	mu     sync.RWMutex
	userId map[string]string
	store  map[string]model.UserResponse
}

func New() *DB {
	userId := make(map[string]string)
	store := make(map[string]model.UserResponse)

	userId["admin"] = "admin"
	store["admin"] = model.UserResponse{
		ID:       "admin",
		Email:    "admin",
		Username: "admin",
		Password: "admin",
		Salt:     nil,
		Admin:    true,
	}
	return &DB{
		mu:     sync.RWMutex{},
		userId: userId,
		store:  store,
	}
}

func (db *DB) CreateUser(u model.UserResponse) error {
	db.mu.RLock()
	defer db.mu.RUnlock()

	if _, ok := db.userId[u.Username]; ok {
		return errors.New(BusysUsername)
	}

	u.ID = uuid.New().String()

	hashedPass, salt := db.HashPass([]byte(u.Password), nil)

	u.Password = fmt.Sprintf("%x", hashedPass)
	u.Salt = salt

	db.userId[u.Username] = u.ID
	db.store[u.ID] = u

	return nil
}

func (db *DB) GetAllUsers() []model.UserResponse {
	db.mu.RLock()
	defer db.mu.RUnlock()

	users := make([]model.UserResponse, 0, len(db.store))
	for _, u := range db.store {
		users = append(users, u)
	}

	return users
}

func (db *DB) GetUserByName(name string) (model.UserResponse, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	id, ok := db.userId[name]
	if !ok {
		return model.UserResponse{}, errors.New(UserNotFound)
	}

	return db.store[id], nil
}

func (db *DB) GetUserByID(id string) (model.UserResponse, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	u, ok := db.store[id]
	if !ok {
		return model.UserResponse{}, errors.New(UserNotFound)
	}

	return u, nil
}

func (db *DB) UpdateUser(newUser model.UserResponse) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	oldUser, ok := db.store[newUser.ID]
	if !ok {
		return errors.New(UserNotFound)
	}

	user := db.updateUserFields(oldUser, newUser)

	if user.Username != oldUser.Username {
		if _, ok = db.userId[newUser.Username]; ok {
			return errors.New(BusysUsername)
		}

		delete(db.userId, oldUser.Username)
		db.userId[newUser.Username] = newUser.ID
	}

	db.store[newUser.ID] = newUser
	return nil
}

func (db *DB) DeleteUser(id string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	u, ok := db.store[id]
	if !ok {
		return errors.New(UserNotFound)
	}

	delete(db.userId, u.Username)
	delete(db.store, u.ID)

	return nil
}

func (db *DB) updateUserFields(oldUser, newUser model.UserResponse) model.UserResponse {
	if newUser.Username == "" {
		newUser.Username = oldUser.Username
	}

	if newUser.Email == "" {
		newUser.Email = oldUser.Email
	}

	if newUser.Password == "" {
		newUser.Password = oldUser.Password
	}

	hashedPass, _ := db.HashPass([]byte(newUser.Password), oldUser.Salt)

	newUser.Password = string(hashedPass)
	newUser.Salt = oldUser.Salt

	newUser.ID = oldUser.ID
	newUser.Admin = oldUser.Admin

	return newUser
}

func (db *DB) HashPass(password, salt []byte) ([]byte, []byte) {
	if salt == nil {
		salt = make([]byte, 8)
		rand.Read(salt)
	}
	hashedPass := argon2.IDKey(password, salt, 1, 64*1024, 4, 32)

	return hashedPass, salt
}

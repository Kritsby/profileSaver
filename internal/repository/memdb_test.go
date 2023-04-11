package repository

import (
	"dev/profileSaver/internal/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDB_CreateUser(t *testing.T) {
	db := New()

	db.userId["admin"] = ""

	tests := []struct {
		name        string
		expectedErr error
		input       model.User
	}{
		{
			name:        "OK",
			expectedErr: nil,
			input: model.User{
				ID:       "",
				Email:    "test@mail.ru",
				Username: "test",
				Password: "test",
				Salt:     nil,
				Admin:    false,
			},
		},
		{
			name:        "NOT_OK",
			expectedErr: ErrUserNameExists,
			input: model.User{
				ID:       "",
				Email:    "",
				Username: "admin",
				Password: "",
				Salt:     nil,
				Admin:    false,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actualErr := db.CreateUser(test.input)

			assert.Equal(t, test.expectedErr, actualErr)
		})
	}
}

func TestDB_GetAllUsers(t *testing.T) {
	db := New()

	db.store["admin"] = model.User{
		ID:       "admin",
		Email:    "admin",
		Username: "admin",
		Password: "admin",
		Salt:     nil,
		Admin:    true,
	}

	expected := []model.User{
		{
			ID:       "admin",
			Email:    "admin",
			Username: "admin",
			Password: "admin",
			Salt:     nil,
			Admin:    true,
		},
	}

	actual := db.GetAllUsers()

	assert.Equal(t, expected, actual)
}

func TestDB_GetUserByID(t *testing.T) {
	db := New()

	db.userId["admin"] = "admin"

	db.store["admin"] = model.User{
		ID:       "admin",
		Email:    "admin",
		Username: "admin",
		Password: "admin",
		Salt:     nil,
		Admin:    true,
	}

	tests := []struct {
		name        string
		expectedErr error
		expectedRes model.User
		input       string
	}{
		{
			name:        "OK",
			expectedErr: nil,
			expectedRes: model.User{
				ID:       "admin",
				Email:    "admin",
				Username: "admin",
				Password: "admin",
				Salt:     nil,
				Admin:    true,
			},
			input: "admin",
		},
		{
			name:        "NOT_OK",
			expectedErr: ErrUserNotFound,
			expectedRes: model.User{},
			input:       "1",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actualRes, actualErr := db.GetUserByID(test.input)

			assert.Equal(t, test.expectedRes, actualRes)
			assert.Equal(t, test.expectedErr, actualErr)
		})
	}
}

func TestDB_GetUserByName(t *testing.T) {
	db := New()

	db.userId["admin"] = "admin"

	db.store["admin"] = model.User{
		ID:       "admin",
		Email:    "admin",
		Username: "admin",
		Password: "admin",
		Salt:     nil,
		Admin:    true,
	}

	tests := []struct {
		name        string
		expectedErr error
		expectedRes model.User
		input       string
	}{
		{
			name:        "OK",
			expectedErr: nil,
			expectedRes: model.User{
				ID:       "admin",
				Email:    "admin",
				Username: "admin",
				Password: "admin",
				Salt:     nil,
				Admin:    true,
			},
			input: "admin",
		},
		{
			name:        "NOT_OK",
			expectedErr: ErrUserNotFound,
			expectedRes: model.User{},
			input:       "1",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actualRes, actualErr := db.GetUserByName(test.input)

			assert.Equal(t, test.expectedRes, actualRes)
			assert.Equal(t, test.expectedErr, actualErr)
		})
	}
}

func TestDB_UpdateUser(t *testing.T) {
	db := New()

	db.userId["admin"] = "admin"
	db.store["admin"] = model.User{
		ID:       "admin",
		Email:    "admin",
		Username: "admin",
		Password: "admin",
		Salt:     nil,
		Admin:    true,
	}

	db.userId["test"] = "admin"
	db.store["test"] = model.User{
		ID:       "test",
		Email:    "test",
		Username: "test",
		Password: "test",
		Salt:     nil,
		Admin:    true,
	}

	tests := []struct {
		name        string
		expectedErr error
		input       model.User
	}{
		{
			name:        "OK",
			expectedErr: nil,
			input: model.User{
				ID:       "test",
				Email:    "",
				Username: "random",
				Password: "",
			},
		},
		{
			name:        "NOT_FOUND",
			expectedErr: ErrUserNotFound,
			input: model.User{
				ID:       "1",
				Email:    "",
				Username: "",
				Password: "",
			},
		},
		{
			name:        "BUSY_USER_NAME",
			expectedErr: ErrUserNameExists,
			input: model.User{
				ID:       "test",
				Email:    "",
				Username: "admin",
				Password: "",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actualErr := db.UpdateUser(test.input)

			assert.Equal(t, test.expectedErr, actualErr)
		})
	}
}

func TestDB_DeleteUser(t *testing.T) {
	db := New()

	db.userId["test"] = "admin"
	db.store["test"] = model.User{
		ID:       "test",
		Email:    "test",
		Username: "test",
		Password: "test",
		Salt:     nil,
		Admin:    true,
	}

	tests := []struct {
		name        string
		input       string
		expectedErr error
	}{
		{
			name:        "OK",
			input:       "test",
			expectedErr: nil,
		},
		{
			name:        "USER_NOT_FOUND",
			input:       "1",
			expectedErr: ErrUserNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actualErr := db.DeleteUser(test.input)

			assert.Equal(t, test.expectedErr, actualErr)
		})
	}
}

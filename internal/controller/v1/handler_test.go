package v1

import (
	"bytes"
	"dev/profileSaver/internal/model"
	mock_repository "dev/profileSaver/internal/repository/mocks"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_handler(t *testing.T) {
	type mockBehavior func(s *mock_repository.MockRepository)

	tests := []struct {
		name                 string
		isAdmin              bool
		handler              string
		method               string
		inputBody            string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:    "OK",
			handler: "CreateUser",
			method:  "POST",
			isAdmin: true,
			mockBehavior: func(s *mock_repository.MockRepository) {
				s.EXPECT().CreateUser(model.User{
					Email:    "test@mail.ru",
					Username: "test",
					Password: "test",
				}).Return(nil)
			},
			inputBody:          `{"email":"test@mail.ru", "username":"test", "password":"test"}`,
			expectedStatusCode: 200,
			expectedResponseBody: `{"data":"user was created"}
`,
		},
		{
			name:    "NOT_OK",
			handler: "CreateUser",
			method:  "POST",
			isAdmin: true,
			mockBehavior: func(s *mock_repository.MockRepository) {
				s.EXPECT().CreateUser(model.User{
					Email:    "test@mail.ru",
					Username: "test",
					Password: "test",
				}).Return(errors.New("error"))
			},
			inputBody:          `{"email":"test@mail.ru", "username":"test", "password":"test"}`,
			expectedStatusCode: 500,
			expectedResponseBody: `{"error":{}}
`,
		},
		{
			name:    "NOT_ADMIN",
			handler: "CreateUser",
			method:  "POST",
			isAdmin: false,
			mockBehavior: func(s *mock_repository.MockRepository) {

			},
			inputBody:            `{"email":"test@mail.ru", "username":"test", "password":"test"}`,
			expectedStatusCode:   401,
			expectedResponseBody: "",
		},
		{
			name:               "BAD_REQUEST",
			handler:            "CreateUser",
			method:             "POST",
			isAdmin:            true,
			mockBehavior:       func(s *mock_repository.MockRepository) {},
			inputBody:          `{1}`,
			expectedStatusCode: 400,
			expectedResponseBody: `{"error":{"Offset":2}}
`,
		},
		{
			name:    "OK",
			method:  "GET",
			handler: "GetAllUsers",
			mockBehavior: func(s *mock_repository.MockRepository) {
				s.EXPECT().GetAllUsers().Return([]model.User{})
			},
			expectedStatusCode: 200,
			expectedResponseBody: `{"data":[]}
`,
		},
		{
			name:    "OK",
			method:  "GET",
			handler: "GetUser",
			mockBehavior: func(s *mock_repository.MockRepository) {
				s.EXPECT().GetUserByID("1").Return(model.User{}, nil)
			},
			expectedStatusCode: 200,
			expectedResponseBody: `{"data":{"id":"","email":"","username":"","password":"","salt":null,"admin":false}}
`,
		},
		{
			name:    "NOT_OK",
			method:  "GET",
			handler: "GetUser",
			mockBehavior: func(s *mock_repository.MockRepository) {
				s.EXPECT().GetUserByID("1").Return(model.User{}, errors.New("error"))
			},
			expectedStatusCode: 500,
			expectedResponseBody: `{"error":{}}
`,
		},
		{
			name:    "OK",
			method:  "PATCH",
			handler: "UpdateUser",
			isAdmin: true,
			mockBehavior: func(s *mock_repository.MockRepository) {
				s.EXPECT().UpdateUser(model.User{
					ID:       "1",
					Email:    "test@mail.ru",
					Username: "test",
					Password: "test",
				}).Return(nil)
			},
			inputBody:          `{"email":"test@mail.ru", "username":"test", "password":"test"}`,
			expectedStatusCode: 200,
			expectedResponseBody: `{"data":"user was updated"}
`,
		},
		{
			name:               "NOT_OK",
			method:             "PATCH",
			handler:            "UpdateUser",
			isAdmin:            true,
			mockBehavior:       func(s *mock_repository.MockRepository) {},
			inputBody:          `{}`,
			expectedStatusCode: 400,
			expectedResponseBody: `{"error":{}}
`,
		},
		{
			name:                 "NOT_ADMIN",
			method:               "PATCH",
			handler:              "UpdateUser",
			isAdmin:              false,
			mockBehavior:         func(s *mock_repository.MockRepository) {},
			inputBody:            `{}`,
			expectedStatusCode:   401,
			expectedResponseBody: "",
		},
		{
			name:               "NOT_ADMIN",
			method:             "PATCH",
			handler:            "UpdateUser",
			isAdmin:            true,
			mockBehavior:       func(s *mock_repository.MockRepository) {},
			inputBody:          `{1}`,
			expectedStatusCode: 400,
			expectedResponseBody: `{"error":{"Offset":2}}
`,
		},
		{
			name:    "OK",
			method:  "DELETE",
			handler: "DeleteUser",
			isAdmin: true,
			mockBehavior: func(s *mock_repository.MockRepository) {
				s.EXPECT().DeleteUser("1").Return(nil)
			},
			expectedStatusCode: 200,
			expectedResponseBody: `{"data":"user was deleted"}
`,
		},
		{
			name:    "NOT_OK",
			method:  "DELETE",
			handler: "DeleteUser",
			isAdmin: true,
			mockBehavior: func(s *mock_repository.MockRepository) {
				s.EXPECT().DeleteUser("1").Return(errors.New("error"))
			},
			expectedStatusCode: 500,
			expectedResponseBody: `{"error":{}}
`,
		},
		{
			name:                 "NOT_ADMIN",
			method:               "DELETE",
			handler:              "DeleteUser",
			isAdmin:              false,
			mockBehavior:         func(s *mock_repository.MockRepository) {},
			expectedStatusCode:   401,
			expectedResponseBody: "",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_repository.NewMockRepository(c)
			repo.EXPECT().IsAuthorized("admin", "admin").Return(true)
			testCase.mockBehavior(repo)

			handlers := New(repo)

			r := handlers.InitRouter()
			var w *httptest.ResponseRecorder
			var req *http.Request
			switch testCase.handler {
			case "GetAllUsers":
				w = httptest.NewRecorder()
				req = httptest.NewRequest("GET", "/v1/user",
					nil)
				req.SetBasicAuth("admin", "admin")
			case "CreateUser":
				repo.EXPECT().GetUserByName("admin").Return(model.User{Admin: testCase.isAdmin}, nil)
				w = httptest.NewRecorder()
				req = httptest.NewRequest("POST", "/v1/user",
					bytes.NewBufferString(testCase.inputBody))
				req.SetBasicAuth("admin", "admin")
			case "UpdateUser":
				repo.EXPECT().GetUserByName("admin").Return(model.User{Admin: testCase.isAdmin}, nil)
				w = httptest.NewRecorder()
				req = httptest.NewRequest("PATCH", "/v1/user/1",
					bytes.NewBufferString(testCase.inputBody))
				req.SetBasicAuth("admin", "admin")
			case "GetUser":
				w = httptest.NewRecorder()
				req = httptest.NewRequest("GET", "/v1/user/1",
					nil)
				req.SetBasicAuth("admin", "admin")
			case "DeleteUser":
				repo.EXPECT().GetUserByName("admin").Return(model.User{Admin: testCase.isAdmin}, nil)
				w = httptest.NewRecorder()
				req = httptest.NewRequest("DELETE", "/v1/user/1",
					nil)
				req.SetBasicAuth("admin", "admin")
			}

			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}

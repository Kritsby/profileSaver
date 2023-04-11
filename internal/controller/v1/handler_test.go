package v1

import (
	"bytes"
	"context"
	"dev/profileSaver/internal/model"
	mock_repository "dev/profileSaver/internal/repository/mocks"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bunrouter"
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
				s.EXPECT().CreateUser(model.UserResponse{}).Return(nil)
			},
			inputBody:          `{}`,
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
				s.EXPECT().CreateUser(model.UserResponse{}).Return(errors.New("error"))
			},
			inputBody:          `{}`,
			expectedStatusCode: 500,
			expectedResponseBody: `{"error":{}}
`,
		},
		{
			name:               "NOT_ADMIN",
			handler:            "CreateUser",
			method:             "POST",
			isAdmin:            false,
			mockBehavior:       func(s *mock_repository.MockRepository) {},
			inputBody:          `{}`,
			expectedStatusCode: 511,
			expectedResponseBody: `{"error":"don't have permission"}
`,
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
				s.EXPECT().GetAllUsers().Return([]model.UserResponse{})
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
				s.EXPECT().GetUserByID("1").Return(model.UserResponse{}, nil)
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
				s.EXPECT().GetUserByID("1").Return(model.UserResponse{}, errors.New("error"))
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
				s.EXPECT().UpdateUser(model.UserResponse{}).Return(nil)
			},
			inputBody:          `{}`,
			expectedStatusCode: 200,
			expectedResponseBody: `{"data":"user was updated"}
`,
		},
		{
			name:    "NOT_OK",
			method:  "PATCH",
			handler: "UpdateUser",
			isAdmin: true,
			mockBehavior: func(s *mock_repository.MockRepository) {
				s.EXPECT().UpdateUser(model.UserResponse{}).Return(errors.New("error"))
			},
			inputBody:          `{}`,
			expectedStatusCode: 500,
			expectedResponseBody: `{"error":{}}
`,
		},
		{
			name:               "NOT_ADMIN",
			method:             "PATCH",
			handler:            "UpdateUser",
			isAdmin:            false,
			mockBehavior:       func(s *mock_repository.MockRepository) {},
			inputBody:          `{}`,
			expectedStatusCode: 511,
			expectedResponseBody: `{"error":"don't have permission"}
`,
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
			name:               "NOT_ADMIN",
			method:             "DELETE",
			handler:            "DeleteUser",
			isAdmin:            false,
			mockBehavior:       func(s *mock_repository.MockRepository) {},
			expectedStatusCode: 511,
			expectedResponseBody: `{"error":"don't have permission"}
`,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_repository.NewMockRepository(c)
			testCase.mockBehavior(repo)

			handlers := New(repo)

			r := bunrouter.New()
			var w *httptest.ResponseRecorder
			var req *http.Request
			switch testCase.handler {
			case "GetAllUsers":
				r.GET("/v1/user", handlers.getAllUsers)

				w = httptest.NewRecorder()
				req = httptest.NewRequest("GET", "/v1/user",
					nil)
			case "CreateUser":
				r.POST("/v1/user", handlers.createUser)

				ctx := context.WithValue(context.Background(), "is_admin", true)

				if !testCase.isAdmin {
					ctx = context.WithValue(context.Background(), "is_admin", false)
				}

				w = httptest.NewRecorder()
				req = httptest.NewRequest("POST", "/v1/user",
					bytes.NewBufferString(testCase.inputBody)).WithContext(ctx)
			case "UpdateUser":
				r.PATCH("/v1/user", handlers.updateUser)

				ctx := context.WithValue(context.Background(), "is_admin", true)

				if !testCase.isAdmin {
					ctx = context.WithValue(context.Background(), "is_admin", false)
				}

				w = httptest.NewRecorder()
				req = httptest.NewRequest("PATCH", "/v1/user",
					bytes.NewBufferString(testCase.inputBody)).WithContext(ctx)
			case "GetUser":
				r.GET("/v1/user/:id", handlers.getUser)

				w = httptest.NewRecorder()
				req = httptest.NewRequest("GET", "/v1/user/1",
					nil)
			case "DeleteUser":
				r.DELETE("/v1/user/:id", handlers.deleteUser)

				ctx := context.WithValue(context.Background(), "is_admin", true)

				if !testCase.isAdmin {
					ctx = context.WithValue(context.Background(), "is_admin", false)
				}

				w = httptest.NewRecorder()
				req = httptest.NewRequest("DELETE", "/v1/user/1",
					nil).WithContext(ctx)
			}

			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}

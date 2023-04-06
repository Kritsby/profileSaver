package v1

import (
	"bytes"
	"context"
	"dev/profileSaver/internal/model"
	mock_repository "dev/profileSaver/internal/repository/mocks"
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
			mockBehavior: func(s *mock_repository.MockRepository) {
				s.EXPECT().CreateUser(model.User{}).Return(nil)
			},
			inputBody:          `{}`,
			expectedStatusCode: 200,
			expectedResponseBody: `{"data":"user was created","params":null,"route":"/v1/user"}
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
			expectedResponseBody: `{"data":[],"params":null,"route":"/v1/user"}
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
			expectedResponseBody: `{"data":{"ID":"","Email":"","Username":"","Password":"","Salt":null,"Admin":false},"params":{"id":"1"},"route":"/v1/user/:id"}
`,
		},
		{
			name:    "OK",
			method:  "PATCH",
			handler: "UpdateUser",
			mockBehavior: func(s *mock_repository.MockRepository) {
				s.EXPECT().UpdateUser(model.User{}).Return(nil)
			},
			inputBody:          `{}`,
			expectedStatusCode: 200,
			expectedResponseBody: `{"data":"user was updated","params":null,"route":"/v1/user"}
`,
		},
		{
			name:    "OK",
			method:  "DELETE",
			handler: "DeleteUser",
			mockBehavior: func(s *mock_repository.MockRepository) {
				s.EXPECT().DeleteUser("1").Return(nil)
			},
			expectedStatusCode: 200,
			expectedResponseBody: `{"data":"user was deleted","params":{"id":"1"},"route":"/v1/user/:id"}
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

				w = httptest.NewRecorder()
				req = httptest.NewRequest("POST", "/v1/user",
					bytes.NewBufferString(testCase.inputBody)).WithContext(ctx)
			case "UpdateUser":
				r.PATCH("/v1/user", handlers.updateUser)

				ctx := context.WithValue(context.Background(), "is_admin", true)

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

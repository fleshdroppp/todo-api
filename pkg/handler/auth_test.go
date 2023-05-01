package handler

import (
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/magiconair/properties/assert"
	"net/http/httptest"
	"testing"
	todo "todo-api"
	"todo-api/pkg/service"
	mock_service "todo-api/pkg/service/mocks"
)

func TestHandler_signUp(t *testing.T) {
	type mockBehaviour func(s *mock_service.MockAuthorization, user todo.User)

	testTable := []struct {
		name                string
		inputBody           string
		inputUser           todo.User
		mockBehaviour       mockBehaviour
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			inputBody: `{"name":"Test","username":"test","password":"qwerty"}`,
			inputUser: todo.User{
				Name:     "Test",
				Username: "test",
				Password: "qwerty",
			},
			mockBehaviour: func(s *mock_service.MockAuthorization, user todo.User) {
				s.EXPECT().CreateUser(user).Return(1, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"id":1}`,
		},
		{
			name:                "Empty Fields",
			inputBody:           `{"username":"test","password":"qwerty"}`,
			mockBehaviour:       func(s *mock_service.MockAuthorization, user todo.User) {},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"invalid input body"}`,
		},
		{
			name:      "Service failure",
			inputBody: `{"name":"Test","username":"test","password":"qwerty"}`,
			inputUser: todo.User{
				Name:     "Test",
				Username: "test",
				Password: "qwerty",
			},
			mockBehaviour: func(s *mock_service.MockAuthorization, user todo.User) {
				s.EXPECT().CreateUser(user).Return(1, errors.New("service failure"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"message":"service failure"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			testCase.mockBehaviour(auth, testCase.inputUser)

			services := &service.Service{Authorization: auth}
			handler := NewHandler(services)

			r := gin.New()
			r.POST("/sign-up", handler.signUp)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/sign-up",
				bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}

}

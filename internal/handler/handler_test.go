package handler

import (
	"bytes"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/julienschmidt/httprouter"
	"github.com/muesli/termenv"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http/httptest"
	"nprn/internal/entity/user/usermodel"
	"nprn/internal/service"
	mock_service "nprn/internal/service/mocks"
	"nprn/pkg/logging"
	"testing"
)

func TestMain(m *testing.M) {
	code := m.Run()
	if code != 0 {
		printWarning("WARNING! TESTS IS FAIL")
	} else {
		printSuccess("OK! ALL TESTS IS PASS")
	}
}

func TestHandler_SignUp(t *testing.T) {
	type mockBehavior func(storage *mock_service.MockUserStorage, user usermodel.UserInternal)

	token, _ := service.GenerateToken("1")
	passHash, _ := service.GeneratePasswordHash("AnnaTestPass")

	testTable := []struct {
		name                string
		inputBody           string
		inputUser           usermodel.UserInternal
		mockBehavior        mockBehavior
		exceptedStatusCode  int
		exceptedRequestBody string
	}{
		{
			name:      "OK",
			inputBody: `{"username":"AnnaTest", "password":"AnnaTestPass", "email":"test@test.com"}`,
			inputUser: usermodel.UserInternal{
				Username:     "AnnaTest",
				PasswordHash: passHash,
				Email:        "test@test.com"},
			mockBehavior: func(storage *mock_service.MockUserStorage, user usermodel.UserInternal) {
				storage.EXPECT().Create(gomock.Any(), user).Return("1", nil)
			},
			exceptedStatusCode:  200,
			exceptedRequestBody: fmt.Sprintf(`{"token":"%s"}`, token),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			userStorage := mock_service.NewMockUserStorage(c)
			testCase.mockBehavior(userStorage, testCase.inputUser)

			logger := logging.GetLogger()

			testService := service.NewService(userStorage, nil, logger)
			testHandler := NewHandler(testService, logger)

			router := httprouter.New()

			router.POST("/auth/sign-up", testHandler.CheckErrorMiddleware(testHandler.SignUp))

			recorder := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/auth/sign-up", bytes.NewBufferString(testCase.inputBody))
			router.ServeHTTP(recorder, req)

			assert.Equal(t, testCase.exceptedStatusCode, recorder.Code)
			assert.Equal(t, testCase.exceptedRequestBody, recorder.Body.String())

		})
	}
}

func TestHandler_SignIn(t *testing.T) {
	type mockBehavior func(storage *mock_service.MockUserStorage, username string, password string)

	token, _ := service.GenerateToken("1")
	passAnna, _ := service.GeneratePasswordHash("AnnaTestPass")

	testTable := []struct {
		name                string
		inputBody           string
		inputUsername       string
		inputPassword       string
		mockBehavior        mockBehavior
		exceptedStatusCode  int
		exceptedRequestBody string
	}{
		{
			name:          "OK",
			inputBody:     `{"username":"AnnaTest", "password":"AnnaTestPass"}`,
			inputUsername: "AnnaTest",
			inputPassword: passAnna,
			mockBehavior: func(storage *mock_service.MockUserStorage, username string, password string) {
				storage.EXPECT().GetOne(gomock.Any(), username, password).Return(usermodel.UserTransfer{ID: "1", Username: "AnnaTest"}, nil)
			},
			exceptedStatusCode:  200,
			exceptedRequestBody: fmt.Sprintf(`{"token":"%s"}`, token),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			userStorage := mock_service.NewMockUserStorage(c)
			testCase.mockBehavior(userStorage, testCase.inputUsername, testCase.inputPassword)

			logger := logging.GetLogger()

			testService := service.NewService(userStorage, nil, logger)
			testHandler := NewHandler(testService, logger)

			router := httprouter.New()

			router.POST("/auth/sign-in", testHandler.CheckErrorMiddleware(testHandler.SignIn))

			recorder := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/auth/sign-in", bytes.NewBufferString(testCase.inputBody))

			router.ServeHTTP(recorder, req)

			assert.Equal(t, testCase.exceptedStatusCode, recorder.Code)
			assert.Equal(t, testCase.exceptedRequestBody, recorder.Body.String())
		})
	}
}

func printWarning(message string) {
	profile := termenv.ColorProfile()

	str := termenv.String(message)
	styleStr := str.Foreground(profile.Color("#ff0000"))

	log.Println(styleStr)
}

func printSuccess(message string) {
	profile := termenv.ColorProfile()

	str := termenv.String(message)
	styleStr := str.Foreground(profile.Color("#0dff00"))

	log.Println(styleStr)
}

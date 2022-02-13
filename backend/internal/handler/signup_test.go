package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Kara4ev/go-web-tmp/internal/model"
	"github.com/Kara4ev/go-web-tmp/internal/model/apperrors"
	"github.com/Kara4ev/go-web-tmp/internal/model/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSignup(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	baseURL := "/api/account"
	url := fmt.Sprintf("%s/signup", baseURL)

	t.Run("Not json data", func(t *testing.T) {
		mockUserService := new(mocks.MockUserService)
		mockUserService.On("Signup", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("*model.User")).Return(nil)

		rr := httptest.NewRecorder()
		router := gin.Default()

		NewHandler(&Config{
			Router:      router,
			UserService: mockUserService,
			BaseUrl:     baseURL,
		})

		reqBody := []byte("plantext")

		request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "text/plan")
		router.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusUnsupportedMediaType, rr.Code)

		mockUserService.AssertNotCalled(t, "Signup")
	})

	t.Run("Email and password required", func(t *testing.T) {

		mockUserService := new(mocks.MockUserService)
		mockUserService.On("Signup", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("*model.User")).Return(nil)

		rr := httptest.NewRecorder()
		router := gin.Default()

		NewHandler(&Config{
			Router:      router,
			UserService: mockUserService,
			BaseUrl:     baseURL,
		})

		reqBody, err := json.Marshal(gin.H{
			"email": "",
		})

		assert.NoError(t, err)

		request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusBadRequest, rr.Code)

		mockUserService.AssertNotCalled(t, "Signup")
	})

	t.Run("Email incorrect", func(t *testing.T) {

		mockUserService := new(mocks.MockUserService)
		mockUserService.On("Signup", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("*model.User")).Return(nil)

		rr := httptest.NewRecorder()
		router := gin.Default()

		NewHandler(&Config{
			Router:      router,
			UserService: mockUserService,
			BaseUrl:     baseURL,
		})

		reqBody, err := json.Marshal(gin.H{
			"email":    "incorrect-email",
			"password": "correct-password",
		})

		assert.NoError(t, err)

		request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusBadRequest, rr.Code)

		mockUserService.AssertNotCalled(t, "Signup")
	})

	t.Run("Password in short", func(t *testing.T) {

		mockUserService := new(mocks.MockUserService)
		mockUserService.On("Signup", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("*model.User")).Return(nil)

		rr := httptest.NewRecorder()
		router := gin.Default()

		NewHandler(&Config{
			Router:      router,
			UserService: mockUserService,
			BaseUrl:     baseURL,
		})

		reqBody, err := json.Marshal(gin.H{
			"email":    "incorrect-email",
			"password": "pas",
		})

		assert.NoError(t, err)

		request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusBadRequest, rr.Code)

		mockUserService.AssertNotCalled(t, "Signup")
	})

	t.Run("Password in long", func(t *testing.T) {

		mockUserService := new(mocks.MockUserService)
		mockUserService.On("Signup", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("*model.User")).Return(nil)

		rr := httptest.NewRecorder()
		router := gin.Default()

		NewHandler(&Config{
			Router:      router,
			UserService: mockUserService,
			BaseUrl:     baseURL,
		})

		reqBody, err := json.Marshal(gin.H{
			"email":    "incorrect-email",
			"password": "long-password-1111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111",
		})

		assert.NoError(t, err)

		request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusBadRequest, rr.Code)

		mockUserService.AssertNotCalled(t, "Signup")
	})

	t.Run("App error returned from user services", func(t *testing.T) {

		apperror := apperrors.NewConflict("User already exists", "correct@email.com")

		mockUserService := new(mocks.MockUserService)
		mockUserService.
			On("Signup", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("*model.User")).
			Return(apperror)

		rr := httptest.NewRecorder()
		router := gin.Default()

		NewHandler(&Config{
			Router:      router,
			UserService: mockUserService,
			BaseUrl:     baseURL,
		})

		reqBody, err := json.Marshal(gin.H{
			"email":    "correct@email.com",
			"password": "correct-password",
		})

		assert.NoError(t, err)

		request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusConflict, rr.Code)

		mockUserService.AssertExpectations(t)
	})

	t.Run("Other error returned from user services", func(t *testing.T) {

		mockUserService := new(mocks.MockUserService)
		mockUserService.On("Signup",
			mock.AnythingOfType("*context.emptyCtx"),
			mock.AnythingOfType("*model.User")).Return(fmt.Errorf("Some error down call chain"))

		rr := httptest.NewRecorder()
		router := gin.Default()

		NewHandler(&Config{
			Router:      router,
			UserService: mockUserService,
			BaseUrl:     baseURL,
		})

		reqBody, err := json.Marshal(gin.H{
			"email":    "correct@email.com",
			"password": "correct-password",
		})

		assert.NoError(t, err)

		request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)

		mockUserService.AssertExpectations(t)
	})

	t.Run("Successful token creation", func(t *testing.T) {

		u := &model.User{
			Email:    "correct@email.com",
			Password: "correct-password",
		}

		mockTokenResp := &model.TokenPair{
			IDToken:      model.IDToken{SS: "idToken"},
			RefreshToken: model.RefreshToken{SS: "refreshToken"},
		}

		mockUserService := new(mocks.MockUserService)
		mockTokenService := new(mocks.MockTokenService)

		mockUserService.
			On("Signup", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("*model.User")).
			Return(nil)

		mockTokenService.
			On("NewPairFromUser", mock.AnythingOfType("*context.emptyCtx"), u, "").
			Return(mockTokenResp, nil)

		rr := httptest.NewRecorder()
		router := gin.Default()

		NewHandler(&Config{
			Router:       router,
			UserService:  mockUserService,
			TokenService: mockTokenService,
			BaseUrl:      baseURL,
		})

		reqBody, err := json.Marshal(gin.H{
			"email":    u.Email,
			"password": u.Password,
		})

		assert.NoError(t, err)

		request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(gin.H{
			"tokens": mockTokenResp,
		})
		assert.NoError(t, err)

		assert.Equal(t, http.StatusCreated, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockUserService.AssertExpectations(t)
		mockTokenService.AssertExpectations(t)
	})

	t.Run("Failed token creation", func(t *testing.T) {

		u := &model.User{
			Email:    "correct@email.com",
			Password: "correct-password",
		}

		mockErrorRespon := apperrors.NewInternal()

		mockUserService := new(mocks.MockUserService)
		mockTokenService := new(mocks.MockTokenService)

		mockUserService.
			On("Signup", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("*model.User")).
			Return(nil)

		mockTokenService.
			On("NewPairFromUser", mock.AnythingOfType("*context.emptyCtx"), u, "").
			Return(nil, mockErrorRespon)

		rr := httptest.NewRecorder()
		router := gin.Default()

		NewHandler(&Config{
			Router:       router,
			UserService:  mockUserService,
			TokenService: mockTokenService,
			BaseUrl:      baseURL,
		})

		reqBody, err := json.Marshal(gin.H{
			"email":    u.Email,
			"password": u.Password,
		})

		assert.NoError(t, err)

		request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(gin.H{
			"error": mockErrorRespon,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockErrorRespon.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockUserService.AssertExpectations(t)
		mockTokenService.AssertExpectations(t)
	})

}

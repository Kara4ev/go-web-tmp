package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Kara4ev/go-web-tmp/internal/model"
	"github.com/Kara4ev/go-web-tmp/internal/model/mocks"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSignout(t *testing.T) {

	gin.SetMode(gin.TestMode)

	baseURL := "/api/account"
	url := fmt.Sprintf("%s/signout", baseURL)

	t.Run("Success", func(t *testing.T) {

		uid, _ := uuid.NewRandom()

		mockTokenService := new(mocks.MockTokenService)
		mockTokenService.
			On("Signout", mock.AnythingOfType("*context.emptyCtx"), uid).
			Return(nil)

		rr := httptest.NewRecorder()
		router := gin.Default()

		router.Use(func(c *gin.Context) {
			c.Set("user", &model.User{
				UID: uid,
			},
			)
		})

		NewHandler(&Config{
			Router:       router,
			TokenService: mockTokenService,
			BaseUrl:      baseURL,
		})

		request, err := http.NewRequest(http.MethodPost, url, nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		mockTokenService.AssertExpectations(t)
	})

	t.Run("error -> return error h.TokenService.Signout", func(t *testing.T) {

		uid, _ := uuid.NewRandom()

		mockTokenService := new(mocks.MockTokenService)
		mockTokenService.
			On("Signout", mock.AnythingOfType("*context.emptyCtx"), uid).
			Return(fmt.Errorf("Some error down call chain"))

		rr := httptest.NewRecorder()
		router := gin.Default()

		router.Use(func(c *gin.Context) {
			c.Set("user", &model.User{
				UID: uid,
			},
			)
		})

		NewHandler(&Config{
			Router:       router,
			TokenService: mockTokenService,
			BaseUrl:      baseURL,
		})

		request, err := http.NewRequest(http.MethodPost, url, nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		mockTokenService.AssertExpectations(t)
	})

	t.Run("error -> not exist user context", func(t *testing.T) {

		mockTokenService := new(mocks.MockTokenService)
		mockTokenService.
			On("Signout", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("*uuid.UUID")).
			Return(nil)

		rr := httptest.NewRecorder()
		router := gin.Default()

		NewHandler(&Config{
			Router:       router,
			TokenService: mockTokenService,
			BaseUrl:      baseURL,
		})

		request, err := http.NewRequest(http.MethodPost, url, nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		mockTokenService.AssertNotCalled(t, "Signout", mock.Anything, mock.Anything)
	})

}

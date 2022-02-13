package handler

import (
	"time"

	"github.com/Kara4ev/go-web-tmp/internal/handler/middleware"
	"github.com/Kara4ev/go-web-tmp/internal/model"
	"github.com/Kara4ev/go-web-tmp/internal/model/apperrors"
	"github.com/Kara4ev/go-web-tmp/pkg/logger"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	UserService  model.UserService
	TokenService model.TokenService
}

type Config struct {
	Router          *gin.Engine
	UserService     model.UserService
	TokenService    model.TokenService
	BaseUrl         string
	TimeoutDuration time.Duration
}

func NewHandler(c *Config) {

	h := &Handler{
		UserService:  c.UserService,
		TokenService: c.TokenService,
	}

	timeoutDuration := c.TimeoutDuration
	if timeoutDuration == 0 {
		timeoutDuration = 5 * time.Minute
	}

	g := c.Router.Group(c.BaseUrl)
	logger.Debug("Gin mode: %s", gin.Mode())
	if gin.Mode() != gin.TestMode {
		g.Use(middleware.Timeout(timeoutDuration, apperrors.NewServiceUnavailable()))
		g.GET("/me", middleware.AuthUser(h.TokenService), h.Me)
		g.POST("/signout", middleware.AuthUser(h.TokenService), h.Signout)
		g.PUT("/details", middleware.AuthUser(h.TokenService), h.Details)
	} else {
		g.GET("/me", h.Me)
		g.POST("/signout", h.Signout)
		g.PUT("/details", h.Details)
	}

	g.POST("/signin", h.Signin)
	g.POST("/signup", h.Signup)
	g.POST("/tokens", h.Tokens)

}

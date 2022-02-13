package middleware

import (
	"strings"

	"github.com/Kara4ev/go-web-tmp/internal/model"
	"github.com/Kara4ev/go-web-tmp/internal/model/apperrors"
	"github.com/Kara4ev/go-web-tmp/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

//import "github.com/go-playground/validator/v10"

type authHeader struct {
	IDToken string `header:"Authorization"`
}

type invalidArgument struct {
	Field string `json:"field"`
	Value string `json:"value"`
	Tag   string `json:"tag"`
	Param string `json:"param"`
}

func AuthUser(s model.TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Debug("middleware AuthUser: execute")
		h := new(authHeader)
		if err := c.ShouldBindHeader(&h); err != nil {
			logger.Debug("middleware AuthUser: error validate request")
			if errs, ok := err.(validator.ValidationErrors); ok {
				var invalidArgs []invalidArgument

				for _, err := range errs {
					invalidArgs = append(invalidArgs, invalidArgument{
						err.Field(),
						err.Value().(string),
						err.Tag(),
						err.Param(),
					})
				}
				err := apperrors.NewBadRequest("invalid request parametrs. See invalidArgs")

				c.JSON(err.Status(), gin.H{
					"error":       err,
					"invalidArgs": invalidArgs,
				})
				c.Abort()
				return
			}

			err := apperrors.NewInternal()
			c.JSON(err.Status(), gin.H{
				"error": err,
			})
			c.Abort()
			return
		}
		logger.Debug("middleware AuthUser: request valid")
		idTokenHeader := strings.Split(h.IDToken, "Bearer ")
		if len(idTokenHeader) != 2 {
			logger.Debug("middleware AuthUser: error token format")
			err := apperrors.NewAuthorization("Must provide Authorization header with format `Bearer {token}`")
			c.JSON(err.Status(), gin.H{
				"error": err,
			})
			c.Abort()
			return
		}
		logger.Debug("middleware AuthUser: bearer format valid")
		user, err := s.ValidateIDToken(idTokenHeader[1])
		if err != nil {
			logger.Debug("middleware AuthUser execute: error token validate")
			err := apperrors.NewAuthorization("Provided token is invalid")
			c.JSON(err.Status(), gin.H{
				"error": err,
			})
			c.Abort()
			return
		}
		logger.Debug("middleware AuthUser: token valide, user uid: %s email: %s", user.UID.String(), user.Email)
		c.Set("user", user)
		c.Next()
	}
}

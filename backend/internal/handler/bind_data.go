package handler

import (
	"fmt"

	"github.com/Kara4ev/go-web-tmp/internal/model/apperrors"
	"github.com/Kara4ev/go-web-tmp/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type invalidArgument struct {
	Filed  string `json:"filed"`
	Valuse string `json:"valuse"`
	Tag    string `json:"tag"`
	Param  string `json:"param"`
}

func bindData(c *gin.Context, req interface{}) bool {
	if c.ContentType() != "application/json" {
		msg := fmt.Sprintf("%s only accepts Content-Type application/json", c.FullPath())
		logger.Warn(msg)
		err := apperrors.NewUnsupportedMediaType(msg)
		c.JSON(err.Status(), gin.H{
			"error": err,
		})
		return false
	}

	if err := c.ShouldBind(req); err != nil {
		logger.Warn("error binding data: %v", err)
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

			err := apperrors.NewBadRequest("invalid request parameters, see invalidArgs")
			c.JSON(err.Status(), gin.H{
				"error":       err,
				"invalidArgs": invalidArgs,
			})
			return false
		}

		fallBack := apperrors.NewInternal()
		c.JSON(fallBack.Status(), gin.H{
			"error": fallBack,
		})
		return false

	}

	return true
}

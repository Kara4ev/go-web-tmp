package handler

import (
	"net/http"

	"github.com/Kara4ev/go-web-tmp/internal/model"
	"github.com/Kara4ev/go-web-tmp/internal/model/apperrors"
	"github.com/Kara4ev/go-web-tmp/pkg/logger"
	"github.com/gin-gonic/gin"
)

func (h *Handler) Signout(c *gin.Context) {

	user, exists := c.Get("user")
	if !exists {
		logger.Error("Unable to extract user from request context for unknown reason: %v\n", c)
		err := apperrors.NewInternal()
		c.JSON(err.Status(), gin.H{
			"error": err,
		})

		return
	}

	ctx := c.Request.Context()
	if err := h.TokenService.Signout(ctx, user.(*model.User).UID); err != nil {
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "user signet out successfully",
	})
}

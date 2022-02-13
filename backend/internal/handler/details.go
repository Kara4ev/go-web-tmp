package handler

import (
	"net/http"

	"github.com/Kara4ev/go-web-tmp/internal/model"
	"github.com/Kara4ev/go-web-tmp/internal/model/apperrors"
	"github.com/Kara4ev/go-web-tmp/pkg/logger"
	"github.com/gin-gonic/gin"
)

type detailsReq struct {
	Name  string `json:"name" binding:"omitempty,max=5 0"`
	Email string `json:"email" binding:"omitempty,email"`
}

// Details handler
func (h *Handler) Details(c *gin.Context) {
	var req detailsReq

	if ok := bindData(c, &req); !ok {
		return
	}

	authUser, exists := c.Get("user")
	if !exists {
		logger.Error("Unable to extract user from request context for unknown reason: %v", c)
		err := apperrors.NewInternal()
		c.JSON(err.Status(), gin.H{
			"error": err,
		})

		return
	}

	ctx := c.Request.Context()
	user := &model.User{
		UID:   authUser.(model.User).UID,
		Email: req.Email,
		Name:  req.Name,
	}

	if err := h.UserService.UpdateDetails(ctx, user); err != nil {
		logger.Error("Failed to update user details: %v", err)
		err := apperrors.NewInternal()
		c.JSON(err.Status(), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})

}

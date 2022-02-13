package handler

import (
	"net/http"

	"github.com/Kara4ev/go-web-tmp/internal/model/apperrors"
	"github.com/Kara4ev/go-web-tmp/pkg/logger"
	"github.com/gin-gonic/gin"
)

type tokenReq struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// Tokens handler
func (h *Handler) Tokens(c *gin.Context) {

	var req tokenReq

	if ok := bindData(c, &req); !ok {
		return
	}

	ctx := c.Request.Context()

	refreshToken, err := h.TokenService.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	u, err := h.UserService.Get(ctx, refreshToken.UID)
	if err != nil {
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	tokens, err := h.TokenService.NewPairFromUser(ctx, u, refreshToken.ID.String())

	if err != nil {
		logger.Warn("failed to create tokens for user %+v, error: %v", u, err)
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tokens": tokens,
	})

}

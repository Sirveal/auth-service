package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) getTokenPair(c *gin.Context) {
	userUUID := c.Query("uuid")
	if userUUID == "" {
		newErrorResponse(c, http.StatusBadRequest, "не указан параметр uuid")
		return
	}

	clientIP := c.ClientIP()
	tokens, err := h.services.Authorization.GenerateTokenPair(userUUID, clientIP)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, tokens)
}

func (h *Handler) refreshTokens(c *gin.Context) {
	var input struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	logrus.Debugf("получен запрос на обновление токена: %s", input.RefreshToken[:10])

	clientIP := c.ClientIP()
	tokens, err := h.services.Authorization.RefreshTokens(input.RefreshToken, clientIP)
	if err != nil {
		logrus.Errorf("не удалось обновить токены: %v", err)
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.JSON(http.StatusOK, tokens)
}

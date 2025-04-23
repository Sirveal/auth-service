package handler

import (
	"helloapp/pkg/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	auth := router.Group("/auth")
	{
		auth.GET("/token", h.getTokenPair)
		auth.POST("/refresh", h.refreshTokens)
	}

	return router
}

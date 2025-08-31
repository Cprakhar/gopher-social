package handler

import (
	"net/http"

	"github.com/cprakhar/gopher-social/internal/config"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	Cfg config.Config
}

func (h *Handler) HealthCheckHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"env":     h.Cfg.Env,
		"version": h.Cfg.Version,
	})
}

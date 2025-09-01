package handler

import (
	"net/http"

	"github.com/cprakhar/gopher-social/internal/config"
	"github.com/cprakhar/gopher-social/internal/store"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	Cfg   config.Config
	Store store.Store
}

func (h *Handler) HealthCheckHandler(ctx *gin.Context) {
	data := map[string]string{
		"status":  "ok",
		"env":     h.Cfg.Env,
		"version": h.Cfg.Version,
	}
	writeJSON(ctx, http.StatusOK, data)
}

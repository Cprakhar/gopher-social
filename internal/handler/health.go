package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//	@BasePath	/v1

// HealthCheck godoc
//	@Summary	health check
//	@Schemes
//	@Description	get the health status
//	@Tags			health
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	map[string]string
//	@Router			/health [get]
func (h *Handler) HealthCheckHandler(ctx *gin.Context) {
	data := map[string]string{
		"status":  "ok",
		"env":     h.Cfg.Env,
		"version": h.Cfg.Version,
	}
	writeJSON(ctx, http.StatusOK, data)
}

package handler

import "github.com/gin-gonic/gin"

type Handler struct {
}

func (h *Handler) HealthCheckHandler(ctx *gin.Context) {
	ctx.Writer.Write([]byte("ok"))
}

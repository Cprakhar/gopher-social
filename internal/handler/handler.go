package handler

import (
	"net/http"

	"github.com/cprakhar/gopher-social/internal/config"
	"github.com/cprakhar/gopher-social/internal/store"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	Cfg    config.Config
	Store  store.Store
	Logger *zap.SugaredLogger
}

func writeJSON(ctx *gin.Context, status int, data any) {
	type envelope struct {
		Data any `json:"data"`
	}

	ctx.JSON(status, envelope{Data: data})
}

func (h *Handler) internalServerErr(ctx *gin.Context, err error) {
	h.Logger.Errorw("internal server error", "method", ctx.Request.Method, "path", ctx.Request.URL.Path, "error", err.Error())
	ctx.JSON(http.StatusInternalServerError, gin.H{"error": "the server encountered a problem"})
}

func (h *Handler) badRequestErr(ctx *gin.Context, err error) {
	h.Logger.Warnw("bad request error", "method", ctx.Request.Method, "path", ctx.Request.URL.Path, "error", err.Error())
	ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
}

func (h *Handler) notFoundErr(ctx *gin.Context, err error) {
	h.Logger.Errorw("not found error", "method", ctx.Request.Method, "path", ctx.Request.URL.Path, "error", err.Error())
	ctx.JSON(http.StatusNotFound, gin.H{"error": "resource not found"})
}

func (h *Handler) conflictErr(ctx *gin.Context, err error) {
	h.Logger.Errorw("conflict error", "method", ctx.Request.Method, "path", ctx.Request.URL.Path, "error", err.Error())
	ctx.JSON(http.StatusConflict, gin.H{"error": "resource already exists"})
}

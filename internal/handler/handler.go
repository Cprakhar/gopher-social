package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func writeJSON(ctx *gin.Context, status int, data any) {
	type envelope struct {
		Data any `json:"data"`
	}

	ctx.JSON(status, envelope{Data: data})
}

func internalServerErr(ctx *gin.Context, err error) {
	log.Printf("internal server error: %s path: %s error: %s", ctx.Request.Method, ctx.Request.URL.Path, err.Error())
	ctx.JSON(http.StatusInternalServerError, gin.H{"error": "the server encountered a problem"})
}

func badRequestErr(ctx *gin.Context, err error) {
	log.Printf("bad request error: %s path: %s error: %s", ctx.Request.Method, ctx.Request.URL.Path, err.Error())

	ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
}

func notFoundErr(ctx *gin.Context, err error) {
	log.Printf("not found error: %s path: %s error: %s", ctx.Request.Method, ctx.Request.URL.Path, err.Error())

	ctx.JSON(http.StatusNotFound, gin.H{"error": "resource not found"})
}
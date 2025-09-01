package handler

import (
	"net/http"

	"github.com/cprakhar/gopher-social/internal/store"
	"github.com/gin-gonic/gin"
)

type CreateUserPayload struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) RegisterUserHandler(ctx *gin.Context) {
	var payload CreateUserPayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := &store.User{
		Username: payload.Username,
		Email:    payload.Email,
		Password: payload.Password,
	}

	if err := h.Store.Users.Create(ctx, user); err != nil {
		internalServerErr(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"id":         user.ID,
		"username":   user.Username,
		"email":      user.Email,
		"created_at": user.CreatedAt,
	})
}

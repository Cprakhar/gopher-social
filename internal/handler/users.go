package handler

import (
	"errors"
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
		Password: []byte(payload.Password),
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

func (h *Handler) GetUserHandler(ctx *gin.Context) {
	user := userFromCtx(ctx)

	writeJSON(ctx, http.StatusOK, user)
}

type FollowUserPayload struct {
	UserID string `json:"user_id" binding:"required"`
}

func (h *Handler) FollowUserHandler(ctx *gin.Context) {
	followingID := userFromCtx(ctx)

	// Get the user ID from auth context or session (stubbed here)
	var payload FollowUserPayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		badRequestErr(ctx, err)
		return
	}

	if err := h.Store.Followers.Follow(ctx, payload.UserID, followingID.ID); err != nil {
		switch {
		case errors.Is(err, store.ErrConflict):
			conflictErr(ctx, err)
			return
		default:
			internalServerErr(ctx, err)
			return
		}
	}

	ctx.Status(http.StatusCreated)
}

func (h *Handler) UnfollowUserHandler(ctx *gin.Context) {
	followingID := userFromCtx(ctx)

	// Get the user ID from auth context or session (stubbed here)
	var payload FollowUserPayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		badRequestErr(ctx, err)
		return
	}

	if err := h.Store.Followers.Unfollow(ctx, payload.UserID, followingID.ID); err != nil {
		internalServerErr(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (h *Handler) UsersContextMiddleware(ctx *gin.Context) {
	id := ctx.Param("id")
	user, err := h.Store.Users.GetByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			notFoundErr(ctx, err)
			ctx.Abort()
			return
		default:
			internalServerErr(ctx, err)
			ctx.Abort()
			return
		}
	}

	ctx.Set("user", user)
	ctx.Next()
}

func userFromCtx(ctx *gin.Context) *store.User {
	user, ok := ctx.Get("user")
	if !ok {
		return nil
	}
	return user.(*store.User)
}
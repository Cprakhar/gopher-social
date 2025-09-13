package handler

import (
	"errors"
	"net/http"

	"github.com/cprakhar/gopher-social/internal/store"
	"github.com/gin-gonic/gin"
)

// GetUser godoc
//
//	@Summary	get a user
//	@Schemes
//	@Description	get a user by id
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"user id"
//	@Success		200	{object}	store.User
//	@Security		ApiKeyAuth
//	@Router			/users/{id} [get]
func (h *Handler) GetUserHandler(ctx *gin.Context) {
	user := userFromCtx(ctx)

	writeJSON(ctx, http.StatusOK, user)
}

// FollowUser godoc
//
//	@Summary	follow a user
//	@Schemes
//	@Description	follow a user by id
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string				true	"user id to follow"
//	@Param			payload	body		FollowUserPayload	true	"follow user payload"
//	@Success		201		{object}	nil
//	@Failure		400		{object}	map[string]string
//	@Failure		409		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Security		ApiKeyAuth
//	@Router			/users/{id}/follow [post]
func (h *Handler) FollowUserHandler(ctx *gin.Context) {
	user := userFromCtx(ctx)
	followingID := ctx.Param("id")

	if err := h.Store.Followers.Follow(ctx, followingID, user.ID); err != nil {
		switch {
		case errors.Is(err, store.ErrConflict):
			h.conflictErr(ctx, err)
			return
		default:
			h.internalServerErr(ctx, err)
			return
		}
	}

	ctx.Status(http.StatusCreated)
}

// UnfollowUser godoc
//
//	@Summary	unfollow a user
//	@Schemes
//	@Description	unfollow a user by id
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string				true	"user id to unfollow"
//	@Param			payload	body		FollowUserPayload	true	"unfollow user payload"
//	@Success		204		{object}	nil
//	@Failure		400		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Security		ApiKeyAuth
//	@Router			/users/{id}/unfollow [post]
func (h *Handler) UnfollowUserHandler(ctx *gin.Context) {
	user := userFromCtx(ctx)
	followingID := ctx.Param("id")

	if err := h.Store.Followers.Unfollow(ctx, followingID, user.ID); err != nil {
		h.internalServerErr(ctx, err)
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
			h.notFoundErr(ctx, err)
			ctx.Abort()
			return
		default:
			h.internalServerErr(ctx, err)
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

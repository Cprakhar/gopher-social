package handler

import (
	"errors"
	"net/http"

	"github.com/cprakhar/gopher-social/internal/store"
	"github.com/gin-gonic/gin"
)

type CreatePostPayload struct {
	Title   string   `json:"title" binding:"required"`
	Content string   `json:"content" binding:"required"`
	Tags    []string `json:"tags"`
}

func (h *Handler) CreatePostHandler(ctx *gin.Context) {
	var payload CreatePostPayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		internalServerErr(ctx, err)
		return
	}

	authorID := "1644b5d2-67c6-4ead-9938-76245ec7b68e"

	post := &store.Post{
		Title:    payload.Title,
		Content:  payload.Content,
		AuthorID: authorID,
		Tags:     payload.Tags,
	}

	if err := h.Store.Posts.Create(ctx, post); err != nil {
		internalServerErr(ctx, err)
		return
	}

	writeJSON(ctx, http.StatusCreated, post)
}

func (h *Handler) GetPostHandler(ctx *gin.Context) {
	post := postFromCtx(ctx)

	comments, err := h.Store.Comments.GetByPostID(ctx, post.ID)
	if err != nil {
		internalServerErr(ctx, err)
		return
	}

	post.Comments = comments

	writeJSON(ctx, http.StatusOK, post)
}

func (h *Handler) DeletePostHandler(ctx *gin.Context) {
	post := postFromCtx(ctx)

	if err := h.Store.Posts.Delete(ctx, post.ID); err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			notFoundErr(ctx, err)
			return
		default:
			internalServerErr(ctx, err)
			return
		}
	}

	ctx.Status(http.StatusNoContent)
}

type UpdatePostPayload struct {
	Title   *string  `json:"title,omitempty"`
	Content *string  `json:"content,omitempty"`
	Tags    []string `json:"tags,omitempty"`
}

func (h *Handler) UpdatePostHandler(ctx *gin.Context) {
	var payload UpdatePostPayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		badRequestErr(ctx, err)
		return
	}

	post := postFromCtx(ctx)

	if payload.Title != nil {
		post.Title = *payload.Title
	}
	if payload.Content != nil {
		post.Content = *payload.Content
	}
	if payload.Tags != nil {
		post.Tags = payload.Tags
	}

	if err := h.Store.Posts.Update(ctx, post); err != nil {
		internalServerErr(ctx, err)
		return
	}

	writeJSON(ctx, http.StatusOK, post)
}

func (h *Handler) PostsContextMiddleware(ctx *gin.Context) {
	id := ctx.Param("id")
	post, err := h.Store.Posts.GetByID(ctx, id)
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

	ctx.Set("post", post)
	ctx.Next()
}

func postFromCtx(ctx *gin.Context) *store.Post {
	post, ok := ctx.Get("post")
	if !ok {
		return nil
	}
	return post.(*store.Post)
}

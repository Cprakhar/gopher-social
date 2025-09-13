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

// CreatePost godoc
//	@Summary	create a post
//	@Schemes
//	@Description	create a new post
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		CreatePostPayload	true	"post payload"
//	@Success		201		{object}	store.Post
//	@Failure		500		{object}	map[string]string
//	@Security		ApiKeyAuth
//	@Router			/posts [post]

func (h *Handler) CreatePostHandler(ctx *gin.Context) {
	var payload CreatePostPayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		h.internalServerErr(ctx, err)
		return
	}

	author := userFromCtx(ctx)

	post := &store.Post{
		Title:    payload.Title,
		Content:  payload.Content,
		AuthorID: author.ID,
		Tags:     payload.Tags,
	}

	if err := h.Store.Posts.Create(ctx, post); err != nil {
		h.internalServerErr(ctx, err)
		return
	}

	writeJSON(ctx, http.StatusCreated, post)
}

// GetPost godoc
//
//	@Summary	get a post
//	@Schemes
//	@Description	get a post by id
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"post id"
//	@Success		200	{object}	store.Post
//	@Failure		500	{object}	map[string]string
//	@Security		ApiKeyAuth
//	@Router			/posts/{id} [get]
func (h *Handler) GetPostHandler(ctx *gin.Context) {
	post := postFromCtx(ctx)

	comments, err := h.Store.Comments.GetByPostID(ctx, post.ID)
	if err != nil {
		h.internalServerErr(ctx, err)
		return
	}

	post.Comments = comments

	writeJSON(ctx, http.StatusOK, post)
}

// DeletePost godoc
//
//	@Summary	delete a post
//	@Schemes
//	@Description	delete a post by id
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"post id"
//	@Success		204	"No Content"
//	@Failure		404	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Security		ApiKeyAuth
//	@Router			/posts/{id} [delete]
func (h *Handler) DeletePostHandler(ctx *gin.Context) {
	post := postFromCtx(ctx)

	if err := h.Store.Posts.Delete(ctx, post.ID); err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			h.notFoundErr(ctx, err)
			return
		default:
			h.internalServerErr(ctx, err)
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

// UpdatePost godoc
//
//	@Summary	update a post
//	@Schemes
//	@Description	update a post by id
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string				true	"post id"
//	@Param			payload	body		UpdatePostPayload	true	"post payload"
//	@Success		200		{object}	store.Post
//	@Failure		400		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Security		ApiKeyAuth
//	@Router			/posts/{id} [patch]
func (h *Handler) UpdatePostHandler(ctx *gin.Context) {
	var payload UpdatePostPayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		h.badRequestErr(ctx, err)
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
		h.internalServerErr(ctx, err)
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
			h.notFoundErr(ctx, err)
			ctx.Abort()
			return
		default:
			h.internalServerErr(ctx, err)
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

package handler

import (
	"net/http"

	"github.com/cprakhar/gopher-social/internal/store"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// GetUserFeed godoc
//
//	@Summary	get user feed
//	@Schemes
//	@Description	get a user's feed with pagination
//	@Tags			feed
//	@Accept			json
//	@Produce		json
//	@Param			limit	query		int		false	"number of posts to return"	default(20)
//	@Param			offset	query		int		false	"number of posts to skip"	default(0)
//	@Param			sort	query		string	false	"sort order"				Enums(asc, desc)	default(desc)
//	@Success		200		{object}	[]store.PostWithMetadata
//	@Failure		404		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Security		ApiKeyAuth
//	@Router			/feed [get]
func (h *Handler) GetUserFeedHandler(ctx *gin.Context) {
	fp := store.PaginatedFeedQuery{
		Limit:  20,
		Offset: 0,
		Sort:   "desc",
	}

	fp, err := fp.Parse(ctx)
	if err != nil {
		h.badRequestErr(ctx, err)
		return
	}

	if err := validator.New().Struct(fp); err != nil {
		h.badRequestErr(ctx, err)
		return
	}

	if err := ctx.ShouldBindQuery(&fp); err != nil {
		h.badRequestErr(ctx, err)
		return
	}

	// Get the user ID from auth context or session (stubbed here)
	user := userFromCtx(ctx)

	feed, err := h.Store.Posts.GetUserFeed(ctx, user.ID, fp)
	if err != nil {
		h.internalServerErr(ctx, err)
		return
	}

	writeJSON(ctx, http.StatusOK, feed)
}

package handler

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/cprakhar/gopher-social/internal/mail"
	"github.com/cprakhar/gopher-social/internal/store"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type CreateUserPayload struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// RegisterUser godoc
//
//	@Summary	register a new user
//	@Schemes
//	@Description	register a new user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		CreateUserPayload	true	"user payload"
//	@Success		201		{object}	map[string]any
//	@Failure		500		{object}	map[string]string
//	@Router			/users [post]
func (h *Handler) RegisterUserHandler(ctx *gin.Context) {
	var payload CreateUserPayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := &store.User{
		Username: payload.Username,
		Email:    payload.Email,
	}

	if err := user.Password.Set(payload.Password); err != nil {
		h.internalServerErr(ctx, err)
		return
	}

	plainToken := uuid.NewString()
	hash := sha256.Sum256([]byte(plainToken))
	hashToken := hex.EncodeToString(hash[:])

	if err := h.Store.Users.CreateAndInvite(ctx, user, hashToken, h.Cfg.Mail.Exp); err != nil {
		h.internalServerErr(ctx, err)
		return
	}

	isProdEnv := h.Cfg.Env == "production"

	activationURL := h.Cfg.WebURL + "/confirm/" + plainToken
	vars := struct {
		Username      string
		ActivationURL string
	}{
		Username:      user.Username,
		ActivationURL: activationURL,
	}
	status, err := h.Mailer.Send(mail.UserWelcomeTemplate, user.Username, user.Email, vars, !isProdEnv)
	if err != nil {
		h.Logger.Errorw("error sending welcome email", "error", err)

		if err := h.Store.Users.Delete(ctx, user.ID); err != nil {
			h.Logger.Errorw("error deleting user after email failure", "error", err)
		}
		return
	}

	h.Logger.Infow("email sent", "status code", status)

	writeJSON(ctx, http.StatusCreated, map[string]any{
		"id":    user.ID,
		"email": user.Email,
		"token": plainToken, // In real-world apps, send this via email
	})
}

// ActivateUser godoc
//
//	@Summary	activate a user
//	@Schemes
//	@Description	activate a user using the token sent via email
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			token	path		string	true	"activation token"
//	@Success		204		{object}	nil
//	@Failure		404		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/users/activate/{token} [put]
func (h *Handler) ActivateUserHandler(ctx *gin.Context) {
	token := ctx.Param("token")

	if err := h.Store.Users.Activate(ctx, token); err != nil {
		switch err {
		case store.ErrNotFound:
			h.notFoundErr(ctx, err)
			return
		default:
			h.internalServerErr(ctx, err)
			return
		}
	}

	ctx.Status(http.StatusNoContent)
}

type CreateUserTokenPayload struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// CreateToken godoc
//
//	@Summary	create a JWT token
//	@Schemes
//	@Description	create a JWT token using email and password
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		CreateUserTokenPayload	true	"credentials payload"
//	@Success		201		{object}	string
//	@Failure		400		{object}	map[string]string
//	@Failure		401		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/tokens [post]
func (h *Handler) CreateTokenHandler(ctx *gin.Context) {
	// parse credentials payload
	var payload CreateUserTokenPayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		h.badRequestErr(ctx, err)
		return
	}

	// fetch the user (check if user exists) from the payload
	user, err := h.Store.Users.GetByEmail(ctx, payload.Email)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			h.unauthorizedErr(ctx, err)
			return
		default:
			h.internalServerErr(ctx, err)
			return
		}
	}
	// generate a token -> add claims

	claims := jwt.RegisteredClaims{
		Issuer:    h.Cfg.Auth.Token.Iss,
		Subject:   user.ID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(h.Cfg.Auth.Token.Exp)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		Audience: jwt.ClaimStrings{
			h.Cfg.Auth.Token.Aud,
		},
		ID: uuid.NewString(),
	}

	tokenStr, err := h.Authenticator.GenerateToken(claims)
	if err != nil {
		h.internalServerErr(ctx, err)
		return
	}

	// send it to the client
	writeJSON(ctx, http.StatusCreated, tokenStr)
}

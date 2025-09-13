package handler

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/cprakhar/gopher-social/internal/store"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func (h *Handler) BasicAuthMiddleware(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		h.unauthorizedBasicErr(ctx, fmt.Errorf("authorization header is missing"))
		ctx.Abort()
		return
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Basic" {
		h.unauthorizedBasicErr(ctx, fmt.Errorf("authorization header is malformed"))
		ctx.Abort()
		return
	}

	decoded, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		h.unauthorizedBasicErr(ctx, err)
		ctx.Abort()
		return
	}

	username := h.Cfg.Auth.Basic.Username
	password := h.Cfg.Auth.Basic.Password

	creds := strings.SplitN(string(decoded), ":", 2)
	if len(creds) != 2 || creds[0] != username || creds[1] != password {
		h.unauthorizedBasicErr(ctx, fmt.Errorf("invalid credentials"))
		ctx.Abort()
		return
	}

	ctx.Next()
}

func (h *Handler) AuthTokenMiddleware(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		h.unauthorizedErr(ctx, fmt.Errorf("authorization header is missing"))
		ctx.Abort()
		return
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		h.unauthorizedErr(ctx, fmt.Errorf("authorization header is malformed"))
		ctx.Abort()
		return
	}

	token, err := h.Authenticator.ValidateToken(parts[1])
	if err != nil {
		h.unauthorizedErr(ctx, err)
		ctx.Abort()
		return
	}

	claims, _ := token.Claims.(jwt.RegisteredClaims)

	userID := claims.Subject
	user, err := h.Store.Users.GetByID(ctx, userID)
	if err != nil {
		h.unauthorizedErr(ctx, err)
		ctx.Abort()
		return
	}

	ctx.Set("user", user)
	ctx.Next()
}

func (h *Handler) CheckPostOwnership(requiredRole string, next gin.HandlerFunc) gin.HandlerFunc {
	return gin.HandlerFunc(func(ctx *gin.Context) {
		user := userFromCtx(ctx)
		post := postFromCtx(ctx)

		if post.AuthorID == user.ID {
			next(ctx)
			return
		}

		allowed, err := h.checkRolePrecedence(ctx, user, requiredRole)
		if err != nil {
			h.internalServerErr(ctx, err)
			ctx.Abort()
			return
		}
		if !allowed {
			h.forbiddenErr(ctx)
			ctx.Abort()
			return
		}
		next(ctx)
	})
}

func (h *Handler) checkRolePrecedence(ctx context.Context, user *store.User, roleName string) (bool, error) {
	role, err := h.Store.Roles.GetByName(ctx, roleName)
	if err != nil {
		return false, err
	}

	return user.Role.Level >= role.Level, nil
}

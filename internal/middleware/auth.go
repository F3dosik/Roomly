package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/avito-internships/test-backend-1-F3dosik/internal/ctxkey"
	"github.com/avito-internships/test-backend-1-F3dosik/internal/domain"
	"github.com/avito-internships/test-backend-1-F3dosik/internal/httputil"
	"github.com/avito-internships/test-backend-1-F3dosik/internal/jwt"
	"go.uber.org/zap"
)

func RequireAuth(logger *zap.SugaredLogger, secretKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				logger.Debug("auth: missing token")
				httputil.HandleError(w, httputil.NewAppError(
					httputil.ErrCodeUnauthorized,
					"unauthorized",
					http.StatusUnauthorized,
				))
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := jwt.ParseToken(tokenString, secretKey)
			if err != nil {
				logger.Debugw("auth error", "error", err)
				httputil.HandleError(w, httputil.NewAppError(
					httputil.ErrCodeUnauthorized,
					"unauthorized",
					http.StatusUnauthorized,
				))
				return
			}

			ctx := context.WithValue(r.Context(), ctxkey.UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, ctxkey.RoleKey, claims.Role)

			next.ServeHTTP(w, r.WithContext(ctx))

		})
	}
}

func RequireRole(role domain.Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRole := r.Context().Value(ctxkey.RoleKey)
			if userRole != role {
				httputil.HandleError(w, httputil.NewAppError(
					httputil.ErrCodeForbidden,
					"access denied",
					http.StatusForbidden,
				))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

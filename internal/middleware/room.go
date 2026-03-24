package middleware

import (
	"context"
	"net/http"

	"github.com/avito-internships/test-backend-1-F3dosik/internal/ctxkey"
	"github.com/avito-internships/test-backend-1-F3dosik/internal/httputil"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func CheckRoomID(logger *zap.SugaredLogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			idStr := chi.URLParam(r, "roomId")
			if idStr == "" {
				logger.Debug("check roomId: roomId is empty")
				httputil.HandleError(w, httputil.NewAppError(
					httputil.ErrCodeInvalidRequest,
					"roomId is empty",
					http.StatusBadRequest,
				))
				return
			}
			id, err := uuid.Parse(idStr)
			if err != nil {
				logger.Debugw("check roomId", "error", err)
				httputil.HandleError(w, httputil.NewAppError(
					httputil.ErrCodeInvalidRequest,
					"invalid roomId",
					http.StatusBadRequest,
				))
				return
			}

			ctx := context.WithValue(r.Context(), ctxkey.RoomIDKey, id)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

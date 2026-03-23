package handler

import (
	"net/http"

	"github.com/avito-internships/test-backend-1-F3dosik/internal/middleware"
	"github.com/avito-internships/test-backend-1-F3dosik/internal/middleware/gzip"
	"github.com/avito-internships/test-backend-1-F3dosik/internal/service"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Handler struct {
	router      chi.Router
	userService service.UserService
	jwtSecret   string
	logger      *zap.SugaredLogger
}

func New(secretKey string, us service.UserService, logger *zap.SugaredLogger) *Handler {
	h := &Handler{
		router:      chi.NewRouter(),
		userService: us,
		logger:      logger,
		jwtSecret:   secretKey,
	}
	h.setupMiddleware()
	h.setupRoutes()
	return h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

func (h *Handler) setupMiddleware() {
	h.router.Use(gzip.WithCompression(h.logger))
	h.router.Use(middleware.WithLogging(h.logger))
}

func (h *Handler) setupRoutes() {
	h.router.Get("/_info", h.Info)

	h.router.Group(func(r chi.Router) {
		r.Use(middleware.RequireJSON(h.logger))
		r.Post("/dummyLogin", h.dummyLogin)
		r.Post("/register", h.register)
	})

}

func (h *Handler) Info(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

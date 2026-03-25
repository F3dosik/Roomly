package handler

import (
	"net/http"

	"github.com/avito-internships/test-backend-1-F3dosik/internal/domain"
	"github.com/avito-internships/test-backend-1-F3dosik/internal/middleware"
	"github.com/avito-internships/test-backend-1-F3dosik/internal/middleware/gzip"
	"github.com/avito-internships/test-backend-1-F3dosik/internal/service"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Handler struct {
	router         chi.Router
	userService    service.UserService
	roomService    service.RoomService
	slotService    service.SlotService
	bookingService service.BookingService
	jwtSecret      string
	logger         *zap.SugaredLogger
}

func New(
	secretKey string,
	us service.UserService, rs service.RoomService,
	ss service.SlotService, bs service.BookingService,
	logger *zap.SugaredLogger,
) *Handler {
	h := &Handler{
		router:         chi.NewRouter(),
		userService:    us,
		roomService:    rs,
		slotService:    ss,
		bookingService: bs,
		logger:         logger,
		jwtSecret:      secretKey,
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
		r.Post("/login", h.login)
	})

	h.router.Group(func(r chi.Router) {
		r.Use(middleware.RequireAuth(h.logger, h.jwtSecret))

		r.Get("/rooms/list", h.getRooms)

		r.With(middleware.CheckRoomID(h.logger)).
			Get("/rooms/{roomId}/slots/list", h.getFreeSlots)

		r.With(middleware.RequireRole(domain.RoleAdmin)).
			Get("/bookings/list", h.listBookings)

		r.With(middleware.RequireRole(domain.RoleUser)).
			Get("/bookings/my", h.getMyBookings)

		r.With(middleware.RequireJSON(h.logger), middleware.RequireRole(domain.RoleAdmin)).
			Post("/rooms/create", h.сreateRoom)

		r.With(middleware.RequireJSON(h.logger), middleware.RequireRole(domain.RoleAdmin)).
			Post("/rooms/{roomId}/schedule/create", h.createSchedule)

		r.With(middleware.RequireJSON(h.logger), middleware.RequireRole(domain.RoleUser)).
			Post("/bookings/create", h.createBooking)

		r.With(middleware.RequireRole(domain.RoleUser)).
			Post("/bookings/{bookingId}/cancel", h.cancelBooking)
	})
}

func (h *Handler) Info(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

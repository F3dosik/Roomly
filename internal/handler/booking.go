package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/avito-internships/test-backend-1-F3dosik/internal/ctxkey"
	"github.com/avito-internships/test-backend-1-F3dosik/internal/domain"
	"github.com/avito-internships/test-backend-1-F3dosik/internal/httputil"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h *Handler) createBooking(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(ctxkey.UserIDKey).(uuid.UUID)

	var req createBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Debugw("createBooking: can't decode JSON body", "error", err)
		httputil.HandleError(w, httputil.NewAppError(
			httputil.ErrCodeInvalidRequest,
			"invalid request body",
			http.StatusBadRequest,
		))
		return
	}

	if req.SlotID == uuid.Nil {
		httputil.HandleError(w, httputil.NewAppError(
			httputil.ErrCodeInvalidRequest,
			"slotId is required",
			http.StatusBadRequest,
		))
		return
	}

	booking, err := h.bookingService.CreateBooking(r.Context(), userID, req.SlotID, req.CreateConferenceLink)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrSlotNotFound):
			h.logger.Debugw("createBooking: slot not found", "error", err)
			httputil.HandleError(w, httputil.NewAppError(
				httputil.ErrCodeSlotNotFound,
				"slot not found",
				http.StatusNotFound,
			))
		case errors.Is(err, domain.ErrSlotAlreadyBooked):
			h.logger.Debugw("createBooking: slot already booked", "error", err)
			httputil.HandleError(w, httputil.NewAppError(
				httputil.ErrCodeSlotAlreadyBooked,
				"slot is already booked",
				http.StatusConflict,
			))
		case errors.Is(err, domain.ErrBookingInPast):
			h.logger.Debugw("createBooking: slot in past", "error", err)
			httputil.HandleError(w, httputil.NewAppError(
				httputil.ErrCodeInvalidRequest,
				"cannot book slot in the past",
				http.StatusBadRequest,
			))
		default:
			h.logger.Errorw("createBooking: internal error", "error", err)
			httputil.HandleError(w, httputil.NewAppError(
				httputil.ErrCodeInternalError,
				http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError,
			))
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(bookingWrapResponse{Booking: toBookingResponse(booking)}); err != nil {
		h.logger.Errorw("createBooking: encode error", "error", err)
	}
}

func (h *Handler) listBookings(w http.ResponseWriter, r *http.Request) {
	page := 1
	pageSize := 20

	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		p, err := strconv.Atoi(pageStr)
		if err != nil || p < 1 {
			httputil.HandleError(w, httputil.NewAppError(
				httputil.ErrCodeInvalidRequest,
				"page must be a positive integer",
				http.StatusBadRequest,
			))
			return
		}
		page = p
	}

	if pageSizeStr := r.URL.Query().Get("pageSize"); pageSizeStr != "" {
		ps, err := strconv.Atoi(pageSizeStr)
		if err != nil || ps < 1 || ps > 100 {
			httputil.HandleError(w, httputil.NewAppError(
				httputil.ErrCodeInvalidRequest,
				"pageSize must be between 1 and 100",
				http.StatusBadRequest,
			))
			return
		}
		pageSize = ps
	}

	bookings, total, err := h.bookingService.ListBookings(r.Context(), page, pageSize)
	if err != nil {
		h.logger.Errorw("listBookings: internal error", "error", err)
		httputil.HandleError(w, httputil.NewAppError(
			httputil.ErrCodeInternalError,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		))
		return
	}

	resp := toListBookingsResponse(bookings, page, pageSize, total)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.logger.Errorw("listBookings: encode error", "error", err)
	}
}

func (h *Handler) getMyBookings(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(ctxkey.UserIDKey).(uuid.UUID)

	bookings, err := h.bookingService.GetMyBookings(r.Context(), userID)
	if err != nil {
		h.logger.Errorw("getMyBookings: internal error", "error", err)
		httputil.HandleError(w, httputil.NewAppError(
			httputil.ErrCodeInternalError,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		))
		return
	}

	resp := toMyBookingResponse(bookings)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.logger.Errorw("getMyBookings: encode error", "error", err)
	}
}

func (h *Handler) cancelBooking(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(ctxkey.UserIDKey).(uuid.UUID)
	bookingIDStr := chi.URLParam(r, "bookingId")
	bookingID, err := uuid.Parse(bookingIDStr)
	if err != nil {
		httputil.HandleError(w, httputil.NewAppError(
			httputil.ErrCodeInvalidRequest,
			"invalid roomId",
			http.StatusBadRequest,
		))
		return
	}

	booking, err := h.bookingService.CancelBooking(r.Context(), userID, bookingID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrBookingNotFound):
			h.logger.Debugw("cancelBooking: not found", "error", err)
			httputil.HandleError(w, httputil.NewAppError(
				httputil.ErrCodeBookingNotFound,
				"booking not found",
				http.StatusNotFound,
			))
		case errors.Is(err, domain.ErrForbidden):
			h.logger.Debugw("cancelBooking: forbidden", "error", err)
			httputil.HandleError(w, httputil.NewAppError(
				httputil.ErrCodeForbidden,
				"cannot cancel another user's booking",
				http.StatusForbidden,
			))
		default:
			h.logger.Errorw("cancelBooking: internal error", "error", err)
			httputil.HandleError(w, httputil.NewAppError(
				httputil.ErrCodeInternalError,
				http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError,
			))
		}
		return
	}

	resp := &bookingWrapResponse{Booking: toBookingResponse(booking)}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.logger.Errorw("cancelBooking: encode error", "error", err)
	}
}

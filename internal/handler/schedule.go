package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/avito-internships/test-backend-1-F3dosik/internal/domain"
	"github.com/avito-internships/test-backend-1-F3dosik/internal/httputil"
	"github.com/avito-internships/test-backend-1-F3dosik/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h *Handler) createSchedule(w http.ResponseWriter, r *http.Request) {
	roomIDStr := chi.URLParam(r, "roomId")
	roomID, err := uuid.Parse(roomIDStr)
	if err != nil {
		httputil.HandleError(w, httputil.NewAppError(
			httputil.ErrCodeInvalidRequest,
			"invalid roomId",
			http.StatusBadRequest,
		))
		return
	}

	var req scheduleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Debugw("createSchedule: can't decode JSON body", "error", err)
		httputil.HandleError(w, httputil.NewAppError(
			httputil.ErrCodeInvalidRequest,
			"invalid request body",
			http.StatusBadRequest,
		))
		return
	}

	schedule, err := h.roomService.CreateSchedule(
		r.Context(), roomID,
		req.DaysOfWeek, req.StartTime, req.EndTime,
	)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrEmptyDaysOfWeek),
			errors.Is(err, service.ErrInvalidDaysOfWeek),
			errors.Is(err, service.ErrInvalidTimeFormat),
			errors.Is(err, service.ErrInvalidTimeRange):
			h.logger.Debugw("createSchedule: invalid request", "error", err)
			httputil.HandleError(w, httputil.NewAppError(
				httputil.ErrCodeInvalidRequest,
				err.Error(),
				http.StatusBadRequest,
			))

		case errors.Is(err, domain.ErrRoomNotFound):
			h.logger.Debugw("createSchedule", "error", err)
			httputil.HandleError(w, httputil.NewAppError(
				httputil.ErrCodeRoomNotFound,
				"room not found",
				http.StatusNotFound,
			))

		case errors.Is(err, domain.ErrScheduleExists):
			h.logger.Debugw("createSchedule", "error", err)
			httputil.HandleError(w, httputil.NewAppError(
				httputil.ErrCodeScheduleExists,
				"schedule for this room already exists and cannot be changed",
				http.StatusConflict,
			))

		default:
			h.logger.Debugw("createSchedule: internal error", "error", err)
			httputil.HandleError(w, httputil.NewAppError(
				httputil.ErrCodeInternalError,
				http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError,
			))
		}
		return
	}

	resp := toGetScheduleResponse(schedule)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.logger.Errorw("createSchedule: encode error", "error", err)
		return
	}

}

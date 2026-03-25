package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/avito-internships/test-backend-1-F3dosik/internal/ctxkey"
	"github.com/avito-internships/test-backend-1-F3dosik/internal/domain"
	"github.com/avito-internships/test-backend-1-F3dosik/internal/httputil"
	"github.com/google/uuid"
)

func (h *Handler) getFreeSlots(w http.ResponseWriter, r *http.Request) {
	dateStr := r.URL.Query().Get("date")
	if dateStr == "" {
		h.logger.Debug("getFreeSlots: date param is empty or not exists")
		httputil.HandleError(w, httputil.NewAppError(
			httputil.ErrCodeInvalidRequest,
			"date param is empty or not exists",
			http.StatusBadRequest,
		))
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		h.logger.Debug("getFreeSlots: inavlid date format")
		httputil.HandleError(w, httputil.NewAppError(
			httputil.ErrCodeInvalidRequest,
			"inavlid date format",
			http.StatusBadRequest,
		))
		return
	}

	roomID := r.Context().Value(ctxkey.RoomIDKey).(uuid.UUID)

	slots, err := h.slotService.GetFreeSlots(r.Context(), roomID, date)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrRoomNotFound):
			h.logger.Debugw("getFreeSlots", "error", err)
			httputil.HandleError(w, httputil.NewAppError(
				httputil.ErrCodeRoomNotFound,
				"room not found",
				http.StatusNotFound,
			))
		default:
			h.logger.Debug("getFreeSlots: internal error", "error", err)
			httputil.HandleError(w, httputil.NewAppError(
				httputil.ErrCodeInternalError,
				"internal error",
				http.StatusInternalServerError,
			))
		}
		return
	}

	resp := toGetSlotResponse(slots)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.logger.Errorw("getFreeSlots: encode error", "error", err)
		return
	}

}

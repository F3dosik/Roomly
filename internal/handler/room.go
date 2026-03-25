package handler

import (
	"encoding/json"
	"net/http"

	"github.com/avito-internships/test-backend-1-F3dosik/internal/httputil"
)

func (h *Handler) getRooms(w http.ResponseWriter, r *http.Request) {
	rooms, err := h.roomService.GetRooms(r.Context())
	if err != nil {
		h.logger.Debugw("getRooms: internal error", "error", err)
		httputil.HandleError(w, httputil.NewAppError(
			httputil.ErrCodeInternalError,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		))
		return
	}

	resp := toGetRoomsResponse(rooms)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.logger.Errorw("getRooms: encode error", "error", err)
		return
	}

}

func (h *Handler) сreateRoom(w http.ResponseWriter, r *http.Request) {
	var req roomRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Debugw("createRoom: can't decode JSON body", "error", err)
		httputil.HandleError(w, httputil.NewAppError(
			httputil.ErrCodeInvalidRequest,
			"invalid request body",
			http.StatusBadRequest,
		))
		return
	}

	room := toRoom(req)
	if err := h.roomService.CreateRoom(r.Context(), room); err != nil {
		h.logger.Debugw("createRoom: internal error", "error", err)
		httputil.HandleError(w, httputil.NewAppError(
			httputil.ErrCodeInternalError,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		))
		return
	}

	resp := toGetRoomResponse(room)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.logger.Errorw("createRoom: encode error", "error", err)
		return
	}
}

package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/avito-internships/test-backend-1-F3dosik/internal/domain"
	"github.com/avito-internships/test-backend-1-F3dosik/internal/httputil"
	"github.com/avito-internships/test-backend-1-F3dosik/internal/service"
)

type Token struct {
	Token string `json:"token"`
}

type dummyLoginRequestBody struct {
	Role domain.Role `json:"role"`
}

func (h *Handler) dummyLogin(w http.ResponseWriter, r *http.Request) {
	var req dummyLoginRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Debugw("dummyLogin: can't decode JSON body", "error", err)
		httputil.HandleError(w, httputil.NewAppError(
			httputil.ErrCodeInvalidRequest,
			"invalid request body",
			http.StatusBadRequest,
		))
		return
	}

	token, err := h.userService.DummyLogin(r.Context(), req.Role)
	if err != nil {
		if errors.Is(err, service.ErrInvalidRole) {
			h.logger.Debugw("dummyLogin: Invalid request (invalid role value)", "role", req.Role)
			httputil.HandleError(w, httputil.NewAppError(
				httputil.ErrCodeInvalidRequest,
				http.StatusText(http.StatusBadRequest),
				http.StatusBadRequest,
			))
		} else {
			h.logger.Debugw("dummyLogin: can't generate token", "error", err)
			httputil.HandleError(w, httputil.NewAppError(
				httputil.ErrCodeInternalError,
				http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError,
			))
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(token); err != nil {
		h.logger.Errorw("dummyLogin: encode error", "error", err)
		return
	}
}

type registerRequestBody struct {
	Email    string      `json:"email"`
	Password string      `json:"password"`
	Role     domain.Role `json:"role"`
}

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	var req registerRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Debugw("register: can't decode JSON body", "error", err)
		httputil.HandleError(w, httputil.NewAppError(
			httputil.ErrCodeInvalidRequest,
			"invalid request body",
			http.StatusBadRequest,
		))
		return
	}

	user, err := h.userService.Register(r.Context(), req.Email, req.Password, req.Role)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidEmail),
			errors.Is(err, service.ErrEmailAlreadyExist),
			errors.Is(err, service.ErrPasswordTooShort),
			errors.Is(err, service.ErrInvalidRole):
			h.logger.Debugw("register: bad request", "error", err)
			httputil.HandleError(w, httputil.NewAppError(
				httputil.ErrCodeInvalidRequest,
				err.Error(),
				http.StatusBadRequest,
			))
		default:
			h.logger.Debugw("register: internal error", "error", err)
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

	if err := json.NewEncoder(w).Encode(toUserResponse(user)); err != nil {
		h.logger.Errorw("dummyLogin: encode error", "error", err)
		return
	}
}

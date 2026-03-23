package httputil

import (
	"encoding/json"
	"errors"
	"net/http"
)

type ErrCode string

const (
	ErrCodeInvalidRequest    ErrCode = "INVALID_REQUEST"
	ErrCodeUnauthorized      ErrCode = "UNAUTHORIZED"
	ErrCodeNotFound          ErrCode = "NOT_FOUND"
	ErrCodeRoomNotFound      ErrCode = "ROOM_NOT_FOUND"
	ErrCodeSlotNotFound      ErrCode = "SLOT_NOT_FOUND"
	ErrCodeSlotAlreadyBooked ErrCode = "SLOT_ALREADY_BOOKED"
	ErrCodeBookingNotFound   ErrCode = "BOOKING_NOT_FOUND"
	ErrCodeForbidden         ErrCode = "FORBIDDEN"
	ErrCodeScheduleExists    ErrCode = "SCHEDULE_EXISTS"
	ErrCodeInternalError     ErrCode = "INTERNAL_ERROR"
)

type ErrorResponse struct {
	Error ErrorBody `json:"error"`
}

type ErrorBody struct {
	Code    ErrCode `json:"code"`
	Message string  `json:"message"`
}

func WriteError(w http.ResponseWriter, code ErrCode, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error: ErrorBody{
			Code:    code,
			Message: message,
		},
	})
}

type AppError struct {
	Code    ErrCode
	Message string
	Status  int
}

func NewAppError(code ErrCode, message string, status int) *AppError {
	return &AppError{Code: code, Message: message, Status: status}
}

func (e *AppError) Error() string {
	return e.Message
}

func HandleError(w http.ResponseWriter, err error) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		WriteError(w, appErr.Code, appErr.Message, appErr.Status)
		return
	}
	WriteError(w, ErrCodeInternalError, "internal server error", http.StatusInternalServerError)
}

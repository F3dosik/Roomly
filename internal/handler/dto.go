package handler

import (
	"time"

	"github.com/avito-internships/test-backend-1-F3dosik/internal/domain"
	"github.com/google/uuid"
)

type Token struct {
	Token string `json:"token"`
}

type dummyLoginRequestBody struct {
	Role domain.Role `json:"role"`
}

type loginRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type registerRequestBody struct {
	Email    string      `json:"email"`
	Password string      `json:"password"`
	Role     domain.Role `json:"role"`
}

type userResponse struct {
	ID        uuid.UUID  `json:"id"`
	Email     string     `json:"email"`
	Role      string     `json:"role"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
}

func toUserResponse(u *domain.User) userResponse {
	return userResponse{
		ID:        u.ID,
		Email:     u.Email,
		Role:      string(u.Role),
		CreatedAt: &u.CreatedAt,
	}
}

type roomRequestBody struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Capacity    *int    `json:"capacity"`
}

type roomResponse struct {
	ID          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	Description *string    `json:"description"`
	Capacity    *int       `json:"capacity"`
	CreatedAt   *time.Time `json:"createdAt"`
}

func toRoomResponse(room *domain.Room) *roomResponse {
	return &roomResponse{
		ID:          room.ID,
		Name:        room.Name,
		Description: room.Description,
		Capacity:    room.Capacity,
		CreatedAt:   room.CreatedAt,
	}
}

type getRoomsResponse struct {
	Rooms []*roomResponse `json:"rooms"`
}

func toGetRoomsResponse(rooms []*domain.Room) *getRoomsResponse {
	resp := getRoomsResponse{Rooms: make([]*roomResponse, 0, len(rooms))}
	for _, room := range rooms {
		resp.Rooms = append(resp.Rooms, toRoomResponse(room))
	}
	return &resp
}

func toRoom(req roomRequestBody) *domain.Room {
	return &domain.Room{
		Name:        req.Name,
		Description: req.Description,
		Capacity:    req.Capacity,
	}
}

type scheduleRequest struct {
	DaysOfWeek []domain.DayOfWeek `json:"daysOfWeek"`
	StartTime  string             `json:"startTime"`
	EndTime    string             `json:"endTime"`
}

type scheduleResponse struct {
	ID         uuid.UUID          `json:"id"`
	RoomID     uuid.UUID          `json:"roomId"`
	DaysOfWeek []domain.DayOfWeek `json:"daysOfWeek"`
	StartTime  string             `json:"startTime"`
	EndTime    string             `json:"endTime"`
}

func toScheduleResponse(s *domain.Schedule) *scheduleResponse {
	return &scheduleResponse{
		ID:         s.ID,
		RoomID:     s.RoomID,
		DaysOfWeek: s.DaysOfWeek,
		StartTime:  s.StartTime.Format("15:04"),
		EndTime:    s.EndTime.Format("15:04"),
	}
}

type slotResponse struct {
	ID     uuid.UUID `json:"id"`
	RoomID uuid.UUID `json:"roomId"`
	Start  time.Time `json:"start"`
	End    time.Time `json:"end"`
}

type bookingResponse struct {
	ID             uuid.UUID  `json:"id"`
	SlotID         uuid.UUID  `json:"slotId"`
	UserID         uuid.UUID  `json:"userId"`
	Status         string     `json:"status"`
	ConferenceLink *string    `json:"confernceLink"`
	CreatedAt      *time.Time `json:"createdAt"`
}

type paginationResponse struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
	Total    int `json:"total"`
}

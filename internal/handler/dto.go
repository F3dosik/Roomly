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

type getRoomResponse struct {
	Room *roomResponse `json:"room"`
}

func toGetRoomResponse(room *domain.Room) *getRoomResponse {
	return &getRoomResponse{
		Room: toRoomResponse(room),
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

type getScheduleResponse struct {
	Schedule *scheduleResponse `json:"schedule"`
}

func toGetScheduleResponse(schedule *domain.Schedule) *getScheduleResponse {
	return &getScheduleResponse{Schedule: toScheduleResponse(schedule)}
}

type slotResponse struct {
	ID     uuid.UUID `json:"id"`
	RoomID uuid.UUID `json:"roomId"`
	Start  time.Time `json:"start"`
	End    time.Time `json:"end"`
}

func toSlotResponse(slot *domain.Slot) *slotResponse {
	return &slotResponse{
		ID:     slot.ID,
		RoomID: slot.RoomID,
		Start:  slot.StartsAt,
		End:    slot.EndsAt,
	}
}

type getSlotResponse struct {
	Slots []*slotResponse `json:"slots"`
}

func toGetSlotResponse(slots []*domain.Slot) *getSlotResponse {
	resp := getSlotResponse{Slots: make([]*slotResponse, 0, len(slots))}
	for _, slot := range slots {
		resp.Slots = append(resp.Slots, toSlotResponse(slot))
	}
	return &resp
}

type createBookingRequest struct {
	SlotID               uuid.UUID `json:"slotId"`
	CreateConferenceLink bool      `json:"createConferenceLink"`
}

type bookingResponse struct {
	ID             uuid.UUID            `json:"id"`
	SlotID         uuid.UUID            `json:"slotId"`
	UserID         uuid.UUID            `json:"userId"`
	Status         domain.BookingStatus `json:"status"`
	ConferenceLink *string              `json:"conferenceLink"`
	CreatedAt      *time.Time           `json:"createdAt"`
}

func toBookingResponse(b *domain.Booking) *bookingResponse {
	return &bookingResponse{
		ID:             b.ID,
		SlotID:         b.SlotID,
		UserID:         b.UserID,
		Status:         b.Status,
		ConferenceLink: b.ConferenceLink,
		CreatedAt:      &b.CreatedAt,
	}
}

type bookingWrapResponse struct {
	Booking *bookingResponse `json:"booking"`
}

type listBookingsResponse struct {
	Bookings   []*bookingResponse `json:"bookings"`
	Pagination paginationResponse `json:"pagination"`
}

func toListBookingsResponse(bookings []*domain.Booking, page, pageSize, total int) *listBookingsResponse {
	resp := &listBookingsResponse{
		Bookings: make([]*bookingResponse, 0, len(bookings)),
		Pagination: paginationResponse{
			Page:     page,
			PageSize: pageSize,
			Total:    total,
		},
	}

	for _, b := range bookings {
		resp.Bookings = append(resp.Bookings, toBookingResponse(b))
	}
	return resp
}

type paginationResponse struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
	Total    int `json:"total"`
}

type myBookingsResponse struct {
	Bookings []*bookingResponse `json:"bookings"`
}

func toMyBookingResponse(bookings []*domain.Booking) *myBookingsResponse {
	resp := myBookingsResponse{Bookings: make([]*bookingResponse, 0, len(bookings))}
	for _, b := range bookings {
		resp.Bookings = append(resp.Bookings, toBookingResponse(b))
	}
	return &resp
}

package ctxkey

type contextKey string

const (
	UserIDKey contextKey = "user_id"
	RoleKey   contextKey = "role"
	RoomIDKey contextKey = "room_id"
)

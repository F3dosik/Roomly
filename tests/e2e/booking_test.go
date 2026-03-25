//go:build e2e

package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const baseURL = "http://roomly:8080"

type tokenResponse struct {
	Token string `json:"token"`
}

type roomResponse struct {
	Room struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"room"`
}

type scheduleResponse struct {
	Schedule struct {
		ID     string `json:"id"`
		RoomID string `json:"roomId"`
	} `json:"schedule"`
}

type slotsResponse struct {
	Slots []struct {
		ID    string `json:"id"`
		Start string `json:"start"`
	} `json:"slots"`
}

type bookingResponse struct {
	Booking struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	} `json:"booking"`
}

type myBookingsResponse struct {
	Bookings []struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	} `json:"bookings"`
}

func post(t *testing.T, url, token string, body any) *http.Response {
	t.Helper()
	b, err := json.Marshal(body)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(b))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	return resp
}

func get(t *testing.T, url, token string) *http.Response {
	t.Helper()
	req, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	return resp
}

func getToken(t *testing.T, role string) string {
	t.Helper()
	resp := post(t, baseURL+"/dummyLogin", "", map[string]string{"role": role})
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var result tokenResponse
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
	require.NotEmpty(t, result.Token)
	return result.Token
}

func TestE2E_RoomScheduleBooking(t *testing.T) {
	adminToken := getToken(t, "admin")
	userToken := getToken(t, "user")

	// 1. создать комнату
	roomResp := post(t, baseURL+"/rooms/create", adminToken, map[string]any{
		"name":        fmt.Sprintf("E2E Room %d", time.Now().UnixNano()),
		"description": "E2E test room",
		"capacity":    4,
	})
	defer roomResp.Body.Close()
	require.Equal(t, http.StatusCreated, roomResp.StatusCode)

	var room roomResponse
	require.NoError(t, json.NewDecoder(roomResp.Body).Decode(&room))
	roomID := room.Room.ID
	require.NotEmpty(t, roomID)

	// 2. создать расписание
	schedResp := post(t, baseURL+"/rooms/"+roomID+"/schedule/create", adminToken, map[string]any{
		"daysOfWeek": []int{1, 2, 3, 4, 5, 6, 7},
		"startTime":  "09:00",
		"endTime":    "18:00",
	})
	defer schedResp.Body.Close()
	require.Equal(t, http.StatusCreated, schedResp.StatusCode)

	var schedule scheduleResponse
	require.NoError(t, json.NewDecoder(schedResp.Body).Decode(&schedule))
	require.NotEmpty(t, schedule.Schedule.ID)

	// 3. подождать воркер или вставить слоты — ждём до 5 секунд
	date := time.Now().UTC().AddDate(0, 0, 1).Format("2006-01-02")
	var slotID string
	require.Eventually(t, func() bool {
		resp := get(t, fmt.Sprintf("%s/rooms/%s/slots/list?date=%s", baseURL, roomID, date), userToken)
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return false
		}
		var slots slotsResponse
		if err := json.NewDecoder(resp.Body).Decode(&slots); err != nil {
			return false
		}
		if len(slots.Slots) == 0 {
			return false
		}
		slotID = slots.Slots[0].ID
		return true
	}, 10*time.Second, 500*time.Millisecond)

	require.NotEmpty(t, slotID)

	// 4. создать бронь
	bookResp := post(t, baseURL+"/bookings/create", userToken, map[string]any{
		"slotId": slotID,
	})
	defer bookResp.Body.Close()
	require.Equal(t, http.StatusCreated, bookResp.StatusCode)

	var booking bookingResponse
	require.NoError(t, json.NewDecoder(bookResp.Body).Decode(&booking))
	require.Equal(t, "active", booking.Booking.Status)
	bookingID := booking.Booking.ID
	require.NotEmpty(t, bookingID)

	// 5. проверить /bookings/my
	myResp := get(t, baseURL+"/bookings/my", userToken)
	defer myResp.Body.Close()
	require.Equal(t, http.StatusOK, myResp.StatusCode)

	var myBookings myBookingsResponse
	require.NoError(t, json.NewDecoder(myResp.Body).Decode(&myBookings))
	require.NotEmpty(t, myBookings.Bookings)

	found := false
	for _, b := range myBookings.Bookings {
		if b.ID == bookingID {
			found = true
			break
		}
	}
	require.True(t, found, "booking not found in /bookings/my")
}

func TestE2E_CancelBooking(t *testing.T) {
	adminToken := getToken(t, "admin")
	userToken := getToken(t, "user")

	// 1. создать комнату и расписание
	roomResp := post(t, baseURL+"/rooms/create", adminToken, map[string]any{
		"name":     fmt.Sprintf("E2E Cancel Room %d", time.Now().UnixNano()),
		"capacity": 2,
	})
	defer roomResp.Body.Close()
	require.Equal(t, http.StatusCreated, roomResp.StatusCode)

	var room roomResponse
	require.NoError(t, json.NewDecoder(roomResp.Body).Decode(&room))
	roomID := room.Room.ID
	schedResp := post(t, baseURL+"/rooms/"+roomID+"/schedule/create", adminToken, map[string]any{
		"daysOfWeek": []int{1, 2, 3, 4, 5, 6, 7},
		"startTime":  "09:00",
		"endTime":    "18:00",
	})
	defer schedResp.Body.Close()
	require.Equal(t, http.StatusCreated, schedResp.StatusCode)

	// 2. дождаться слотов
	date := time.Now().UTC().AddDate(0, 0, 1).Format("2006-01-02")
	var slotID string
	require.Eventually(t, func() bool {
		resp := get(t, fmt.Sprintf("%s/rooms/%s/slots/list?date=%s", baseURL, roomID, date), userToken)
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return false
		}
		var slots slotsResponse
		if err := json.NewDecoder(resp.Body).Decode(&slots); err != nil {
			return false
		}
		if len(slots.Slots) == 0 {
			return false
		}
		slotID = slots.Slots[0].ID
		return true
	}, 10*time.Second, 500*time.Millisecond)

	// 3. создать бронь
	bookResp := post(t, baseURL+"/bookings/create", userToken, map[string]any{
		"slotId": slotID,
	})
	defer bookResp.Body.Close()
	require.Equal(t, http.StatusCreated, bookResp.StatusCode)

	var booking bookingResponse
	require.NoError(t, json.NewDecoder(bookResp.Body).Decode(&booking))
	bookingID := booking.Booking.ID

	// 4. отменить бронь
	cancelResp := post(t, baseURL+"/bookings/"+bookingID+"/cancel", userToken, nil)
	defer cancelResp.Body.Close()
	require.Equal(t, http.StatusOK, cancelResp.StatusCode)

	var cancelled bookingResponse
	require.NoError(t, json.NewDecoder(cancelResp.Body).Decode(&cancelled))
	require.Equal(t, "cancelled", cancelled.Booking.Status)

	// 5. идемпотентность — повторная отмена
	cancelResp2 := post(t, baseURL+"/bookings/"+bookingID+"/cancel", userToken, nil)
	defer cancelResp2.Body.Close()
	require.Equal(t, http.StatusOK, cancelResp2.StatusCode)

	var cancelled2 bookingResponse
	require.NoError(t, json.NewDecoder(cancelResp2.Body).Decode(&cancelled2))
	require.Equal(t, "cancelled", cancelled2.Booking.Status)

	// 6. проверить что бронь не в /bookings/my (она в прошлом не будет, но статус cancelled)
	myResp := get(t, baseURL+"/bookings/my", userToken)
	defer myResp.Body.Close()
	require.Equal(t, http.StatusOK, myResp.StatusCode)
}

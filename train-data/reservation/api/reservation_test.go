package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AndreasKl/train-reservation-kata/train-data/reservation"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var reservationPayload = `{"train_id": "express_2000", "seats": ["1A", "1B"], "booking_reference": "75bcd15"}`
var reservationRequest = func(payload string) *http.Request {
	return httptest.NewRequest(http.MethodPost, "/reserve", bytes.NewBufferString(payload))
}

func TestReservationApi_ReserveSeats(t *testing.T) {
	api := NewReservationApi(reservation.NewTrainServiceWithDefaultTrains())

	recorder := httptest.NewRecorder()
	api.ReserveSeats(recorder, reservationRequest(reservationPayload), nil)

	assert.Equal(t, 200, recorder.Code)

	responseJson := map[string]map[string]any{}
	err := json.Unmarshal(recorder.Body.Bytes(), &responseJson)
	require.NoError(t, err)

	assert.Equal(t, map[string]any{"booking_reference": "75bcd15", "coach": "A", "seat_number": "1"}, responseJson["seats"]["1A"])
	assert.Equal(t, map[string]any{"booking_reference": "75bcd15", "coach": "B", "seat_number": "1"}, responseJson["seats"]["1B"])
	assert.Equal(t, map[string]any{"booking_reference": "", "coach": "A", "seat_number": "2"}, responseJson["seats"]["2A"])
}

func TestReservationApi_ReserveSeats_AlreadyBooked(t *testing.T) {
	api := NewReservationApi(reservation.NewTrainServiceWithDefaultTrains())

	recorder := httptest.NewRecorder()
	api.ReserveSeats(recorder, reservationRequest(reservationPayload), nil)
	assert.Equal(t, http.StatusOK, recorder.Code)

	recorder = httptest.NewRecorder()
	api.ReserveSeats(recorder, reservationRequest(reservationPayload), nil)
	assert.Equal(t, http.StatusConflict, recorder.Code)

	assert.Equal(t, "", recorder.Body.String())
}

func TestReservationApi_ReserveSeats_UnknownTrain(t *testing.T) {
	api := NewReservationApi(reservation.NewTrainServiceWithDefaultTrains())

	recorder := httptest.NewRecorder()
	api.ReserveSeats(recorder, reservationRequest(`{"train_id": "not_known", "seats": ["1A", "1B"], "booking_reference": "75bcd15"}`), nil)
	assert.Equal(t, http.StatusNotFound, recorder.Code)
	assert.Equal(t, "", recorder.Body.String())
}

func TestReservationApi_ResetAllReservations(t *testing.T) {
	api := NewReservationApi(reservation.NewTrainServiceWithDefaultTrains())

	recorder := httptest.NewRecorder()
	api.ResetAllReservations(recorder, httptest.NewRequest(http.MethodPost, "/reset", nil), nil)

	assert.Equal(t, 200, recorder.Code)
}

func TestReservationApi_FetchDataForTrainById(t *testing.T) {
	api := NewReservationApi(reservation.NewTrainServiceWithDefaultTrains())

	recorder := httptest.NewRecorder()
	api.FetchDataForTrainById(recorder, httptest.NewRequest(http.MethodGet, "/train", nil), httprouter.Params{httprouter.Param{Key: "trainID", Value: "express_2000"}})

	assert.Equal(t, 200, recorder.Code)
}

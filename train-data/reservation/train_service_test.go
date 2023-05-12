package reservation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var trains = Trains{
	"train-a": Train{
		Seats: Seats{
			"1A": Seat{BookingReference: ""},
			"1B": Seat{BookingReference: ""},
		},
	},
	"train-b": Train{
		Seats: Seats{
			"1A": Seat{BookingReference: ""},
			"1B": Seat{BookingReference: ""},
		},
	},
}

func TestNewTrainServiceWithDefaultTrains_TwoDefaultTrainsAreLoaded(t *testing.T) {
	trainService := NewTrainServiceWithDefaultTrains()
	assert.Len(t, trainService.Trains(), 2)
}

func TestReserveSeats_AllOrNothing(t *testing.T) {
	trainService := NewTrainService(trains)

	err := trainService.ReserveSeats(
		Reservation{TrainID: "train-a", Seats: []string{"1A", "1A"}, BookingReference: "1234"},
	)
	require.ErrorIs(t, err, ErrSeatAlreadyReserved)

	assertEmptySeat(t, trainService, "train-a", "1A")
	assertEmptySeat(t, trainService, "train-a", "1B")
}

func TestReserveSeats_UnknownTrain(t *testing.T) {
	trainService := NewTrainService(trains)

	err := trainService.ReserveSeats(
		Reservation{TrainID: "unknown-train", Seats: []string{"1A", "1A"}, BookingReference: "1234"},
	)
	require.ErrorIs(t, err, ErrTrainNotFound)
}

func TestResetAllReservations(t *testing.T) {
	trainService := NewTrainService(trains)

	err := trainService.ReserveSeats(
		Reservation{TrainID: "train-a", Seats: []string{"1A", "1B"}, BookingReference: "1234"},
	)
	require.NoError(t, err)
	assert.Equal(t, "1234", trainService.trains["train-a"].Seats["1A"].BookingReference)
	assert.Equal(t, "1234", trainService.trains["train-a"].Seats["1B"].BookingReference)

	trainService.ResetAllReservations()

	assertEmptySeat(t, trainService, "train-a", "1A")
	assertEmptySeat(t, trainService, "train-a", "1B")
}

func TestTrainByID(t *testing.T) {
	trainService := NewTrainService(trains)

	train, ok := trainService.TrainByID("train-a")
	assert.True(t, ok)

	assert.Equal(t, Seats{"1A": {Coach: "", SeatNumber: "", BookingReference: ""}, "1B": {Coach: "", SeatNumber: "", BookingReference: ""}}, train.Seats)
}

func TestTrainByID_UnknownTrain(t *testing.T) {
	trainService := NewTrainService(trains)

	_, ok := trainService.TrainByID("unknown-train")
	assert.False(t, ok)
}

func assertEmptySeat(t *testing.T, trainService *trainService, trainID string, seatID string) {
	t.Helper()

	assert.Equal(t, "", trainService.trains[trainID].Seats[seatID].BookingReference)
}

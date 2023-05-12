package reservation

import (
	_ "embed"
	"encoding/json"
	"log"
)

//go:embed trains.json
var trainData []byte

func loadTrains() (trains Trains, err error) {
	return trains, json.Unmarshal(trainData, &trains)
}

type trainService struct {
	trains Trains
}

func NewTrainServiceWithDefaultTrains() *trainService {
	trains, err := loadTrains()
	if err != nil {
		panic(err)
	}
	log.Printf("Prepared a few trains. Choooo, chooo:\n%+v\n\n", trains)
	return NewTrainService(trains)
}

func NewTrainService(trains Trains) *trainService {
	return &trainService{trains: trains}
}

func (t *trainService) Trains() Trains {
	return t.trains
}

func (t *trainService) TrainByID(trainID string) (Train, bool) {
	train, ok := t.trains[trainID]
	return train, ok
}

func (t *trainService) ResetAllReservations() {
	for trainName, train := range t.trains {
		for seatID, seat := range train.Seats {
			seat.BookingReference = ""
			t.trains[trainName].Seats[seatID] = seat
		}
	}
}

func (t *trainService) ReserveSeats(res Reservation) error {
	train, ok := t.TrainByID(res.TrainID)
	if !ok {
		return ErrTrainNotFound
	}

	scratchTrain := train.copy()
	for _, seatID := range res.Seats {
		seat := scratchTrain.Seats[seatID]
		if seat.BookingReference != "" {
			return ErrSeatAlreadyReserved
		}

		seat.BookingReference = res.BookingReference
		scratchTrain.Seats[seatID] = seat
	}
	t.trains[res.TrainID] = scratchTrain
	return nil
}

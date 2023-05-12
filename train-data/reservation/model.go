package reservation

import "errors"

var ErrTrainNotFound = errors.New("train not found")
var ErrSeatAlreadyReserved = errors.New("seat already reserved")
var ErrInvalidReservation = errors.New("invalid reservation")

type Trains map[string]Train

type Train struct {
	Seats Seats `json:"seats"`
}

func (t *Train) copy() Train {
	clonedSeats := Seats{}
	for seatID, seat := range t.Seats {
		clonedSeats[seatID] = seat
	}
	return Train{Seats: clonedSeats}
}

type Seats map[string]Seat

type Seat struct {
	Coach            string `json:"coach"`
	SeatNumber       string `json:"seat_number"`
	BookingReference string `json:"booking_reference"`
}

type Reservation struct {
	TrainID          string   `json:"train_id"`
	BookingReference string   `json:"booking_reference"`
	Seats            []string `json:"seats"`
}

func (r Reservation) Validate() error {
	if r.TrainID == "" || r.BookingReference == "" || len(r.Seats) == 0 {
		return ErrInvalidReservation
	}
	return nil
}

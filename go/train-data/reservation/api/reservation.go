package api

import (
	_ "embed"
	"encoding/json"
	"net/http"

	"github.com/AndreasKl/train-reservation-kata/train-data/reservation"
	"github.com/julienschmidt/httprouter"
)

type TrainsSupplier interface {
	Trains() reservation.Trains
	TrainByID(string) (reservation.Train, bool)
}

type TrainSeatsReserver interface {
	ReserveSeats(reservation.Reservation) error
}

type TrainReservationsResetter interface {
	ResetAllReservations()
}

type TrainService interface {
	TrainsSupplier
	TrainSeatsReserver
	TrainReservationsResetter
}

func NewReservationApi(trainService TrainService) *ReservationApi {
	return &ReservationApi{trainService: trainService}
}

type ReservationApi struct {
	trainService TrainService
}

func (a *ReservationApi) FetchDataForTrainById(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	trainID := ps.ByName("trainID")
	if train, ok := a.trainService.TrainByID(trainID); ok {
		_ = json.NewEncoder(w).Encode(train)
		return
	}
	w.WriteHeader(http.StatusNotFound)
}

func (a *ReservationApi) ReserveSeats(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var res reservation.Reservation
	if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := res.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := a.trainService.ReserveSeats(res); err != nil {
		http.Error(w, err.Error(), statusCodeForError(err))
		return
	}

	train, ok := a.trainService.TrainByID(res.TrainID)
	if !ok {
		http.Error(w, "Not able to find the train that received a reservation.", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(train); err != nil {
		http.Error(w, err.Error(), statusCodeForError(err))
	}
}

func (a *ReservationApi) ResetAllReservations(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	a.trainService.ResetAllReservations()
	_ = json.NewEncoder(w).Encode(a.trainService.Trains())
}

func statusCodeForError(err error) int {
	switch err {
	case reservation.ErrTrainNotFound:
		return http.StatusNotFound
	case reservation.ErrSeatAlreadyReserved:
		return http.StatusConflict
	}
	return http.StatusInternalServerError
}

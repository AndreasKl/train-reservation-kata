package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/julienschmidt/httprouter"
)

const (
	serverReadHeaderTimeout = 5 * time.Second
	serverWriteTimeout      = 10 * time.Second
	serverReadTimeOut       = 120 * time.Second
)

//go:embed trains.json
var trainData []byte

type application struct {
	*http.Server
}

func newApplication() *application {
	trains, err := prepareTrains()
	if err != nil {
		log.Fatalf("Failed to prepare trains. Cause: %s\n", err.Error())
	}

	log.Printf("Loaded a few trains. Choooo, chooo:\n\n%+v", trains)

	router := httprouter.New()
	router.GET("/data_for_train/:trainID", fetchDataForTrainById(trains))
	router.POST("/reserve", reserveSingleSeat(trains))
	router.POST("/reset", resetAllReservations(trains))

	return &application{
		&http.Server{
			Addr:              ":8080",
			ReadTimeout:       serverReadTimeOut,
			ReadHeaderTimeout: serverReadHeaderTimeout,
			WriteTimeout:      serverWriteTimeout,
			Handler:           router,
		}}
}

func (a *application) start() {
	log.Println("Application starting.")

	a.startHTTPServer()
	log.Println("Application started.")
}

func (a *application) startHTTPServer() {
	go func() {
		if err := a.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				log.Println("Server shutdown.")
			} else {
				log.Fatalf("Unexpected server error on server shutdown. Cause: %s\n", err.Error())
			}
		}
	}()
}

func (a *application) stop() error {
	log.Print("Application shutting down.")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := a.Shutdown(ctx)

	log.Print("Application shutdown.")
	return err
}

func main() {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	_ = os.Setenv("TZ", "UTC")

	app := newApplication()
	app.start()
	<-shutdown
	_ = app.stop()
}

type Trains map[string]Train

type Train struct {
	Seats Seats
}

type Seats map[string]Seat

type Seat struct {
	Coach            string `json:"coach"`
	SeatNumber       string `json:"seat_number"`
	BookingReference string `json:"booking_reference"`
}

func prepareTrains() (trains Trains, err error) {
	return trains, json.Unmarshal(trainData, &trains)
}

func fetchDataForTrainById(trains Trains) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		trainID := ps.ByName("trainID")
		if train, ok := trains[trainID]; ok {
			_ = json.NewEncoder(w).Encode(train)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}
}

type Reservation struct {
	TrainID          string   `json:"train_id"`
	BookingReference string   `json:"booking_reference"`
	Seats            []string `json:"seats"`
}

func reserveSingleSeat(trains Trains) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var reservations Reservation
		if err := json.NewDecoder(r.Body).Decode(&reservations); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if reservations.TrainID == "" || reservations.BookingReference == "" || len(reservations.Seats) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if train, ok := trains[reservations.TrainID]; ok {
			// FIXME: Build a train service to handle this logic and
			// consider all or nothing semantics or partial booking.

			for _, seatID := range reservations.Seats {
				seat := train.Seats[seatID]
				if seat.BookingReference != "" {
					w.WriteHeader(http.StatusConflict)
					return
				}

				trains[reservations.TrainID].Seats[seatID] = Seat{
					Coach:            seat.Coach,
					SeatNumber:       seat.SeatNumber,
					BookingReference: reservations.BookingReference,
				}
			}
			w.WriteHeader(http.StatusNoContent)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}
}

func resetAllReservations(trains Trains) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		for trainName, train := range trains {
			for seatID, seat := range train.Seats {
				trains[trainName].Seats[seatID] = Seat{
					Coach:            seat.Coach,
					SeatNumber:       seat.SeatNumber,
					BookingReference: "",
				}
			}
		}

		_ = json.NewEncoder(w).Encode(trains)
	}
}

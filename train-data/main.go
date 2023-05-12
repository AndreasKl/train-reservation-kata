package main

import (
	"context"
	_ "embed"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/AndreasKl/train-reservation-kata/train-data/reservation"
	"github.com/AndreasKl/train-reservation-kata/train-data/reservation/api"
	"github.com/julienschmidt/httprouter"
)

type application struct {
	*http.Server
}

func newApplication(production bool) *application {
	trainService := reservation.NewTrainServiceWithDefaultTrains()
	reservationApi := api.NewReservationApi(trainService)

	router := httprouter.New()
	router.GET("/data_for_train/:trainID", reservationApi.FetchDataForTrainById)
	router.POST("/reserve", reservationApi.ReserveSeats)
	router.POST("/reset", reservationApi.ResetAllReservations)

	return &application{
		&http.Server{
			Addr:    getPort(production),
			Handler: router,
		}}
}

func (a *application) start() {
	log.Println("Application starting.")

	a.startHTTPServer()
	log.Printf("Application started. Listening on port %s.\n", a.Addr)
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

	app := newApplication(true)
	app.start()
	<-shutdown
	_ = app.stop()
}

func getPort(production bool) string {
	if production {
		return ":8080"
	}

	port, err := findRandomFreePort()
	if err != nil {
		log.Fatalf("Could not find a free port. Cause: %s\n", err.Error())
	}
	return ":" + strconv.Itoa(port)
}

func findRandomFreePort() (port int, err error) {
	var a *net.TCPAddr
	if a, err = net.ResolveTCPAddr("tcp", "localhost:0"); err == nil {
		var listener *net.TCPListener
		if listener, err = net.ListenTCP("tcp", a); err == nil {
			defer listener.Close()
			return listener.Addr().(*net.TCPAddr).Port, nil
		}
	}
	return 0, errors.New("could not find free port")
}

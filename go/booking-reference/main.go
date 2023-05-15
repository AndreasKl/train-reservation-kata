package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/AndreasKl/train-reservation-kata/booking-reference/reference"
)

var defaultStartingPoint int64 = 12345678

type application struct {
	*http.Server
}

func newApplication(production bool) *application {
	referenceController := reference.NewController(configureStartingPoint())
	return &application{
		&http.Server{
			Addr:    getPort(production),
			Handler: http.HandlerFunc(referenceController.GenerateNext),
		},
	}
}

func configureStartingPoint() int64 {
	startingPoint, err := strconv.ParseInt(os.Getenv("STARTING_POINT"), 10, 0)
	if err != nil {
		log.Printf("Environment variable STARTING_POINT not set or invalid, defaulting to '%d'.\n", defaultStartingPoint)
		return defaultStartingPoint
	}
	return startingPoint
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

func getPort(production bool) string {
	if production {
		return ":8082"
	}

	port, err := findRandomFreePort()
	if err != nil {
		log.Fatalf("Could not find a free port. Cause: %s\n", err.Error())
	}
	return ":" + strconv.Itoa(port)
}

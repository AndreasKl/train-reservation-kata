package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/google/uuid"
)

type application struct {
	*http.Server
}

var defaultStartingPoint int64 = 12345678

func newApplication(production bool) *application {
	controller := NewController(configureStartingPoint())
	return &application{
		&http.Server{
			Addr:    getPort(production),
			Handler: controller,
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

type ID uuid.UUID

type Reference struct {
	ID ID
}

type Controller struct {
	startingPoint atomic.Int64
}

func NewController(startingPoint int64) *Controller {
	c := &Controller{}
	c.startingPoint.Store(startingPoint)
	return c
}

type ReferenceResponse struct {
	Value string `json:"value"`
}

func (c *Controller) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	reference, err := json.Marshal(c.fetchReference())
	if err != nil {
		log.Println("Not able convert to json.")
		resp.WriteHeader(500)
		return
	}

	resp.Header().Add("Content-Type", "application/json")
	_, err = resp.Write(reference)

	if err != nil {
		log.Println("Not able to send response.")
		return
	}
}

func (c *Controller) fetchReference() ReferenceResponse {
	reference := fmt.Sprintf("%016x", c.startingPoint.Add(1))
	return ReferenceResponse{Value: reference}
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

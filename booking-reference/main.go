package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/google/uuid"
)

const (
	serverReadHeaderTimeout = 5 * time.Second
	serverWriteTimeout      = 10 * time.Second
	serverReadTimeOut       = 120 * time.Second
)

type application struct {
	server *http.Server
}

var defaultStartingPoint int64 = 12345678

func newApplication() *application {
	controller := NewController(configureStartingPoint())
	return &application{
		server: &http.Server{
			Addr:              ":8080",
			ReadTimeout:       serverReadTimeOut,
			ReadHeaderTimeout: serverReadHeaderTimeout,
			WriteTimeout:      serverWriteTimeout,
			Handler:           controller,
		}}
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
	log.Println("Application started.")
}

func (a *application) startHTTPServer() {
	go func() {
		if err := a.server.ListenAndServe(); err != nil {
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
	err := a.server.Shutdown(ctx)

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

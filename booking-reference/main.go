package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

const (
	serverReadHeaderTimeout = 5 * time.Second
	serverWriteTimeout      = 10 * time.Second
	serverReadTimeOut       = 120 * time.Second
)

type application struct {
	server *http.Server
}

type MyHandler string

func (h MyHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	_, err := resp.Write([]byte(h))
	if err != nil {
		log.Warn().Err(err).Msg("Not able to send response.")
		return
	}
}

func newApplication() *application {
	return &application{
		server: &http.Server{
			Addr:              ":8080",
			ReadTimeout:       serverReadTimeOut,
			ReadHeaderTimeout: serverReadHeaderTimeout,
			WriteTimeout:      serverWriteTimeout,
			Handler:           MyHandler("Hello"),
		}}
}

func (a *application) start() {
	log.Logger.Info().Msg("Application starting.")

	a.startHTTPServer()

	log.Logger.Info().Msg("Application started.")
}

func (a *application) startHTTPServer() {
	go func() {
		if err := a.server.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				log.Logger.Info().Msg("Server shutdown.")
			} else {
				log.Logger.Panic().Err(err).Msg("Unexpected server error.")
			}
		}
	}()
}

func (a *application) stop() error {
	log.Logger.Info().Msg("Application shutting down.")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := a.server.Shutdown(ctx)

	log.Logger.Info().Msg("Application shutdown.")
	return err
}

func main() {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	_ = os.Setenv("TZ", "UTC")
	setupLogging()

	app := newApplication()
	app.start()
	<-shutdown
	_ = app.stop()
}

func setupLogging() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	log.Logger = log.With().Caller().Str("application", "booking-reference").Logger()
}

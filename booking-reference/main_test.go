package main

import (
	"fmt"
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/require"
)

func TestApplicationStart_ExpectRunningWebServer(t *testing.T) {
	port, err := getFreePort()
	require.NoError(t, err)

	app := newApplication()
	app.server.Addr = fmt.Sprintf(":%s", port)
	app.start()

	require.Eventually(t,
		func() bool {
			_, err = resty.New().R().Get(fmt.Sprintf("http://127.0.0.1%s", app.server.Addr))
			return err == nil
		},
		1*time.Second,
		100*time.Millisecond,
	)
	require.NoError(t, app.stop())
}

func getFreePort() (string, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return "", err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return "", err
	}
	defer func() { _ = l.Close() }()
	return strconv.Itoa(l.Addr().(*net.TCPAddr).Port), nil
}

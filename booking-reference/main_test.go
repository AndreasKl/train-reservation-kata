package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/require"
)

func TestApplicationStart_ExpectRunningWebServer(t *testing.T) {
	app := newApplication(false)
	app.start()

	require.Eventually(t,
		func() bool {
			_, err := resty.New().R().Get(fmt.Sprintf("http://127.0.0.1%s", app.Addr))
			return err == nil
		},
		2*time.Second,
		5*time.Millisecond,
	)
	require.NoError(t, app.stop())
}

package api

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestController_ServeHTTP(t *testing.T) {
	controller := NewController(213123)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/", nil)

	controller.ServeHTTP(recorder, request)

	assert.Equal(t, 200, recorder.Code)
	assert.Equal(t, "0000000000034084", recorder.Body.String())

	recorderForSecondRequest := httptest.NewRecorder()
	controller.ServeHTTP(recorderForSecondRequest, request)

	assert.Equal(t, 200, recorderForSecondRequest.Code)
	assert.Equal(t, "0000000000034085", recorderForSecondRequest.Body.String())
}

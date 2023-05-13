package reference

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestController_GenerateNext(t *testing.T) {
	controller := NewController(213123)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/", nil)

	controller.GenerateNext(recorder, request)

	assert.Equal(t, 200, recorder.Code)
	assert.JSONEq(t, `{"value": "0000000000034084"}`, recorder.Body.String())

	recorderForSecondRequest := httptest.NewRecorder()
	controller.GenerateNext(recorderForSecondRequest, request)

	assert.Equal(t, 200, recorderForSecondRequest.Code)
	assert.JSONEq(t, `{"value": "0000000000034085"}`, recorderForSecondRequest.Body.String())
}

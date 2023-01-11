package api

import (
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

type anotherType struct {
	ID        int `json:"id,omitempty"`
	Name      string
	CreatedAt *time.Time `json:"createdAt,omitempty"`
}

type referenceResponse struct {
	ID          string      `json:"id"`
	AnotherType anotherType `json:"anotherType"`
}

type Controller struct {
}

func NewController() *Controller {
	return &Controller{}
}

func (c Controller) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	_, err := resp.Write([]byte("reference"))
	if err != nil {
		log.Warn().Err(err).Msg("Not able to send response.")
		return
	}
}

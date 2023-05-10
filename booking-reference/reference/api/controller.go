package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/rs/zerolog/log"
)

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
	marshal, err := json.Marshal(c.fetchReference())
	if err != nil {
		log.Warn().Err(err).Msg("Not able convert to json.")
		resp.WriteHeader(500)
		return
	}

	resp.Header().Add("Content-Type", "application/json")
	_, err = resp.Write(marshal)

	if err != nil {
		log.Warn().Err(err).Msg("Not able to send response.")
		return
	}
}

func (c *Controller) fetchReference() ReferenceResponse {
	reference := fmt.Sprintf("%016x", c.startingPoint.Add(1))
	return ReferenceResponse{Value: reference}
}

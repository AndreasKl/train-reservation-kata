package reference

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type ReferenceResponse struct {
	Value string `json:"value"`
}

type Controller struct {
	startingPoint atomic.Int64
}

func NewController(startingPoint int64) *Controller {
	c := &Controller{}
	c.startingPoint.Store(startingPoint)
	return c
}

func (c *Controller) GenerateNext(w http.ResponseWriter, r *http.Request) {
	if err := json.NewEncoder(w).Encode(c.fetchReference()); err != nil {
		log.Printf("Not able to convert the response to json. Error: %s\n", err.Error())
		http.Error(w, "Not able convert to response to json.", 500)
		return
	}
}

func (c *Controller) fetchReference() ReferenceResponse {
	reference := fmt.Sprintf("%016x", c.startingPoint.Add(1))
	return ReferenceResponse{Value: reference}
}

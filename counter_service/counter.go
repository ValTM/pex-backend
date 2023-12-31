package counter_service

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/render"

	"pex-backend/utils"
)

type counterPayload struct {
	Counter string
	Uuid    string
}

// Bind checks for bad payload data
func (c counterPayload) Bind(*http.Request) error {
	if c.Uuid == "" {
		return fmt.Errorf("uuid is mandatory")
	}
	if c.Counter == "" {
		return fmt.Errorf("counter is mandatory")
	}
	return nil
}

func (c *CounterService) incrementHandler(w http.ResponseWriter, r *http.Request) {
	payload := &counterPayload{}
	if err := render.Bind(r, payload); err != nil {
		_ = render.Render(w, r, utils.RenderError(err, http.StatusBadRequest))
		return
	}
	currentClientMap, ok := c.countersMap[payload.Uuid]
	if ok {
		// increment existing counter if it exists
		currentClientMap[payload.Counter] = currentClientMap[payload.Counter] + 1
	} else {
		// initialize a new counter to 1
		newCounterMap := make(counter)
		newCounterMap[payload.Counter] = 1
		c.countersMap[payload.Uuid] = newCounterMap
	}
	_, err := w.Write([]byte(strconv.Itoa(c.countersMap[payload.Uuid][payload.Counter])))
	if err != nil {
		log.Printf("error writing response: %d", err)
	}
}

func (c *CounterService) decrementHandler(w http.ResponseWriter, r *http.Request) {
	payload := &counterPayload{}
	if err := render.Bind(r, payload); err != nil {
		_ = render.Render(w, r, utils.RenderError(err, http.StatusBadRequest))
		return
	}
	currentClientMap, ok := c.countersMap[payload.Uuid]
	if ok {
		if currentClientMap[payload.Counter]-1 < 0 {
			_ = render.Render(w, r, utils.RenderError(errors.New("counter value can not be less than zero"), http.StatusBadRequest))
			return
		}
		currentClientMap[payload.Counter] = currentClientMap[payload.Counter] - 1
	} else {
		// the user is trying to decrement a non-existing counter
		// the default value is 0, so it would go negative, which we don't allow
		_ = render.Render(w, r, utils.RenderError(errors.New("counter value can not be less than zero"), http.StatusBadRequest))
		return
	}
	_, err := w.Write([]byte(strconv.Itoa(c.countersMap[payload.Uuid][payload.Counter])))
	if err != nil {
		log.Printf("error writing response: %d", err)
	}
}

func (c *CounterService) resetHandler(w http.ResponseWriter, r *http.Request) {
	payload := &counterPayload{}
	if err := render.Bind(r, payload); err != nil {
		_ = render.Render(w, r, utils.RenderError(err, http.StatusBadRequest))
		return
	}
	currentClientMap, ok := c.countersMap[payload.Uuid]
	if ok {
		currentClientMap[payload.Counter] = 0
	}
	// if not found, it's okay to not do anything and lie to the client with a 0, since the counter does not exist,
	// and on increment it will be created, or an error will be thrown when decrementing (since it goes to -1)
	_, err := w.Write([]byte(strconv.Itoa(0)))
	if err != nil {
		log.Printf("error writing response: %d", err)
	}
}

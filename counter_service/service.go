package counter_service

import (
	"github.com/go-chi/chi/v5"
)

type counter map[string]int
type CounterService struct {
	countersMap map[string]counter
}

func (c *CounterService) InitService(r *chi.Mux) {
	c.countersMap = make(map[string]counter)
	c.addRoutes(r)
}

func (c *CounterService) addRoutes(r chi.Router) {
	r.Route("/counter", func(rr chi.Router) {
		rr.Post("/inc", c.incrementHandler)
		rr.Post("/dec", c.decrementHandler)
		rr.Post("/reset", c.resetHandler)
	})
}

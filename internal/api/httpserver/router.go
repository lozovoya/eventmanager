package httpserver

import (
	v1 "EventManager/internal/api/v1"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
)

func NewRouter (mux *chi.Mux, eventController *v1.Event) chi.Mux {
	mux.Use(middleware.Logger)
	mux.Route("/api/v1", func(router chi.Router) {
		router.Post("/events", eventController.InEvent)
		router.Get("/snapshot", eventController.GetSnapShot)
	})

	log.Println("NewRouter: new router is activated")
	return *mux
}

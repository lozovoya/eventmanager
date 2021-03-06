package main

import (
	"EventManager/internal/api/httpserver"
	v1 "EventManager/internal/api/v1"
	"EventManager/internal/bus"
	"EventManager/internal/cache"
	"fmt"
	"github.com/go-chi/chi/v5"
	"log"
	"net"
	"net/http"
	"os"
)

const (
	defaultPort     = "9999"
	defaultHost     = "0.0.0.0"
	defaultCacheDSN = "redis://eventcache:6379/0"
	//defaultCacheDSN = "redis://localhost:6379/0"
	defaultBusDSN = "amqp://guest:guest@eventbus:5672/"
)

func main() {
	port, ok := os.LookupEnv("EV_PORT")
	if !ok {
		port = defaultPort
	}

	host, ok := os.LookupEnv("EV_HOST")
	if !ok {
		host = defaultHost
	}

	cacheDSN, ok := os.LookupEnv("EV_CACHE")
	if !ok {
		cacheDSN = defaultCacheDSN
	}

	busDSN, ok := os.LookupEnv("EV_BUS")
	if !ok {
		busDSN = defaultBusDSN

	}

	if err := execute(net.JoinHostPort(host, port), cacheDSN, busDSN); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func execute(addr, cacheDSN, busDSN string) (err error) {
	cacheCallPool := cache.InitCache(cacheDSN)
	cacheCall := cache.NewCallCache(cacheCallPool)

	busConn, err := bus.InitBus(busDSN)
	if err != nil {
		return fmt.Errorf("Execute: %w", err)
	}
	defer busConn.Close()
	busEvent := bus.NewEventBus(busConn)

	eventsController := v1.NewEventsController(cacheCall, busEvent)

	router := httpserver.NewRouter(chi.NewRouter(), eventsController)
	server := http.Server{
		Addr:    addr,
		Handler: &router,
	}

	return server.ListenAndServe()
}

package main

import (
	"EventManager/internal/api/httpserver"
	v1 "EventManager/internal/api/v1"
	"EventManager/internal/cache"
	"github.com/go-chi/chi/v5"
	"log"
	"net"
	"net/http"
	"os"
)

const (
	defaultPort     = "9999"
	defaultHost     = "0.0.0.0"
	//defaultCacheDSN = "redis://eventscache:6379/0"
	defaultCacheDSN = "redis://localhost:6379/0"
	//PRIVATEKEY      = "./keys/private.key"
	//PUBLICKEY       = "./keys/public.key"
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

	if err := execute(net.JoinHostPort(host, port), cacheDSN); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func execute (addr, cacheDSN string) (err error) {
	cacheCallPool := cache.InitCache(cacheDSN)
	cacheCall := cache.NewCallCache(cacheCallPool)

	eventsController := v1.NewEventsController(cacheCall)

	router := httpserver.NewRouter(chi.NewRouter(), eventsController)
	server := http.Server{
		Addr:              addr,
		Handler:           &router,
	}

	return server.ListenAndServe()
}
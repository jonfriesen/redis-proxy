package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/jonfriesen/redis-proxy/api"
	"github.com/jonfriesen/redis-proxy/cache"
	"github.com/jonfriesen/redis-proxy/storage/redis"
)

func main() {
	cacheSize := flag.Int("cachesize", 10, "Represents the amount of keys that can be help in the cache")
	recordExpiry := flag.Int("expiry", 60000, "Represents the expiration limit in milliseconds")
	host := flag.String("host", "0.0.0.0", "Host IP address")
	port := flag.String("port", "4000", "Host port")
	redisHost := flag.String("redis-host", "0.0.0.0", "Redis IP address")
	redisPort := flag.String("redis-port", "6379", "Redis port")
	flag.Parse()

	cache := cache.New(int32(*cacheSize), time.Duration(*recordExpiry)*time.Millisecond)

	dataSource, err := redis.New(*redisHost, *redisPort)
	if err != nil {
		log.Fatalf("Failed to created Redis connection with %+v", err)
	}

	handler := api.New(dataSource, cache)

	server := &http.Server{
		Addr:    fmt.Sprintf("%v:%v", *host, *port),
		Handler: handler,
	}

	// graceful shutdown sequence
	go func() {
		sigquit := make(chan os.Signal, 1)
		signal.Notify(sigquit, os.Interrupt, os.Kill)

		sig := <-sigquit
		log.Printf("caught sig: %+v", sig)
		log.Printf("Gracefully shutting down server...")

		if err := server.Shutdown(context.Background()); err != nil {
			log.Printf("Unable to shut down server: %v", err)
		} else {
			log.Println("Server stopped")
		}
	}()

	log.Printf("Magic is happening on port %v", *port)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Printf("%v", err)
	} else {
		log.Println("Server closed")
	}

}

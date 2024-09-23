package cmd

import (
	"category-svc/feature/shared"
	"category-svc/feature/ticket"
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"time"
)

func runHTTPServer(ctx context.Context) {
	// Load configuration
	cfg := shared.LoadConfig("config/app.yaml")

	// Set needed dependencies
	cacheClient := redis.NewClient(&redis.Options{
		Addr:                  cfg.RedisConfig.Addr,
		MinIdleConns:          int(cfg.RedisConfig.MinConn),
		ConnMaxIdleTime:       5 * time.Minute,
		ContextTimeoutEnabled: true,
	})
	defer cacheClient.Close()

	ticket.SetCachePool(cacheClient)

	// Create a new server
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	ticket.HttpRoute(mux)

	srv := &http.Server{
		Addr:         cfg.Server.Addr(),
		Handler:      mux,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalln("unable to start server", err)
		}
	}()

	log.Println("server started")

	// Wait for signal to shut down
	<-ctx.Done()

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err := srv.Shutdown(ctxShutDown); err != nil {
		log.Fatalln("unable to shutdown server", err)
	}

	log.Println("server shutdown")
}

package server

import (
	"context"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"time"
)

/*
Supporting gracefully server shutdown alongside Golang's native *http and Gorilla Mux.
Implementation taken from https://github.com/gorilla/mux/blob/master/README.md#graceful-shutdown.
*/
func Start(port string, connectionTimeout time.Duration, onRouterSetup func() http.Handler) {
	srv := &http.Server{
		Addr: "0.0.0.0:" + port,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      onRouterSetup(),
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		log.Info("Starting server on port:", port)
		if err := srv.ListenAndServe(); err != nil {
			panic(err)
		}
	}()

	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	wait := time.Second * connectionTimeout
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)

	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	serverDidShutdown()
}

func serverDidShutdown() {
	log.Info("Server did shutdown")
	os.Exit(0)
}

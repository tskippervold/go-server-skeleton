package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/tskippervold/golang-base-server/internal/app/auth"
	"github.com/tskippervold/golang-base-server/internal/app/handlers"
	"github.com/tskippervold/golang-base-server/internal/db"
	"github.com/tskippervold/golang-base-server/internal/utils/log"

	env "github.com/tskippervold/golang-base-server/internal/app"

	"github.com/gorilla/mux"
)

func main() {
	logger := log.NewLogger()

	// Parsing the specified run argument `-config=`
	flagConfig := flag.String("config", "local.yml", "Config yaml file")
	flag.Parse()

	// Load a `Config` struct containing the yaml values.
	logger.Info("Loading config:", *flagConfig)
	config, err := env.LoadConfig(".././configs/" + *flagConfig)
	if err != nil {
		panic(err)
	}

	/*if err := db.ConnectToTest(); err != nil {
		panic(err)
	}*/

	logger.Info("Connecting to database:", config.Database.Host, config.Database.Name)
	if err := db.Connect(
		config.Database.Host,
		config.Database.Port,
		config.Database.Name,
		config.Database.User,
		config.Database.Pass,
	); err != nil {
		panic(err)
	}

	logger.Info("Migrating database")
	if err := db.Migrate(); err != nil {
		panic(err)
	}

	logger.Info("Starting server on port:", config.Server.Port)
	startServer(config.Server.Port, config.Server.ConnectionTimeout, onRouterSetup, onServerShutdown)
}

func onRouterSetup() http.Handler {
	env := env.Env{
		Log: log.NewLogger(),
		DB:  db.GetConnection(),
	}

	r := mux.NewRouter()
	r.Use(env.Log.HTTPRequestMiddleware)

	rAuth := r.PathPrefix("/auth").Subrouter()
	auth.Setup(rAuth, &env)

	rAPI := r.PathPrefix("/api").Subrouter()
	rAPI.Use(auth.JWTMiddleware)
	handlers.Setup(rAPI, &env)

	return r
}

/*
Supporting gracefully server shutdown alongside Golang's native *http and Gorilla Mux.
Implementation taken from https://github.com/gorilla/mux/blob/master/README.md#graceful-shutdown.
*/
func startServer(port string, connectionTimeout time.Duration, onRouterSetup func() http.Handler, onShutdown func()) {
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
	onShutdown()
}

func onServerShutdown() {
	logger := log.NewLogger()
	logger.Info("Server did shutdown")
	os.Exit(0)
}

package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	env "github.com/tskippervold/golang-base-server/internal/app"
	"github.com/tskippervold/golang-base-server/internal/app/handlers"
	"github.com/tskippervold/golang-base-server/internal/utils"
)

var logger = utils.NewLogger()

func main() {
	// Parsing the specified run argument `-config=`
	flagConfig := flag.String("config", "local.yml", "Config yaml file")
	flag.Parse()

	// Load a `Config` struct containing the yaml values.
	logger.Info("Loading config:", *flagConfig)
	config, err := env.LoadConfig("./configs/" + *flagConfig)
	if err != nil {
		panic(err)
	}

	/*
		Uncomment bellow when you are ready to use a Postgres database.
	*/
	/*logger.Info("Connecting to database:", config.Database.Host, config.Database.Name)
	dbConnection, err := env.ConnectDatabase(config.Database.Host, config.Database.Name, config.Database.Name, config.Database.User, config.Database.Pass)
	if err != nil {
		panic(err)
	}

	e := &env.Env{
		DB:  dbConnection,
		Log: logger,
	}

	logger.Info("Migrating database")
	if err := e.MigrateDatabase(); err != nil {
		panic(err)
	}*/

	e := &env.Env{
		Log: logger,
	}

	logger.Info("Starting server on port:", config.Server.Port)
	startServer(config.Server.Port, config.Server.ConnectionTimeout, func() http.Handler {
		r := mux.NewRouter()
		r.Use(e.Log.HTTPRequestMiddleware)
		handlers.Setup(r.PathPrefix("/api").Subrouter(), e)
		return r
	}, onServerShutdown)
}

/*
Supporting gracefully server shutdown alongside Golang's native *http and Gorilla Mux.
Implementation taken from https://github.com/gorilla/mux/blob/master/README.md#graceful-shutdown.
*/
func startServer(port string, connectionTimeout time.Duration, handler func() http.Handler, onShutdown func()) {
	srv := &http.Server{
		Addr: "0.0.0.0:" + port,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      handler(),
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
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
	logger.Info("Server did shutdown")
	os.Exit(0)
}

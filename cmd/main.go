package main

import (
	"flag"
	"net/http"

	"github.com/rs/cors"
	"github.com/tskippervold/golang-base-server/internal/app/auth"
	"github.com/tskippervold/golang-base-server/internal/app/handlers"
	"github.com/tskippervold/golang-base-server/internal/db"
	"github.com/tskippervold/golang-base-server/internal/server"
	"github.com/tskippervold/golang-base-server/internal/utils/log"

	env "github.com/tskippervold/golang-base-server/internal/app"

	"github.com/gorilla/mux"
)

func main() {
	logger := log.NewLogger()

	// Parsing the specified run argument `-config=`
	flagConfig := flag.String("config", "./configs/local.yml", "Config yaml file")
	flag.Parse()

	// Load a `Config` struct containing the yaml values.
	logger.Info("Loading config:", *flagConfig)
	config, err := env.LoadConfig(*flagConfig)
	if err != nil {
		logger.Error(err)
		panic(err)
	}

	logger.Info("Connecting to database:", config.Database.Host, config.Database.Name)

	/*if err := db.ConnectToTest(); err != nil {
		panic(err)
	}*/

	/*if err := db.Connect(
		config.Database.Host,
		config.Database.Port,
		config.Database.Name,
		config.Database.User,
		config.Database.Pass,
	); err != nil {
		logger.Error(err)
		panic(err)
	}

	logger.Info("Migrating database")
	if err := db.Migrate(); err != nil {
		logger.Error(err)
		panic(err)
	}*/

	logger.Info("Starting server on port:", config.Server.Port)
	server.Start(config.Server.Port, config.Server.ConnectionTimeout, onRouterSetup)
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

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
		// Enable Debugging for testing, consider disabling in production
		Debug: false,
	})

	return c.Handler(r)
}

package main

import (
	"flag"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
	"github.com/tskippervold/go-server-skeleton/internal/app"
	"github.com/tskippervold/go-server-skeleton/internal/app/handlers"
	"github.com/tskippervold/go-server-skeleton/internal/server"
	"github.com/tskippervold/go-server-skeleton/internal/utils"
	"net/http"
)

func main() {
	// Parsing the specified run argument `-config=`
	flagConfig := flag.String("config", "", "Config yaml file")
	flag.Parse()

	if *flagConfig == "" {
		log.Panic("No config yaml file provided. Please run application with `-config=<path_to_yaml>`")
	}

	// Load a `Config` struct containing the yaml values.
	log.Info("Loading config:", *flagConfig)
	config, err := app.LoadConfig(*flagConfig)
	if err != nil {
		log.Panic(err)
	}

	/*log.Info("Connecting to database:", config.Database.Host, config.Database.Name)

	if err := db.Connect(
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

	server.Start(config.Server.Port, config.Server.ConnectionTimeout, onRouterSetup)
}

func onRouterSetup() http.Handler {
	env := app.Environment{}
	r := mux.NewRouter()
	r.Use(utils.LogRequestsMiddleware)

	rAPI := r.PathPrefix("/api").Subrouter()
	handlers.Setup(rAPI, &env)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{"GET", "PUT", "POST", "DELETE"},
		// Enable Debugging for testing, consider disabling in production
		Debug: false,
	})

	return c.Handler(r)
}

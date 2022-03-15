package main

import (
	"os"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	_ "github.com/joho/godotenv/autoload"
	"gitlab.autoiterative.com/group-zealous-ishizaka-gates/backendify/cache"
	"gitlab.autoiterative.com/group-zealous-ishizaka-gates/backendify/providers"
	"gitlab.autoiterative.com/group-zealous-ishizaka-gates/backendify/routes"
	"gitlab.autoiterative.com/group-zealous-ishizaka-gates/backendify/server"
)

var (
	serverPort string
)

func init() {
	serverPort = os.Getenv("SERVER_PORT")
}

func main() {
	pdrs := providers.New(os.Args[1:])

	s := server.New(
		server.UseMidlewares(
			server.HealthcheckMiddleware(), // put it first to avoid log the healthcheck, internally match with the /status endpoint.
			middleware.StripSlashes,        // match paths with a trailing slash, strip it, and continue routing through the mux
			middleware.Recoverer,           // recover from panics without crashing server
		),
		server.ListenOn(serverPort), // if serverPort is empty by default it takes the port 9000
	)

	c := cache.New(24*time.Hour, 0)

	s.Route("/", func(r chi.Router) {
		// before to attend the request we need to be sure that the
		r.Use(server.ValidateQueryParametersMiddleware([]routes.RequiredQueryParameter{routes.CompanyID, routes.CountryCode}))

		// Register the routes
		r.Get("/company", routes.CompanyRoute(pdrs, c))
	})

	// start the server
	s.Start()
}

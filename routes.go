package main

import (
	"log"

	"github.com/ant0ine/go-json-rest/rest"
	. "github.com/yetis-br/travelPlanning/models"
)

// NewRouter return the configured routes for this API
func NewRouter() rest.App {
	trips := Trips{
		Store: map[string]*Trip{},
	}

	router, err := rest.MakeRouter(
		rest.Post("/login", jwtMiddleware.LoginHandler),
		rest.Get("/refresh_token", jwtMiddleware.RefreshHandler),
		rest.Get("/trips", trips.GetAllTrips),
		rest.Get("/trips/:id", trips.GetTrip),
		rest.Post("/trips", trips.PostTrip),
	)
	if err != nil {
		log.Fatal(err)
	}

	return router
}

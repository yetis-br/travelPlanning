package main

import (
	"log"

	"github.com/StephanDollberg/go-json-rest-middleware-jwt"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/yetis-br/travelPlanning/models"
)

// NewRouter return the configured routes for this API
func NewRouter(jwtmiddleware *jwt.JWTMiddleware) rest.App {
	trips := models.Trips{
		Store: map[string]*models.Trip{},
	}

	router, err := rest.MakeRouter(
		rest.Post("/login", jwtmiddleware.LoginHandler),
		rest.Get("/refresh_token", jwtmiddleware.RefreshHandler),
		rest.Get("/trips", trips.GetAllTrips),
		rest.Post("/trips", trips.PostTrip),
		rest.Get("/trips/:id", trips.GetTrip),
	)
	if err != nil {
		log.Fatal(err)
	}

	return router
}

package main

import (
	"log"

	"github.com/StephanDollberg/go-json-rest-middleware-jwt"
	"github.com/ant0ine/go-json-rest/rest"
	r "github.com/dancannon/gorethink"
	"github.com/yetis-br/travelPlanning/models"
)

// NewRouter return the configured routes for this API
func NewRouter(jwtmiddleware *jwt.JWTMiddleware, session *r.Session) rest.App {

	trips := models.Trips{
		Conn: session,
	}

	router, err := rest.MakeRouter(
		rest.Post("/login", jwtmiddleware.LoginHandler),
		rest.Get("/refresh_token", jwtmiddleware.RefreshHandler),
		rest.Get("/trips", trips.GetAllTrips),
		rest.Get("/trips/:id", trips.GetTrip),
		rest.Post("/trips", trips.PostTrip),
	)
	if err != nil {
		log.Fatal(err)
	}

	return router
}

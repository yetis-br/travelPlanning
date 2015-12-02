package main

import (
	"log"
	"travelPlanning/models"

	"github.com/ant0ine/go-json-rest/rest"
)

// NewRouter return the configured routes for this API
func NewRouter() rest.App {
	travels := models.Travels{}

	router, err := rest.MakeRouter(
		rest.Get("/travels", travels.GetAllTravels),
	)
	if err != nil {
		log.Fatal(err)
	}

	return router
}

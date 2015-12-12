package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/yetis-br/travelPlanning/models"
)

var _ = Describe("TravelPlanning API", func() {

	var loginCredentials = map[string]string{"username": "admin", "password": "admin"}
	var token string

	Describe("Interacting with login functionality", func() {
		Context("If login works correctly", func() {
			token = Login(loginCredentials)
			It("Return 200 Status", func() {
				Expect(Response.Code).To(Equal(200))
			})
			It("The created Token is not empty", func() {
				Expect(token).NotTo(BeEmpty())
			})
		})
	})

	trips := models.Trips{
		Store: map[string]*models.Trip{},
	}

	Describe("Working with Trips functionality", func() {
		Context("When loading all trips without token", func() {
			It("Return 401 Status", func() {
				APIRequest("/Trips", trips.GetAllTrips, "GET", nil, "")
				Expect(Response.Code).To(Equal(401))
			})
		})

		Context("When loading all trips with token", func() {
			It("Return 200 Status and Empty", func() {
				APIRequest("/Trips", trips.GetAllTrips, "GET", nil, token)
				Expect(Response.Code).To(Equal(200))
			})
			It("Return zero records", func() {
				APIRequest("/Trips", trips.GetAllTrips, "GET", nil, token)
				Expect(Response.Body).To(MatchJSON(`[]`))
			})
		})
	})
})

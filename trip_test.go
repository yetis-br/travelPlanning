package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/yetis-br/travelPlanning/models"
)

var _ = Describe("Trip API", func() {

	var token string = Login()

	Describe("Working with Trips functionality", func() {
		Context("When loading all trips without token", func() {
			It("Return 401 Status", func() {
				APIRequest("/trips", "GET", nil, "")
				Expect(Response.Code).To(Equal(401))
			})
		})

		Context("When loading all trips with token", func() {
			It("Return 200 Status and Empty", func() {
				APIRequest("/trips", "GET", nil, token)
				Expect(Response.Code).To(Equal(200))
			})
			It("Return zero records", func() {
				APIRequest("/trips", "GET", nil, token)
				Expect(Response.Body).To(MatchJSON(`[]`))
			})
		})

		Context("When posting a new trip with token", func() {
			trip := models.Trip{}
			trip.Title = "Europa 2017"

			It("Return 200 Status", func() {
				APIRequest("/trips", "POST", trip, token)
				Expect(Response.Code).To(Equal(200))
			})
			It("Return one record", func() {
				Expect(Response.Body).To(MatchJSON(`{"id": "0", "title": "Europa 2017"}`))
			})
		})
	})
})

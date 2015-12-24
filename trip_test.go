package main

import (
	"encoding/json"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/yetis-br/travelPlanning/models"
)

var _ = Describe("Trip API", func() {

	var token string = Login()
	var id string

	Describe("Working with Trips functionality", func() {
		Context("When loading all trips without token", func() {
			It("Return 401 Status", func() {
				APIRequest("/trips", "GET", nil, "")
				Expect(Response.Code).To(Equal(401))
			})
		})

		Context("When fetching all trips with token", func() {
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
			tripInput := models.Trip{}
			tripInput.Title = "Europa 2017"

			It("Return 201 Status", func() {
				APIRequest("/trips", "POST", tripInput, token)
				Expect(Response.Code).To(Equal(201))
			})
			It("Return one record", func() {
				tripOutput := models.Trip{}
				json.Unmarshal(Response.Body.Bytes(), &tripOutput)
				Expect(tripInput.Title).To(Equal(tripOutput.Title))
				id = tripOutput.ID
			})
		})

		Context("When fetching only one trip with token", func() {
			It("Return 200 Status", func() {
				APIRequest("/trips/"+id, "GET", nil, token)
				Expect(Response.Code).To(Equal(200))
			})
			It("Return one record", func() {
				tripOutput := models.Trip{}
				json.Unmarshal(Response.Body.Bytes(), &tripOutput)
				Expect("Europa 2017").To(Equal(tripOutput.Title))
			})
		})

		Context("When adding a new trip place", func() {
			places := []models.TripPlace{
				{ID: "00000001", Status: "active", Order: 0},
				{ID: "00000002", Status: "active", Order: 1},
				{ID: "00000003", Status: "active", Order: 2},
				{ID: "00000004", Status: "active", Order: 3},
			}

			It("Return 200 Status", func() {
				APIRequest("/trips/"+id+"/updatePlaces", "PATCH", places, token)
				Expect(Response.Code).To(Equal(200))
			})
			It("Return four places", func() {
				tripOutput := models.Trip{}
				json.Unmarshal(Response.Body.Bytes(), &tripOutput)
				Expect(4).To(Equal(tripOutput.TotalPlaces))
			})
		})

		Context("When delete a trip place", func() {
			It("Return 200 Status", func() {
				APIRequest("/trips/"+id+"/deletePlace/3", "DELETE", nil, token)
				Expect(Response.Code).To(Equal(200))
			})
			It("Return four places", func() {
				tripOutput := models.Trip{}
				APIRequest("/trips/"+id, "GET", nil, token)
				json.Unmarshal(Response.Body.Bytes(), &tripOutput)
				Expect(3).To(Equal(tripOutput.TotalPlaces))
			})
		})
	})
})

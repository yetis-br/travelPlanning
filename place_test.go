package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Place API", func() {

	//var token string = Login()

	Describe("Working with Places functionality", func() {
		Context("When loading all trips without token", func() {
			It("Return 401 Status", func() {
				APIRequest("/places", "GET", nil, "")
				Expect(Response.Code).To(Equal(401))
			})
		})
	})
})

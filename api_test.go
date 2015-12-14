package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TravelPlanning API", func() {

	var token string = Login()

	Describe("Interacting with login functionality", func() {
		Context("If login works correctly", func() {
			It("Return 200 Status", func() {
				Expect(Response.Code).To(Equal(200))
			})
			It("The created Token is not empty", func() {
				Expect(token).NotTo(BeEmpty())
			})
		})
	})
})

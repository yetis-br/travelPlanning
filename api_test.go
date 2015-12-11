package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var loginCredentials = map[string]string{"username": "admin", "password": "admin"}

var _ = Describe("Interacting with login functionality", func() {
	Context("If login works correctly", func() {
		token := Login(loginCredentials)
		It("Return 200 Status", func() {
			Expect(Response.Code).To(Equal(200))
		})
		It("The created Token is not empty", func() {
			Expect(token).NotTo(BeEmpty())
		})
	})
})

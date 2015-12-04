package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestTravelPlanning(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "TravelPlanning Suite")
}

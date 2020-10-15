package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSezzleTest(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "SezzleTest Suite")
}

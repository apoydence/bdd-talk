package bdd_talk_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestBddTalk(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "BddTalk Suite")
}

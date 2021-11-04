package signatures_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSignatures(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Signatures Test Suite")
}

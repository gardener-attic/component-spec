// SPDX-FileCopyrightText: 2021 SAP SE or an SAP affiliate company and Gardener contributors.
//
// SPDX-License-Identifier: Apache-2.0

package oci_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	cdv2 "github.com/gardener/component-spec/bindings-go/apis/v2"
	"github.com/gardener/component-spec/bindings-go/oci"
)

func TestConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "oci Test Suite")
}

var _ = Describe("helper", func(){


	Context("OCIRef", func() {

		It("should correctly parse a repository url without a protocol and a component", func() {
			repoCtx := cdv2.RepositoryContext{BaseURL: "example.com"}
			ref, err := oci.OCIRef(repoCtx, "somecomp", "v0.0.0")
			Expect(err).ToNot(HaveOccurred())
			Expect(ref).To(Equal("example.com/component-descriptors/somecomp:v0.0.0"))
		})

		It("should correctly parse a repository url with a protocol and a component", func() {
			repoCtx := cdv2.RepositoryContext{BaseURL: "http://example.com"}
			ref, err := oci.OCIRef(repoCtx, "somecomp", "v0.0.0")
			Expect(err).ToNot(HaveOccurred())
			Expect(ref).To(Equal("example.com/component-descriptors/somecomp:v0.0.0"))
		})

		It("should correctly parse a repository url without a protocol and a port and a component", func() {
			repoCtx := cdv2.RepositoryContext{BaseURL: "example.com:443"}
			ref, err := oci.OCIRef(repoCtx, "somecomp", "v0.0.0")
			Expect(err).ToNot(HaveOccurred())
			Expect(ref).To(Equal("example.com:443/component-descriptors/somecomp:v0.0.0"))
		})

		It("should correctly parse a repository url with a protocol and a port and a component", func() {
			repoCtx := cdv2.RepositoryContext{BaseURL: "http://example.com:443"}
			ref, err := oci.OCIRef(repoCtx, "somecomp", "v0.0.0")
			Expect(err).ToNot(HaveOccurred())
			Expect(ref).To(Equal("example.com:443/component-descriptors/somecomp:v0.0.0"))
		})


	})

})

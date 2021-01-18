// SPDX-FileCopyrightText: 2020 SAP SE or an SAP affiliate company and Gardener contributors.
//
// SPDX-License-Identifier: Apache-2.0

package cdutils_test

import (
	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	cdv2 "github.com/gardener/component-spec/bindings-go/apis/v2"
	"github.com/gardener/component-spec/bindings-go/apis/v2/cdutils"
	"github.com/gardener/component-spec/bindings-go/codec"
)

var _ = Describe("resource utils", func() {

	Context("#GetImageReferenceFromList", func() {
		It("should return the image from a component descriptor list", func() {
			data, err := ioutil.ReadFile("../../../../language-independent/test-resources/component_descriptor_v2.yaml")
			Expect(err).ToNot(HaveOccurred())
			cd := cdv2.ComponentDescriptor{}
			Expect(codec.Decode(data, &cd)).To(Succeed())

			imageAccess, err := cdutils.GetImageReferenceFromList(
				&cdv2.ComponentDescriptorList{Components: []cdv2.ComponentDescriptor{cd}},
				"github.com/gardener/gardener", "apiserver")
			Expect(err).ToNot(HaveOccurred())
			Expect(imageAccess).To(Equal("eu.gcr.io/gardener-project/gardener/apiserver:v1.7.4"))
		})

		It("should return an error if no component matches the given name", func() {
			data, err := ioutil.ReadFile("../../../../language-independent/test-resources/component_descriptor_v2.yaml")
			Expect(err).ToNot(HaveOccurred())
			cd := cdv2.ComponentDescriptor{}
			Expect(codec.Decode(data, &cd)).To(Succeed())

			_, err = cdutils.GetImageReferenceFromList(
				&cdv2.ComponentDescriptorList{Components: []cdv2.ComponentDescriptor{cd}},
				"github.com/gardener/nocomp", "apiserver")
			Expect(err).To(HaveOccurred())
		})
	})

	Context("#GetImageReferenceByName", func() {
		It("should return the image from a component descriptor", func() {
			data, err := ioutil.ReadFile("../../../../language-independent/test-resources/component_descriptor_v2.yaml")
			Expect(err).ToNot(HaveOccurred())
			cd := &cdv2.ComponentDescriptor{}
			Expect(codec.Decode(data, cd)).To(Succeed())

			imageAccess, err := cdutils.GetImageReferenceByName(cd, "apiserver")
			Expect(err).ToNot(HaveOccurred())

			Expect(imageAccess).To(Equal("eu.gcr.io/gardener-project/gardener/apiserver:v1.7.4"))
		})

		It("should return an error if no resource matches the given name", func() {
			data, err := ioutil.ReadFile("../../../../language-independent/test-resources/component_descriptor_v2.yaml")
			Expect(err).ToNot(HaveOccurred())
			cd := &cdv2.ComponentDescriptor{}
			Expect(codec.Decode(data, cd)).To(Succeed())

			_, err = cdutils.GetImageReferenceByName(cd, "noimage")
			Expect(err).To(HaveOccurred())
		})
	})

})

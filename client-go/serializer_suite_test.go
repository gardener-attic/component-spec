// Copyright 2020 Copyright (c) 2020 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package client_go_test

import (
	"io/ioutil"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	client_go "github.com/gardener/component-spec/client-go"
	v2 "github.com/gardener/component-spec/client-go/apis/v2"
)

func TestConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Utils Test Suite")
}

var _ = Describe("serializer", func() {

	It("should decode simple component descriptor without overwrites", func() {
		data, err := ioutil.ReadFile("./testdata/01-data.yaml")
		Expect(err).ToNot(HaveOccurred())

		var cd v2.ComponentDescriptor
		err = client_go.Decode(data, &cd)
		Expect(err).ToNot(HaveOccurred())

		Expect(cd.Components).To(HaveLen(1))
		Expect(cd.Components[0]).To(BeAssignableToTypeOf(&v2.GardenerComponent{}))
		Expect(cd.Components[0].GetDependencies()[0]).To(BeAssignableToTypeOf(&v2.OCIImage{}))
	})

	It("should decode simple component descriptor and apply override", func() {
		data, err := ioutil.ReadFile("./testdata/02-data.yaml")
		Expect(err).ToNot(HaveOccurred())

		var cd v2.ComponentDescriptor
		err = client_go.Decode(data, &cd)
		Expect(err).ToNot(HaveOccurred())

		Expect(cd.Components).To(HaveLen(1))

		effective, err := cd.ApplyOverwrites()
		Expect(err).ToNot(HaveOccurred())

		Expect(effective.Components).To(HaveLen(1))
		Expect(effective.Components[0].GetDependencies()).To(ConsistOf(cd.OverwriteDeclarations[0].Overwrites[0].DependencyOverwrites[0]))
	})

	It("should decode a generic component", func() {
		data, err := ioutil.ReadFile("./testdata/06-generic-component.yaml")
		Expect(err).ToNot(HaveOccurred())

		var cd v2.ComponentDescriptor
		err = client_go.Decode(data, &cd)
		Expect(err).ToNot(HaveOccurred())

		Expect(cd.Components).To(HaveLen(1))
		Expect(cd.Components[0].GetDependencies()[0]).To(BeAssignableToTypeOf(&v2.GenericComponent{}))

		genComp := cd.Components[0].GetDependencies()[0].(*v2.GenericComponent)
		Expect(genComp.AdditionalData).To(HaveKeyWithValue("url", "example.com"))
	})

	It("should decode a generic component and overwrite its data", func() {
		data, err := ioutil.ReadFile("./testdata/07-generic-component-overwrite.yaml")
		Expect(err).ToNot(HaveOccurred())

		var cd v2.ComponentDescriptor
		err = client_go.Decode(data, &cd)
		Expect(err).ToNot(HaveOccurred())

		effective, err := cd.ApplyOverwrites()
		Expect(err).ToNot(HaveOccurred())

		genComp := effective.Components[0].GetDependencies()[0].(*v2.GenericComponent)
		Expect(genComp.AdditionalData).To(HaveKeyWithValue("url", "other.com"))
	})

	It("should throw an error if a non resolvable component is defined as resolvable component", func() {
		data, err := ioutil.ReadFile("./testdata/03-fail.yaml")
		Expect(err).ToNot(HaveOccurred())

		var cd v2.ComponentDescriptor
		err = client_go.Decode(data, &cd)
		Expect(err).To(HaveOccurred())
	})

	It("should overwrite additional attributes of a ociComponent", func() {
		data, err := ioutil.ReadFile("./testdata/04-ociComp-overwrite.yaml")
		Expect(err).ToNot(HaveOccurred())

		var cd v2.ComponentDescriptor
		err = client_go.Decode(data, &cd)
		Expect(err).ToNot(HaveOccurred())

		effective, err := cd.ApplyOverwrites()
		Expect(err).ToNot(HaveOccurred())

		Expect(effective.Components).To(HaveLen(1))
		Expect(effective.Components[0]).To(BeAssignableToTypeOf(&v2.OCIComponent{}))

		gardenerComp := effective.Components[0].(*v2.OCIComponent)
		Expect(gardenerComp.Repository).To(Equal("eu.gcr.io/gardener-project/gardener"))
	})

	It("should apply multiple overwrites to a component and return the latest overwrite", func() {
		data, err := ioutil.ReadFile("./testdata/05-ociComp-multiple-overwrites.yaml")
		Expect(err).ToNot(HaveOccurred())

		var cd v2.ComponentDescriptor
		err = client_go.Decode(data, &cd)
		Expect(err).ToNot(HaveOccurred())

		effective, err := cd.ApplyOverwrites()
		Expect(err).ToNot(HaveOccurred())

		Expect(effective.Components).To(HaveLen(1))
		Expect(effective.Components[0]).To(BeAssignableToTypeOf(&v2.OCIComponent{}))

		gardenerComp := effective.Components[0].(*v2.OCIComponent)
		Expect(gardenerComp.Repository).To(Equal("eu.gcr.io/gardener-project/latest/gardener"))
		Expect(effective.Components[0].GetDependencies()).To(ConsistOf(cd.OverwriteDeclarations[0].Overwrites[0].DependencyOverwrites[0]))
	})

})

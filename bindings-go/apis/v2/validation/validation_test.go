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

package validation_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/util/validation/field"

	v2 "github.com/gardener/component-spec/bindings-go/apis/v2"
	"github.com/gardener/component-spec/bindings-go/apis/v2/validation"
)

func TestConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "V2 Test Suite")
}

var _ = Describe("Validation", func() {

	var (
		ociImage1 *v2.OCIImage
		ociImage2 *v2.OCIImage
	)

	BeforeEach(func() {
		ociImage1 = &v2.OCIImage{
			ComponentMetadata: v2.ComponentMetadata{
				Type:    v2.OCIImageType,
				Name:    "image1",
				Version: "1.2.3",
			},
			ImageReference: "docker/image1:1.2.3",
		}
		ociImage2 = &v2.OCIImage{
			ComponentMetadata: v2.ComponentMetadata{
				Type:    v2.OCIImageType,
				Name:    "image2",
				Version: "1.2.3",
			},
			ImageReference: "docker/image2:1.2.3",
		}
	})

	Context("#ValidateDependencies", func() {

		var list v2.ComponentList

		BeforeEach(func() {
			list = v2.ComponentList{
				ociImage1, ociImage2,
			}
		})

		It("should pass if no components are defined", func() {
			errList := validation.ValidateDependencies(&field.Path{}, nil)
			Expect(errList).To(BeEmpty())
		})

		It("should pass if no duplicates are defined", func() {
			errList := validation.ValidateDependencies(&field.Path{}, list)
			Expect(errList).To(BeEmpty())
		})

		It("should forbid if there are duplicated components", func() {
			list = append(list, ociImage1)
			errList := validation.ValidateDependencies(&field.Path{}, list)
			Expect(errList).ToNot(BeEmpty())
		})
	})
})

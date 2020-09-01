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

package validation

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	"k8s.io/apimachinery/pkg/util/validation/field"

	v2 "github.com/gardener/component-spec/bindings-go/apis/v2"
)

func TestConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "V2 Test Suite")
}

var _ = Describe("Validation", func() {

	var (
		comp *v2.ComponentDescriptor

		ociImage1    *v2.Resource
		ociRegistry1 *v2.OCIRegistryAccess
		ociImage2    *v2.Resource
		ociRegistry2 *v2.OCIRegistryAccess
	)

	BeforeEach(func() {
		ociRegistry1 = &v2.OCIRegistryAccess{
			ObjectType: v2.ObjectType{
				Type: v2.OCIRegistryType,
			},
			ImageReference: "docker/image1:1.2.3",
		}
		ociImage1 = &v2.Resource{
			ObjectMeta: v2.ObjectMeta{
				Name:    "image1",
				Version: "1.2.3",
			},
			ObjectType: v2.ObjectType{
				Type: v2.OCIImageType,
			},
			Access: ociRegistry1,
		}
		ociRegistry2 = &v2.OCIRegistryAccess{
			ObjectType: v2.ObjectType{
				Type: v2.OCIRegistryType,
			},
			ImageReference: "docker/image1:1.2.3",
		}
		ociImage2 = &v2.Resource{
			ObjectMeta: v2.ObjectMeta{
				Name:    "image2",
				Version: "1.2.3",
			},
			ObjectType: v2.ObjectType{
				Type: v2.OCIImageType,
			},
			Access: ociRegistry2,
		}

		comp = &v2.ComponentDescriptor{
			Metadata: v2.Metadata{
				Version: v2.SchemaVersion,
			},
			ComponentSpec: v2.ComponentSpec{
				ObjectMeta: v2.ObjectMeta{
					Name:    "my-comp",
					Version: "1.2.3",
				},
				Provider:            v2.ExternalProvider,
				RepositoryContexts:  nil,
				Sources:             nil,
				ComponentReferences: nil,
				LocalResources:      nil,
				ExternalResources:   []v2.Resource{*ociImage1, *ociImage2},
			},
		}
	})

	Context("#Metadata", func() {

		It("should forbid if the component schemaVersion is missing", func() {
			comp := v2.ComponentDescriptor{
				Metadata: v2.Metadata{},
			}

			errList := validate(nil, &comp)
			Expect(errList).To(ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
				"Type":  Equal(field.ErrorTypeRequired),
				"Field": Equal("meta.schemaVersion"),
			}))))
		})

		It("should pass if the component schemaVersion is defined", func() {
			errList := validate(nil, comp)
			Expect(errList).ToNot(ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
				"Type":  Equal(field.ErrorTypeRequired),
				"Field": Equal("meta.schemaVersion"),
			}))))
		})

	})

	Context("#Provider", func() {
		It("should forbid if a component's provider is invalid", func() {
			comp.Provider = "custom"
			errList := validate(nil, comp)
			Expect(errList).To(ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
				"Type":  Equal(field.ErrorTypeInvalid),
				"Field": Equal("component.provider"),
			}))))
		})
	})

	Context("#ObjectMeta", func() {
		It("should forbid if the component's version is missing", func() {
			comp := v2.ComponentDescriptor{}
			errList := validate(nil, &comp)
			Expect(errList).To(ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
				"Type":  Equal(field.ErrorTypeRequired),
				"Field": Equal("component.name"),
			}))))
			Expect(errList).To(ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
				"Type":  Equal(field.ErrorTypeRequired),
				"Field": Equal("component.version"),
			}))))
		})

		It("should forbid if the component's name is missing", func() {
			comp := v2.ComponentDescriptor{}
			errList := validate(nil, &comp)
			Expect(errList).To(ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
				"Type":  Equal(field.ErrorTypeRequired),
				"Field": Equal("component.name"),
			}))))
		})

	})

	Context("#Sources", func() {
		It("should forbid if a duplicated component's source is defined", func() {
			comp.Sources = []v2.Resource{
				{
					ObjectMeta: v2.ObjectMeta{Name: "a"},
				},
				{
					ObjectMeta: v2.ObjectMeta{Name: "a"},
				},
			}
			errList := validate(nil, comp)
			Expect(errList).To(ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
				"Type":  Equal(field.ErrorTypeDuplicate),
				"Field": Equal("component.sources[1]"),
			}))))
		})
	})

	Context("#ComponentReferences", func() {
		It("should pass if a reference is set", func() {
			comp.ComponentReferences = []v2.ComponentReference{
				{
					ComponentName: "test",
					ObjectMeta: v2.ObjectMeta{
						Name:    "test",
						Version: "1.2.3",
					},
				},
			}
			errList := validate(nil, comp)
			Expect(errList).ToNot(ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
				"Type":  Equal(field.ErrorTypeRequired),
				"Field": Equal("component.componentReferences[0].name"),
			}))))
			Expect(errList).ToNot(ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
				"Type":  Equal(field.ErrorTypeRequired),
				"Field": Equal("component.componentReferences[0].version"),
			}))))
		})

		It("should forbid if a reference's name is missing", func() {
			comp.ComponentReferences = []v2.ComponentReference{
				{
					ComponentName: "test",
					ObjectMeta: v2.ObjectMeta{
						Version: "1.2.3",
					},
				},
			}
			errList := validate(nil, comp)
			Expect(errList).To(ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
				"Type":  Equal(field.ErrorTypeRequired),
				"Field": Equal("component.componentReferences[0].name"),
			}))))
		})

		It("should forbid if a reference's component name is missing", func() {
			comp.ComponentReferences = []v2.ComponentReference{
				{
					ObjectMeta: v2.ObjectMeta{
						Name:    "test",
						Version: "1.2.3",
					},
				},
			}
			errList := validate(nil, comp)
			Expect(errList).To(ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
				"Type":  Equal(field.ErrorTypeRequired),
				"Field": Equal("component.componentReferences[0].componentName"),
			}))))
		})

		It("should forbid if a reference's version is missing", func() {
			comp.ComponentReferences = []v2.ComponentReference{
				{
					ComponentName: "test",
					ObjectMeta: v2.ObjectMeta{
						Name: "test",
					},
				},
			}
			errList := validate(nil, comp)
			Expect(errList).To(ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
				"Type":  Equal(field.ErrorTypeRequired),
				"Field": Equal("component.componentReferences[0].version"),
			}))))
		})

		It("should forbid if a duplicated component reference is defined", func() {
			comp.ComponentReferences = []v2.ComponentReference{
				{
					ObjectMeta: v2.ObjectMeta{
						Name: "test",
					},
				},
				{
					ObjectMeta: v2.ObjectMeta{
						Name: "test",
					},
				},
			}
			errList := validate(nil, comp)
			Expect(errList).To(ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
				"Type":  Equal(field.ErrorTypeDuplicate),
				"Field": Equal("component.componentReferences[1]"),
			}))))
		})
	})

	Context("#LocalResources", func() {
		It("should forbid if a local resource's version differs from the version of the parent", func() {
			comp.LocalResources = []v2.Resource{
				{
					ObjectMeta: v2.ObjectMeta{
						Name:    "locRes",
						Version: "0.0.1",
					},
				},
			}
			errList := validate(nil, comp)
			Expect(errList).To(ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
				"Type":  Equal(field.ErrorTypeInvalid),
				"Field": Equal("component.localResources[0].version"),
			}))))
		})

		It("should forbid if a duplicated local resource is defined", func() {
			comp.LocalResources = []v2.Resource{
				{
					ObjectMeta: v2.ObjectMeta{
						Name: "test",
					},
				},
				{
					ObjectMeta: v2.ObjectMeta{
						Name: "test",
					},
				},
			}
			errList := validate(nil, comp)
			Expect(errList).To(ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
				"Type":  Equal(field.ErrorTypeDuplicate),
				"Field": Equal("component.localResources[1]"),
			}))))
		})
	})

	Context("#ExternalResources", func() {
		It("should forbid if a duplicated local resource is defined", func() {
			comp.ExternalResources = []v2.Resource{
				{
					ObjectMeta: v2.ObjectMeta{
						Name: "test",
					},
				},
				{
					ObjectMeta: v2.ObjectMeta{
						Name: "test",
					},
				},
			}
			errList := validate(nil, comp)
			Expect(errList).To(ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
				"Type":  Equal(field.ErrorTypeDuplicate),
				"Field": Equal("component.externalResources[1]"),
			}))))
		})
	})

	Context("#labels", func() {

		It("should forbid if labels are defined multiple times in the same context", func() {
			comp.ComponentReferences = []v2.ComponentReference{
				{
					ComponentName: "test",
					ObjectMeta: v2.ObjectMeta{
						Name:    "test",
						Version: "1.2.3",
						Labels: []v2.Label{
							{
								Name:  "l1",
								Value: []byte{},
							},
							{
								Name:  "l1",
								Value: []byte{},
							},
						},
					},
				},
			}

			errList := validate(nil, comp)
			Expect(errList).To(ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
				"Type":  Equal(field.ErrorTypeDuplicate),
				"Field": Equal("component.componentReferences[0].labels[1]"),
			}))))
		})

		It("should pass if labels are defined multiple times in the same context with differnet names", func() {
			comp.ComponentReferences = []v2.ComponentReference{
				{
					ComponentName: "test",
					ObjectMeta: v2.ObjectMeta{
						Name:    "test",
						Version: "1.2.3",
						Labels: []v2.Label{
							{
								Name:  "l1",
								Value: []byte{},
							},
							{
								Name:  "l2",
								Value: []byte{},
							},
						},
					},
				},
			}

			errList := validate(nil, comp)
			Expect(errList).ToNot(ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
				"Type":  Equal(field.ErrorTypeDuplicate),
				"Field": Equal("component.componentReferences[0].labels[1]"),
			}))))
		})

	})
})

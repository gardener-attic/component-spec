// Copyright 2022 Copyright (c) 2022 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file.
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

package signatures_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	cdv2 "github.com/gardener/component-spec/bindings-go/apis/v2"
	"github.com/gardener/component-spec/bindings-go/apis/v2/signatures"
)

var _ = Describe("Normalise/Hash component-descriptor", func() {
	var baseCd cdv2.ComponentDescriptor
	correctBaseCdHash := "aa32547cf0cbead58bc9a27d6c0545d6a4965f9ff2de9f09ce1e6d777f53fbaf"
	//corresponding normalised CD:
	//[{"component":[{"componentReferences":[[{"componentName":"compRefNameComponentName"},{"digest":[{"hashAlgorithm":"sha256"},{"normalisationAlgorithm":"jsonNormalisation/v1"},{"value":"00000000000000"}]},{"extraIdentity":[{"refKey":"refName"}]},{"name":"compRefName"},{"version":"v0.0.2compRef"}],[{"componentName":"compRefNameComponentName"},{"digest":[{"hashAlgorithm":"sha256"},{"normalisationAlgorithm":"jsonNormalisation/v1"},{"value":"00000000000000"}]},{"name":"compRefWithNoExtraIdentity"},{"version":"v0.0.3compRef"}]]},{"name":"CD-Name"},{"provider":""},{"resources":[[{"digest":[{"hashAlgorithm":"sha256"},{"normalisationAlgorithm":"ociArtifactDigest/v1"},{"value":"00000000000000"}]},{"extraIdentity":[{"key":"value"}]},{"name":"Resource1"},{"relation":""},{"type":""},{"version":"v0.0.3resource"}],[{"digest":[{"hashAlgorithm":"sha256"},{"normalisationAlgorithm":"ociArtifactDigest/v1"},{"value":"00000000000000"}]},{"extraIdentity":[{"key":"value"}]},{"name":"ResourceWithNoExtraIdentity"},{"relation":""},{"type":""},{"version":"v0.0.4resource"}]]},{"version":"v0.0.1"}]}]
	BeforeEach(func() {
		baseCd = cdv2.ComponentDescriptor{
			Metadata: cdv2.Metadata{
				Version: "v2",
			},
			ComponentSpec: cdv2.ComponentSpec{
				ObjectMeta: cdv2.ObjectMeta{
					Name:    "CD-Name",
					Version: "v0.0.1",
				},
				ComponentReferences: []cdv2.ComponentReference{
					{
						Name:          "compRefName",
						ComponentName: "compRefNameComponentName",
						Version:       "v0.0.2compRef",
						ExtraIdentity: cdv2.Identity{
							"refKey": "refName",
						},
						Digest: &cdv2.DigestSpec{
							HashAlgorithm:          signatures.SHA256,
							NormalisationAlgorithm: string(cdv2.JsonNormalisationV1),
							Value:                  "00000000000000",
						},
					},
					{
						// ExtraIdentity is nil -> should be left out completely from normalisation
						Name:          "compRefWithNoExtraIdentity",
						ComponentName: "compRefNameComponentName",
						Version:       "v0.0.3compRef",
						Digest: &cdv2.DigestSpec{
							HashAlgorithm:          signatures.SHA256,
							NormalisationAlgorithm: string(cdv2.JsonNormalisationV1),
							Value:                  "00000000000000",
						},
					},
				},
				Resources: []cdv2.Resource{
					{
						IdentityObjectMeta: cdv2.IdentityObjectMeta{
							Name:    "Resource1",
							Version: "v0.0.3resource",
							ExtraIdentity: cdv2.Identity{
								"key": "value",
							},
						},
						Digest: &cdv2.DigestSpec{
							HashAlgorithm:          signatures.SHA256,
							NormalisationAlgorithm: string(cdv2.OciArtifactDigestV1),
							Value:                  "00000000000000",
						},
						Access: cdv2.NewUnstructuredType(cdv2.OCIRegistryType, map[string]interface{}{"imageRef": "ref"}),
					},
					{
						IdentityObjectMeta: cdv2.IdentityObjectMeta{
							Name:    "ResourceWithNoExtraIdentity",
							Version: "v0.0.4resource",
							ExtraIdentity: cdv2.Identity{
								"key": "value",
							},
						},
						Digest: &cdv2.DigestSpec{
							HashAlgorithm:          signatures.SHA256,
							NormalisationAlgorithm: string(cdv2.OciArtifactDigestV1),
							Value:                  "00000000000000",
						},
						Access: cdv2.NewUnstructuredType(cdv2.OCIRegistryType, map[string]interface{}{"imageRef": "ref:v0.0.4"}),
					},
				},
			},
		}
	})

	Describe("missing componentReference Digest", func() {
		It("should fail to hash", func() {
			baseCd.ComponentSpec.ComponentReferences[0].Digest = nil
			hasher, err := signatures.HasherForName(signatures.SHA256)
			Expect(err).To(BeNil())
			hash, err := signatures.HashForComponentDescriptor(baseCd, *hasher)
			Expect(hash).To(BeNil())
			Expect(err).ToNot(BeNil())
		})
	})
	Describe("should give the correct hash", func() {
		It("with sha256", func() {
			hasher, err := signatures.HasherForName(signatures.SHA256)
			Expect(err).To(BeNil())
			hash, err := signatures.HashForComponentDescriptor(baseCd, *hasher)
			Expect(err).To(BeNil())
			Expect(hash.Value).To(Equal(correctBaseCdHash))
		})
	})
	Describe("should remove empty component refs/resources lists during normalisation", func() {
		It("with sha256", func() {
			expectedHash := "44460fad9d46f9281018858a94bf80ae348e5c24ea4d9955ade89a16fb587edf"
			//corresponding normalised CD:
			//[{"component":[{"name":"CD-Name"},{"provider":""},{"version":"v0.0.1"}]}]
			cdWithEmptyLists := cdv2.ComponentDescriptor{
				Metadata: cdv2.Metadata{
					Version: "v2",
				},
				ComponentSpec: cdv2.ComponentSpec{
					ObjectMeta: cdv2.ObjectMeta{
						Name:    "CD-Name",
						Version: "v0.0.1",
					},
					// ComponentReferences & Resource empty -> should be left out completely from normalisation
					ComponentReferences: []cdv2.ComponentReference{},
					Resources:           []cdv2.Resource{},
				},
			}
			hasher, err := signatures.HasherForName(signatures.SHA256)
			Expect(err).To(BeNil())
			hash, err := signatures.HashForComponentDescriptor(cdWithEmptyLists, *hasher)
			Expect(err).To(BeNil())
			Expect(hash.Value).To(Equal(expectedHash))
		})
	})
	Describe("should ignore modifications in unhashed fields", func() {
		It("should succeed with signature changes", func() {
			baseCd.Signatures = append(baseCd.Signatures, cdv2.Signature{
				Name: "TestSig",
				Digest: cdv2.DigestSpec{
					HashAlgorithm:          signatures.SHA256,
					NormalisationAlgorithm: string(cdv2.JsonNormalisationV1),
					Value:                  "00000",
				},
				Signature: cdv2.SignatureSpec{
					Algorithm: "test",
					Value:     "0000",
				},
			})
			hasher, err := signatures.HasherForName(signatures.SHA256)
			Expect(err).To(BeNil())
			hash, err := signatures.HashForComponentDescriptor(baseCd, *hasher)
			Expect(err).To(BeNil())
			Expect(hash.Value).To(Equal(correctBaseCdHash))
		})
		It("should succeed with source changes", func() {
			baseCd.Sources = append(baseCd.Sources, cdv2.Source{
				IdentityObjectMeta: cdv2.IdentityObjectMeta{
					Name:    "source1",
					Version: "v0.0.0",
				},
			})
			hasher, err := signatures.HasherForName(signatures.SHA256)
			Expect(err).To(BeNil())
			hash, err := signatures.HashForComponentDescriptor(baseCd, *hasher)
			Expect(err).To(BeNil())
			Expect(hash.Value).To(Equal(correctBaseCdHash))
		})
		It("should succeed with resource access reference changes", func() {
			access, err := cdv2.NewUnstructured(cdv2.NewOCIRegistryAccess("ociRef/path/to/image"))
			Expect(err).To(BeNil())
			baseCd.Resources[0].Access = &access
			hasher, err := signatures.HasherForName(signatures.SHA256)
			Expect(err).To(BeNil())
			hash, err := signatures.HashForComponentDescriptor(baseCd, *hasher)
			Expect(err).To(BeNil())
			Expect(hash.Value).To(Equal(correctBaseCdHash))
		})

	})
	Describe("should correctly handle empty access and digest", func() {
		It("should be equal hash for access.type == None and access == nil", func() {
			baseCd.Resources[0].Access = nil
			baseCd.Resources[0].Digest = nil

			hasher, err := signatures.HasherForName(signatures.SHA256)
			Expect(err).To(BeNil())
			hash, err := signatures.HashForComponentDescriptor(baseCd, *hasher)
			Expect(err).To(BeNil())

			//add access to resource
			access := cdv2.NewEmptyUnstructured("None")
			Expect(err).To(BeNil())
			baseCd.Resources[0].Access = access
			hash2, err := signatures.HashForComponentDescriptor(baseCd, *hasher)
			Expect(err).To(BeNil())
			Expect(hash).To(Equal(hash2))
		})
		It("should fail if digest is empty", func() {
			baseCd.Resources[0].Digest = nil

			hasher, err := signatures.HasherForName(signatures.SHA256)
			Expect(err).To(BeNil())
			_, err = signatures.HashForComponentDescriptor(baseCd, *hasher)
			Expect(err).To(HaveOccurred())
		})
		It("should succed if digest is empty and access is nil", func() {
			baseCd.Resources[0].Access = nil
			baseCd.Resources[0].Digest = nil

			hasher, err := signatures.HasherForName(signatures.SHA256)
			Expect(err).To(BeNil())
			_, err = signatures.HashForComponentDescriptor(baseCd, *hasher)
			Expect(err).To(BeNil())
		})
		It("should fail if first is nil access and an access is added but a digest is missing", func() {
			baseCd.Resources[0].Access = nil
			baseCd.Resources[0].Digest = nil

			hasher, err := signatures.HasherForName(signatures.SHA256)
			Expect(err).To(BeNil())
			_, err = signatures.HashForComponentDescriptor(baseCd, *hasher)
			Expect(err).To(BeNil())

			//add ociRegistryAccess
			access, err := cdv2.NewUnstructured(cdv2.NewOCIRegistryAccess("ociRef/path/to/image"))
			Expect(err).To(BeNil())
			baseCd.Resources[0].Access = &access
			_, err = signatures.HashForComponentDescriptor(baseCd, *hasher)
			Expect(err).To(HaveOccurred())
		})
		It("should fail if first is none access.type and an access is added but a digest is missing", func() {
			baseCd.Resources[0].Access = cdv2.NewEmptyUnstructured("None")
			baseCd.Resources[0].Digest = nil

			hasher, err := signatures.HasherForName(signatures.SHA256)
			Expect(err).To(BeNil())
			_, err = signatures.HashForComponentDescriptor(baseCd, *hasher)
			Expect(err).To(BeNil())

			//add ociRegistryAccess
			access, err := cdv2.NewUnstructured(cdv2.NewOCIRegistryAccess("ociRef/path/to/image"))
			Expect(err).To(BeNil())
			baseCd.Resources[0].Access = &access
			_, err = signatures.HashForComponentDescriptor(baseCd, *hasher)
			Expect(err).To(HaveOccurred())
		})
		It("should fail if access is nil and digest is set", func() {
			baseCd.Resources[0].Access = nil

			hasher, err := signatures.HasherForName(signatures.SHA256)
			Expect(err).To(BeNil())
			_, err = signatures.HashForComponentDescriptor(baseCd, *hasher)
			Expect(err).To(HaveOccurred())
		})
		It("should fail if access.type is None and digest is set", func() {
			baseCd.Resources[0].Access = cdv2.NewEmptyUnstructured("None")

			hasher, err := signatures.HasherForName(signatures.SHA256)
			Expect(err).To(BeNil())
			_, err = signatures.HashForComponentDescriptor(baseCd, *hasher)
			Expect(err).To(HaveOccurred())
		})
	})
	Describe("add digest to cd", func() {
		It("should succed if existing digest match calculated", func() {
			err := signatures.AddDigestsToComponentDescriptor(context.TODO(), &baseCd, func(ctx context.Context, cd cdv2.ComponentDescriptor, cr cdv2.ComponentReference) (*cdv2.DigestSpec, error) {
				return &cdv2.DigestSpec{
					HashAlgorithm:          signatures.SHA256,
					NormalisationAlgorithm: string(cdv2.JsonNormalisationV1),
					Value:                  "00000000000000",
				}, nil
			}, func(ctx context.Context, cd cdv2.ComponentDescriptor, r cdv2.Resource) (*cdv2.DigestSpec, error) {
				return &cdv2.DigestSpec{
					HashAlgorithm:          signatures.SHA256,
					NormalisationAlgorithm: string(cdv2.OciArtifactDigestV1),
					Value:                  "00000000000000",
				}, nil
			})
			Expect(err).To(BeNil())
		})
		It("should fail if calcuated componentReference digest is different", func() {
			err := signatures.AddDigestsToComponentDescriptor(context.TODO(), &baseCd, func(ctx context.Context, cd cdv2.ComponentDescriptor, cr cdv2.ComponentReference) (*cdv2.DigestSpec, error) {
				return &cdv2.DigestSpec{
					HashAlgorithm:          signatures.SHA256,
					NormalisationAlgorithm: string(cdv2.JsonNormalisationV1),
					Value:                  "00000000000000-different",
				}, nil
			}, func(ctx context.Context, cd cdv2.ComponentDescriptor, r cdv2.Resource) (*cdv2.DigestSpec, error) {
				return &cdv2.DigestSpec{
					HashAlgorithm:          signatures.SHA256,
					NormalisationAlgorithm: string(cdv2.OciArtifactDigestV1),
					Value:                  "00000000000000",
				}, nil
			})
			Expect(err).To(HaveOccurred())
		})
		It("should fail if calcuated resource digest is different", func() {
			err := signatures.AddDigestsToComponentDescriptor(context.TODO(), &baseCd, func(ctx context.Context, cd cdv2.ComponentDescriptor, cr cdv2.ComponentReference) (*cdv2.DigestSpec, error) {
				return &cdv2.DigestSpec{
					HashAlgorithm:          signatures.SHA256,
					NormalisationAlgorithm: string(cdv2.JsonNormalisationV1),
					Value:                  "00000000000000",
				}, nil
			}, func(ctx context.Context, cd cdv2.ComponentDescriptor, r cdv2.Resource) (*cdv2.DigestSpec, error) {
				return &cdv2.DigestSpec{
					HashAlgorithm:          signatures.SHA256,
					NormalisationAlgorithm: string(cdv2.OciArtifactDigestV1),
					Value:                  "00000000000000-different",
				}, nil
			})
			Expect(err).To(HaveOccurred())
		})
		It("should add digest if missing", func() {
			baseCd.ComponentReferences[0].Digest = nil
			baseCd.Resources[0].Digest = nil

			err := signatures.AddDigestsToComponentDescriptor(context.TODO(), &baseCd, func(ctx context.Context, cd cdv2.ComponentDescriptor, cr cdv2.ComponentReference) (*cdv2.DigestSpec, error) {
				return &cdv2.DigestSpec{
					HashAlgorithm:          signatures.SHA256,
					NormalisationAlgorithm: string(cdv2.JsonNormalisationV1),
					Value:                  "00000000000000",
				}, nil
			}, func(ctx context.Context, cd cdv2.ComponentDescriptor, r cdv2.Resource) (*cdv2.DigestSpec, error) {
				return &cdv2.DigestSpec{
					HashAlgorithm:          signatures.SHA256,
					NormalisationAlgorithm: string(cdv2.OciArtifactDigestV1),
					Value:                  "00000000000000",
				}, nil
			})
			Expect(err).To(BeNil())

			Expect(baseCd.ComponentReferences[0].Digest).To(Equal(&cdv2.DigestSpec{
				HashAlgorithm:          signatures.SHA256,
				NormalisationAlgorithm: string(cdv2.JsonNormalisationV1),
				Value:                  "00000000000000",
			}))
			Expect(baseCd.Resources[0].Digest).To(Equal(&cdv2.DigestSpec{
				HashAlgorithm:          signatures.SHA256,
				NormalisationAlgorithm: string(cdv2.OciArtifactDigestV1),
				Value:                  "00000000000000",
			}))
		})
		It("should preserve the EXCLUDE-FROM-SIGNATURE digest", func() {
			baseCd.Resources[0].Digest = cdv2.NewExcludeFromSignatureDigest()

			err := signatures.AddDigestsToComponentDescriptor(context.TODO(), &baseCd, func(ctx context.Context, cd cdv2.ComponentDescriptor, cr cdv2.ComponentReference) (*cdv2.DigestSpec, error) {
				return &cdv2.DigestSpec{
					HashAlgorithm:          signatures.SHA256,
					NormalisationAlgorithm: string(cdv2.JsonNormalisationV1),
					Value:                  "00000000000000",
				}, nil
			}, func(ctx context.Context, cd cdv2.ComponentDescriptor, r cdv2.Resource) (*cdv2.DigestSpec, error) {
				return &cdv2.DigestSpec{
					HashAlgorithm:          signatures.SHA256,
					NormalisationAlgorithm: string(cdv2.OciArtifactDigestV1),
					Value:                  "00000000000000",
				}, nil
			})
			Expect(err).To(BeNil())

			Expect(baseCd.ComponentReferences[0].Digest).To(Equal(&cdv2.DigestSpec{
				HashAlgorithm:          signatures.SHA256,
				NormalisationAlgorithm: string(cdv2.JsonNormalisationV1),
				Value:                  "00000000000000",
			}))
			Expect(baseCd.Resources[0].Digest).To(Equal(cdv2.NewExcludeFromSignatureDigest()))
		})
	})
})

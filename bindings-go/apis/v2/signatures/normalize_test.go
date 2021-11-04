package signatures_test

import (
	"crypto/sha256"
	"encoding/hex"

	v2 "github.com/gardener/component-spec/bindings-go/apis/v2"
	"github.com/gardener/component-spec/bindings-go/apis/v2/signatures"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Normalise/Hash component-descriptor", func() {
	var baseCd v2.ComponentDescriptor
	correctBaseCdHash := "dcefe8d7b4a43f5d7569234adbd696dfbfe3f1fc402ac3bed3f5e587624b9ac0"
	//corresponding normalised CD:
	//[{"component":[{"componentReferences":[[{"digest":[{"algorithm":"sha256"},{"value":"00000000000000"}]},{"extraIdentity":[{"refKey":"refName"}]},{"name":"compRefName"},{"version":"v0.0.2compRef"}]]},{"name":"CD-Name"},{"resources":[[{"digest":[{"algorithm":"sha256"},{"value":"00000000000000"}]},{"extraIdentity":[{"key":"value"}]},{"name":"Resource1"},{"version":"v0.0.3resource"}]]},{"version":"v0.0.1"}]},{"meta":[{"schemaVersion":"v2"}]}]
	BeforeEach(func() {
		baseCd = v2.ComponentDescriptor{
			Metadata: v2.Metadata{
				Version: "v2",
			},
			ComponentSpec: v2.ComponentSpec{
				ObjectMeta: v2.ObjectMeta{
					Name:    "CD-Name",
					Version: "v0.0.1",
				},
				ComponentReferences: []v2.ComponentReference{
					{
						Name:          "compRefName",
						ComponentName: "compRefNameComponentName",
						Version:       "v0.0.2compRef",
						ExtraIdentity: v2.Identity{
							"refKey": "refName",
						},
						Digest: v2.DigestSpec{
							Algorithm: "sha256",
							Value:     "00000000000000",
						},
					},
				},
				Resources: []v2.Resource{
					{
						IdentityObjectMeta: v2.IdentityObjectMeta{
							Name:    "Resource1",
							Version: "v0.0.3resource",
							ExtraIdentity: v2.Identity{
								"key": "value",
							},
						},
						Digest: v2.DigestSpec{
							Algorithm: "sha256",
							Value:     "00000000000000",
						},
					},
				},
			},
		}
	})

	Describe("missing componentReference Digest", func() {
		It("should fail to hash", func() {
			baseCd.ComponentSpec.ComponentReferences[0].Digest = v2.DigestSpec{}
			hash, err := signatures.HashForComponentDescriptor(baseCd, sha256.New())
			Expect(hash).To(BeNil())
			Expect(err).ToNot(BeNil())
		})
	})
	Describe("missing resource Digest", func() {
		It("should fail to hash", func() {
			baseCd.ComponentSpec.Resources[0].Digest = v2.DigestSpec{}
			hash, err := signatures.HashForComponentDescriptor(baseCd, sha256.New())
			Expect(hash).To(BeNil())
			Expect(err).ToNot(BeNil())
		})
	})
	Describe("should give the correct hash hash", func() {
		It("with sha256", func() {
			hash, err := signatures.HashForComponentDescriptor(baseCd, sha256.New())
			Expect(err).To(BeNil())
			Expect(hex.EncodeToString(hash)).To(Equal(correctBaseCdHash))
		})
	})
	Describe("should ignore modifications in unhashed fields", func() {
		It("should succed with signature changes", func() {
			baseCd.Signatures = append(baseCd.Signatures, v2.Signature{
				Name:                 "TestSig",
				NormalisationVersion: "v1",
				Digest: v2.DigestSpec{
					Algorithm: "sha256",
					Value:     "00000",
				},
				Signature: v2.SignatureSpec{
					Algorithm: "test",
					Data:      "0000",
				},
			})
			hash, err := signatures.HashForComponentDescriptor(baseCd, sha256.New())
			Expect(err).To(BeNil())
			Expect(hex.EncodeToString(hash)).To(Equal(correctBaseCdHash))
		})
		It("should succed with source changes", func() {
			baseCd.Sources = append(baseCd.Sources, v2.Source{
				IdentityObjectMeta: v2.IdentityObjectMeta{
					Name:    "source1",
					Version: "v0.0.0",
				},
			})
			hash, err := signatures.HashForComponentDescriptor(baseCd, sha256.New())
			Expect(err).To(BeNil())
			Expect(hex.EncodeToString(hash)).To(Equal(correctBaseCdHash))
		})
		It("should succed with resource access reference changes", func() {
			access, err := v2.NewUnstructured(v2.NewOCIRegistryAccess("ociRef/path/to/image"))
			Expect(err).To(BeNil())
			baseCd.Resources[0].Access = &access
			hash, err := signatures.HashForComponentDescriptor(baseCd, sha256.New())
			Expect(err).To(BeNil())
			Expect(hex.EncodeToString(hash)).To(Equal(correctBaseCdHash))
		})
	})
})

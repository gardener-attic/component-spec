package signatures_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	v2 "github.com/gardener/component-spec/bindings-go/apis/v2"
	"github.com/gardener/component-spec/bindings-go/apis/v2/signatures"
)

var _ = Describe("Normalise/Hash component-descriptor", func() {
	var baseCd v2.ComponentDescriptor
	correctBaseCdHash := "f21bf11eb5109f39d63ef5b178dd1b5ef38f144340f18d661fb1fb42f78ea7ee"
	//corresponding normalised CD:
	//[{"component":[{"componentReferences":[[{"componentName":"compRefNameComponentName"},{"digest":[{"hashAlgorithm":"sha256"},{"normalisationAlgorithm":"jsonNormalisation/V1"},{"value":"00000000000000"}]},{"extraIdentity":[{"refKey":"refName"}]},{"name":"compRefName"},{"version":"v0.0.2compRef"}]]},{"name":"CD-Name"},{"provider":""},{"resources":[[{"digest":[{"hashAlgorithm":"sha256"},{"normalisationAlgorithm":"manifestDigest/V1"},{"value":"00000000000000"}]},{"extraIdentity":[{"key":"value"}]},{"name":"Resource1"},{"type",""},{"version":"v0.0.3resource"}]]},{"version":"v0.0.1"}]},{"meta":[{"schemaVersion":"v2"}]}]
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
						Digest: &v2.DigestSpec{
							HashAlgorithm:          signatures.SHA256,
							NormalisationAlgorithm: string(v2.JsonNormalisationV1),
							Value:                  "00000000000000",
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
						Digest: &v2.DigestSpec{
							HashAlgorithm:          signatures.SHA256,
							NormalisationAlgorithm: string(v2.ManifestDigestV1),
							Value:                  "00000000000000",
						},
						Access: v2.NewUnstructuredType(v2.OCIRegistryType, map[string]interface{}{"imageRef": "ref"}),
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
	Describe("should ignore modifications in unhashed fields", func() {
		It("should succeed with signature changes", func() {
			baseCd.Signatures = append(baseCd.Signatures, v2.Signature{
				Name: "TestSig",
				Digest: v2.DigestSpec{
					HashAlgorithm:          signatures.SHA256,
					NormalisationAlgorithm: string(v2.JsonNormalisationV1),
					Value:                  "00000",
				},
				Signature: v2.SignatureSpec{
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
			baseCd.Sources = append(baseCd.Sources, v2.Source{
				IdentityObjectMeta: v2.IdentityObjectMeta{
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
			access, err := v2.NewUnstructured(v2.NewOCIRegistryAccess("ociRef/path/to/image"))
			Expect(err).To(BeNil())
			baseCd.Resources[0].Access = &access
			hasher, err := signatures.HasherForName(signatures.SHA256)
			Expect(err).To(BeNil())
			hash, err := signatures.HashForComponentDescriptor(baseCd, *hasher)
			Expect(err).To(BeNil())
			Expect(hash.Value).To(Equal(correctBaseCdHash))
		})
	})
})

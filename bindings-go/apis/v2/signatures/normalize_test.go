package signatures_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	v2 "github.com/gardener/component-spec/bindings-go/apis/v2"
	"github.com/gardener/component-spec/bindings-go/apis/v2/signatures"
)

var _ = Describe("Normalise/Hash component-descriptor", func() {
	var baseCd v2.ComponentDescriptor
	correctBaseCdHash := "6c571bb6e351ae755baa7f26cbd1f600d2968ab8b88e25a3bab277e53afdc3ad"
	//corresponding normalised CD:
	//[{"component":[{"componentReferences":[[{"componentName":"compRefNameComponentName"},{"digest":[{"hashAlgorithm":"sha256"},{"normalisationAlgorithm":"jsonNormalisation/v1"},{"value":"00000000000000"}]},{"extraIdentity":[{"refKey":"refName"}]},{"name":"compRefName"},{"version":"v0.0.2compRef"}]]},{"name":"CD-Name"},{"provider":""},{"resources":[[{"digest":[{"hashAlgorithm":"sha256"},{"normalisationAlgorithm":"manifestDigest/v1"},{"value":"00000000000000"}]},{"extraIdentity":[{"key":"value"}]},{"name":"Resource1"},{"relation": ""},{"type",""},{"version":"v0.0.3resource"}]]},{"version":"v0.0.1"}]},{"meta":[{"schemaVersion":"v2"}]}]
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
							NormalisationAlgorithm: string(v2.OciArtifactDigestV1),
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
	Describe("should correctly handle empty access and digest", func() {
		It("should be equal hash for access.type == None and access == nil", func() {
			baseCd.Resources[0].Access = nil
			baseCd.Resources[0].Digest = nil

			hasher, err := signatures.HasherForName(signatures.SHA256)
			Expect(err).To(BeNil())
			hash, err := signatures.HashForComponentDescriptor(baseCd, *hasher)
			Expect(err).To(BeNil())

			//add access to resource
			access := v2.NewEmptyUnstructured("None")
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
			access, err := v2.NewUnstructured(v2.NewOCIRegistryAccess("ociRef/path/to/image"))
			Expect(err).To(BeNil())
			baseCd.Resources[0].Access = &access
			_, err = signatures.HashForComponentDescriptor(baseCd, *hasher)
			Expect(err).To(HaveOccurred())
		})
		It("should fail if first is none access.type and an access is added but a digest is missing", func() {
			baseCd.Resources[0].Access = v2.NewEmptyUnstructured("None")
			baseCd.Resources[0].Digest = nil

			hasher, err := signatures.HasherForName(signatures.SHA256)
			Expect(err).To(BeNil())
			_, err = signatures.HashForComponentDescriptor(baseCd, *hasher)
			Expect(err).To(BeNil())

			//add ociRegistryAccess
			access, err := v2.NewUnstructured(v2.NewOCIRegistryAccess("ociRef/path/to/image"))
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
			baseCd.Resources[0].Access = v2.NewEmptyUnstructured("None")

			hasher, err := signatures.HasherForName(signatures.SHA256)
			Expect(err).To(BeNil())
			_, err = signatures.HashForComponentDescriptor(baseCd, *hasher)
			Expect(err).To(HaveOccurred())
		})
	})
})

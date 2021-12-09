package signatures_test

import (
	"crypto/sha256"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	v2 "github.com/gardener/component-spec/bindings-go/apis/v2"
	"github.com/gardener/component-spec/bindings-go/apis/v2/signatures"
)

type TestSigner struct{}

func (s TestSigner) Sign(componentDescriptor v2.ComponentDescriptor, digest v2.DigestSpec) (*v2.SignatureSpec, error) {
	return &v2.SignatureSpec{
		Algorithm: "testSignAlgorithm",
		Value:     fmt.Sprintf("%s:%s-signed", digest.HashAlgorithm, digest.Value),
	}, nil
}

type TestVerifier struct{}

func (v TestVerifier) Verify(componentDescriptor v2.ComponentDescriptor, signature v2.Signature) error {
	if signature.Signature.Value != fmt.Sprintf("%s:%s-signed", signature.Digest.HashAlgorithm, signature.Digest.Value) {
		return fmt.Errorf("signature verification failed: Invalid signature")
	}
	return nil
}

type TestSHA256Hasher signatures.Hasher

var _ = Describe("Sign/Verify component-descriptor", func() {
	var baseCd v2.ComponentDescriptor
	testSHA256Hasher := signatures.Hasher{
		HashFunction:  sha256.New(),
		AlgorithmName: "sha256",
	}
	signatureName := "testSignatureName"
	correctBaseCdHash := "64c04405dcd03a6f345584adb860cad4f4ed6dba1943d5535db3407b2bf9f000"

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
							HashAlgorithm:          "sha256",
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
							HashAlgorithm:          "sha256",
							NormalisationAlgorithm: string(v2.ManifestDigestV1),
							Value:                  "00000000000000",
						},
					},
				},
			},
		}
	})

	Describe("sign component-descriptor", func() {
		It("should add one signature", func() {
			err := signatures.SignComponentDescriptor(&baseCd, TestSigner{}, testSHA256Hasher, signatureName)
			Expect(err).To(BeNil())
			Expect(len(baseCd.Signatures)).To(BeIdenticalTo(1))
			Expect(baseCd.Signatures[0].Name).To(BeIdenticalTo(signatureName))
			Expect(baseCd.Signatures[0].Digest.NormalisationAlgorithm).To(BeIdenticalTo(string(v2.JsonNormalisationV1)))
			Expect(baseCd.Signatures[0].Digest.HashAlgorithm).To(BeIdenticalTo("sha256"))
			Expect(baseCd.Signatures[0].Digest.Value).To(BeIdenticalTo(correctBaseCdHash))
			Expect(baseCd.Signatures[0].Signature.Algorithm).To(BeIdenticalTo("testSignAlgorithm"))
			Expect(baseCd.Signatures[0].Signature.Value).To(BeIdenticalTo(fmt.Sprintf("%s:%s-signed", "sha256", correctBaseCdHash)))
		})
	})
	Describe("verify component-descriptor signature", func() {
		It("should verify one signature", func() {
			err := signatures.SignComponentDescriptor(&baseCd, TestSigner{}, testSHA256Hasher, signatureName)
			Expect(err).To(BeNil())
			Expect(len(baseCd.Signatures)).To(BeIdenticalTo(1))
			err = signatures.VerifySignedComponentDescriptor(&baseCd, TestVerifier{}, signatureName)
			Expect(err).To(BeNil())
		})
		It("should reject an invalid signature", func() {
			err := signatures.SignComponentDescriptor(&baseCd, TestSigner{}, testSHA256Hasher, signatureName)
			Expect(err).To(BeNil())
			Expect(len(baseCd.Signatures)).To(BeIdenticalTo(1))
			baseCd.Signatures[0].Signature.Value = "invalidSignature"
			err = signatures.VerifySignedComponentDescriptor(&baseCd, TestVerifier{}, signatureName)
			Expect(err).ToNot(BeNil())
		})
		It("should reject a missing signature", func() {
			err := signatures.VerifySignedComponentDescriptor(&baseCd, TestVerifier{}, signatureName)
			Expect(err).ToNot(BeNil())
		})

		It("should validate the correct signature if multiple are present", func() {
			err := signatures.SignComponentDescriptor(&baseCd, TestSigner{}, testSHA256Hasher, signatureName)
			Expect(err).To(BeNil())
			Expect(len(baseCd.Signatures)).To(BeIdenticalTo(1))

			baseCd.Signatures = append(baseCd.Signatures, v2.Signature{
				Name: "testSignAlgorithmNOTRight",
				Digest: v2.DigestSpec{
					HashAlgorithm:          "testAlgorithm",
					NormalisationAlgorithm: string(v2.JsonNormalisationV1),
					Value:                  "testValue",
				},
				Signature: v2.SignatureSpec{
					Algorithm: "testSigning",
					Value:     "AdditionalSignature",
				},
			})
			err = signatures.VerifySignedComponentDescriptor(&baseCd, TestVerifier{}, signatureName)
			Expect(err).To(BeNil())
		})
	})

	Describe("verify normalised component-descriptor digest with signed digest ", func() {
		It("should reject an invalid hash", func() {
			err := signatures.SignComponentDescriptor(&baseCd, TestSigner{}, testSHA256Hasher, signatureName)
			Expect(err).To(BeNil())
			Expect(len(baseCd.Signatures)).To(BeIdenticalTo(1))
			baseCd.Signatures[0].Digest.Value = "invalidHash"
			err = signatures.VerifySignedComponentDescriptor(&baseCd, TestVerifier{}, signatureName)
			Expect(err).ToNot(BeNil())
		})
		It("should reject a missing hash", func() {
			err := signatures.SignComponentDescriptor(&baseCd, TestSigner{}, testSHA256Hasher, signatureName)
			Expect(err).To(BeNil())
			Expect(len(baseCd.Signatures)).To(BeIdenticalTo(1))
			baseCd.Signatures[0].Digest = v2.DigestSpec{}
			err = signatures.VerifySignedComponentDescriptor(&baseCd, TestVerifier{}, signatureName)
			Expect(err).ToNot(BeNil())
		})
	})
})

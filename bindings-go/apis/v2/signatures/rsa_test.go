package signatures_test

import (
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
	"os"
	"os/exec"
	"path"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	v2 "github.com/gardener/component-spec/bindings-go/apis/v2"
	"github.com/gardener/component-spec/bindings-go/apis/v2/signatures"
)

var _ = Describe("RSA Sign/Verify", func() {
	var pathPrivateKey string
	var pathPublicKey string
	var stringToHashAndSign string
	var dir string

	BeforeEach(func() {
		var err error
		dir, err = ioutil.TempDir("", "component-spec-test")
		Expect(err).To(BeNil())

		// openssl genrsa -out private.key 4096
		pathPrivateKey = path.Join(dir, "private.key")
		createPrivateKeyCommand := exec.Command("openssl", "genrsa", "-out", pathPrivateKey, "4096")
		err = createPrivateKeyCommand.Run()
		Expect(err).To(BeNil())

		// openssl rsa -in private.key -outform PEM -pubout -out public.key
		pathPublicKey = path.Join(dir, "public.key")
		createPublicKeyCommand := exec.Command("openssl", "rsa", "-in", pathPrivateKey, "-outform", "PEM", "-pubout", "-out", pathPublicKey)
		err = createPublicKeyCommand.Run()
		Expect(err).To(BeNil())

		stringToHashAndSign = "TestStringToSign"
	})

	AfterEach(func() {
		os.RemoveAll(dir)
	})

	Describe("RSA Sign with private key", func() {
		It("should create a signature", func() {
			hashOfString := sha256.Sum256([]byte(stringToHashAndSign))

			signer, err := signatures.CreateRsaSignerFromKeyFile(pathPrivateKey)
			Expect(err).To(BeNil())
			signature, err := signer.Sign(v2.ComponentDescriptor{}, v2.DigestSpec{
				HashAlgorithm:          "sha256",
				NormalisationAlgorithm: string(v2.JsonNormalisationV1),
				Value:                  hex.EncodeToString(hashOfString[:]),
			})
			Expect(err).To(BeNil())
			Expect(signature.Algorithm).To(BeIdenticalTo("RSASSA-PKCS1-V1_5-SIGN"))
			Expect(signature.Value).NotTo(BeNil())
		})
		It("should should fail on unknown Digest algorithm", func() {
			hashOfString := sha256.Sum256([]byte(stringToHashAndSign))

			signer, err := signatures.CreateRsaSignerFromKeyFile(pathPrivateKey)
			Expect(err).To(BeNil())
			signature, err := signer.Sign(v2.ComponentDescriptor{}, v2.DigestSpec{
				HashAlgorithm:          "unknown",
				NormalisationAlgorithm: string(v2.JsonNormalisationV1),
				Value:                  hex.EncodeToString(hashOfString[:]),
			})
			Expect(err).ToNot(BeNil())
			Expect(signature).To(BeNil())
		})

	})
	Describe("RSA Sign verify public key", func() {
		It("should verify a signature", func() {
			hashOfString := sha256.Sum256([]byte(stringToHashAndSign))

			signer, err := signatures.CreateRsaSignerFromKeyFile(pathPrivateKey)
			Expect(err).To(BeNil())
			digest := v2.DigestSpec{
				HashAlgorithm:          "sha256",
				NormalisationAlgorithm: string(v2.JsonNormalisationV1),
				Value:                  hex.EncodeToString(hashOfString[:]),
			}
			signature, err := signer.Sign(v2.ComponentDescriptor{}, digest)
			Expect(err).To(BeNil())
			Expect(signature.Algorithm).To(BeIdenticalTo("RSASSA-PKCS1-V1_5-SIGN"))
			Expect(signature.Value).NotTo(BeNil())

			verifier, err := signatures.CreateRsaVerifierFromKeyFile(pathPublicKey)
			Expect(err).To(BeNil())
			err = verifier.Verify(v2.ComponentDescriptor{}, v2.Signature{
				Digest:    digest,
				Signature: *signature,
			})
			Expect(err).To(BeNil())
		})
		It("should deny a signature from a wrong actor", func() {
			hashOfString := sha256.Sum256([]byte(stringToHashAndSign))

			//generate a wrong key (e.g. from a malicious actor)
			pathWrongPrivateKey := path.Join(dir, "privateWrong.key")
			createWrongPrivateKeyCommand := exec.Command("openssl", "genrsa", "-out", pathWrongPrivateKey, "4096")
			err := createWrongPrivateKeyCommand.Run()
			Expect(err).To(BeNil())

			signer, err := signatures.CreateRsaSignerFromKeyFile(pathWrongPrivateKey)
			Expect(err).To(BeNil())
			digest := v2.DigestSpec{
				HashAlgorithm:          "sha256",
				NormalisationAlgorithm: string(v2.JsonNormalisationV1),
				Value:                  hex.EncodeToString(hashOfString[:]),
			}
			signature, err := signer.Sign(v2.ComponentDescriptor{}, digest)
			Expect(err).To(BeNil())
			Expect(signature.Algorithm).To(BeIdenticalTo("RSASSA-PKCS1-V1_5-SIGN"))
			Expect(signature.Value).NotTo(BeNil())

			verifier, err := signatures.CreateRsaVerifierFromKeyFile(pathPublicKey)
			Expect(err).To(BeNil())
			err = verifier.Verify(v2.ComponentDescriptor{}, v2.Signature{
				Digest:    digest,
				Signature: *signature,
			})
			Expect(err.Error()).To(BeIdenticalTo("signature verification failed, crypto/rsa: verification error"))
		})
	})
})

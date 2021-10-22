package signatures

import (
	"encoding/hex"
	"fmt"

	v2 "github.com/gardener/component-spec/bindings-go/apis/v2"
)

// SignComponentDescriptor signs the given component-descriptor with the signer.
// The component-descriptor has to contain digests for componentReferences and resources.
func SignComponentDescriptor(cd *v2.ComponentDescriptor, signer Signer) error {
	hashCd, err := HashForComponentDescriptor(*cd)
	if err != nil {
		return fmt.Errorf("failed getting hash for cd: %w", err)
	}
	decodedHash, err := hex.DecodeString(hashCd)
	if err != nil {
		return fmt.Errorf("failed decoding hash to bytes")
	}

	signature, err := signer.Sign(*cd, decodedHash)
	if err != nil {
		return fmt.Errorf("failed signing hash of normalised component descriptor, %w", err)
	}
	cd.Signatures = append(cd.Signatures, v2.Signature{NormalisationType: v2.NormalisationTypeV1, Digest: v2.DigestSpec{
		Algorithm: "sha256",
		Value:     hashCd,
	},
		Signature: *signature,
	})
	return nil
}

// VerifySignedComponentDescriptor verifies the signature and hash of the component-descriptor.
// Returns error if verification fails.
func VerifySignedComponentDescriptor(cd *v2.ComponentDescriptor, verifier Verifier) error {
	//Verify hash with signature
	err := verifier.Verify(*cd, cd.Signatures[0]) //TODO: select by name
	if err != nil {
		return fmt.Errorf("failed verifying: %w", err)
	}

	//Verify normalised cd to given (and verified) hash
	hashCd, err := HashForComponentDescriptor(*cd)
	if err != nil {
		return fmt.Errorf("failed getting hash for cd: %w", err)
	}
	if hashCd != cd.Signatures[0].Digest.Value {
		return fmt.Errorf("normalised component-descriptor does not match signed hash")
	}

	return nil
}

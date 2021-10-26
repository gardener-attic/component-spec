package signatures

import (
	"encoding/hex"
	"fmt"

	v2 "github.com/gardener/component-spec/bindings-go/apis/v2"
)

// SignComponentDescriptor signs the given component-descriptor with the signer.
// The component-descriptor has to contain digests for componentReferences and resources.
func SignComponentDescriptor(cd *v2.ComponentDescriptor, signer Signer, hasher Hasher, signatureName string) error {
	hashCd, err := HashForComponentDescriptor(*cd, hasher.HashFunction)
	if err != nil {
		return fmt.Errorf("failed getting hash for cd: %w", err)
	}

	digest := v2.DigestSpec{
		Algorithm: hasher.AlgorithmName,
		Value:     hex.EncodeToString(hashCd),
	}

	signature, err := signer.Sign(*cd, digest)
	if err != nil {
		return fmt.Errorf("failed signing hash of normalised component descriptor, %w", err)
	}
	cd.Signatures = append(cd.Signatures, v2.Signature{
		Name:              signatureName,
		NormalisationType: v2.NormalisationTypeV1,
		Digest:            digest,
		Signature:         *signature,
	})
	return nil
}

// VerifySignedComponentDescriptor verifies the signature (selected by signatureName) and hash of the component-descriptor (as specified in the signature).
// Returns error if verification fails.
func VerifySignedComponentDescriptor(cd *v2.ComponentDescriptor, verifier Verifier, signatureName string) error {
	//find matching signature

	matchingSignature, err := selectSignatureByName(cd, signatureName)
	if err != nil {
		return fmt.Errorf("failed checking signature: %w", err)
	}

	//Verify hash with signature
	err = verifier.Verify(*cd, *matchingSignature)
	if err != nil {
		return fmt.Errorf("failed verifying: %w", err)
	}

	//get hasher by algorithm name
	hasher, err := HasherForName(matchingSignature.Digest.Algorithm)
	if err != nil {
		return fmt.Errorf("failed creating hasher for %s: %w", matchingSignature.Digest.Algorithm, err)
	}

	//Verify normalised cd to given (and verified) hash
	hashCd, err := HashForComponentDescriptor(*cd, hasher.HashFunction)
	if err != nil {
		return fmt.Errorf("failed getting hash for cd: %w", err)
	}
	if hex.EncodeToString(hashCd) != matchingSignature.Digest.Value {
		return fmt.Errorf("normalised component-descriptor does not match signed hash")
	}

	return nil
}

func selectSignatureByName(cd *v2.ComponentDescriptor, signatureName string) (*v2.Signature, error) {
	for _, signature := range cd.Signatures {
		if signature.Name == signatureName {
			return &signature, nil
		}
	}
	return nil, fmt.Errorf("signature with name %s not found in component-descriptor", signatureName)

}

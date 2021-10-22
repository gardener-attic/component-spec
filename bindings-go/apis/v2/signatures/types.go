package signatures

import v2 "github.com/gardener/component-spec/bindings-go/apis/v2"

type Signer interface {
	// Sign returns the signature for the data for the component-descriptor
	Sign(componentDescriptor v2.ComponentDescriptor, data []byte) (*v2.SignatureSpec, error)
}

type Verifier interface {
	// Verify checks the signature, returns an error on verification failure
	Verify(componentDescriptor v2.ComponentDescriptor, signature v2.Signature) error
}

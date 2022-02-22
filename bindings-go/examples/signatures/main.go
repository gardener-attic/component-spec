package main

import (
	"context"
	"fmt"

	v2 "github.com/gardener/component-spec/bindings-go/apis/v2"
	"github.com/gardener/component-spec/bindings-go/apis/v2/signatures"
)

func main() {
	cd := v2.ComponentDescriptor{
		Metadata: v2.Metadata{
			Version: "v2",
		},
		ComponentSpec: v2.ComponentSpec{
			ObjectMeta: v2.ObjectMeta{
				Name:    "CD-Name<html>cool</html> Unicode ♥ unprintable characters \u0007 \u0031",
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
	ctx := context.TODO()
	err := signatures.AddDigestsToComponentDescriptor(ctx, &cd, func(ctx context.Context, cd v2.ComponentDescriptor, cr v2.ComponentReference) (*v2.DigestSpec, error) {
		return &v2.DigestSpec{
			HashAlgorithm:          "testing",
			NormalisationAlgorithm: string(v2.JsonNormalisationV1),
			Value:                  string(cr.GetIdentityDigest()),
		}, nil
	}, func(ctx context.Context, cd v2.ComponentDescriptor, r v2.Resource) (*v2.DigestSpec, error) {
		return &v2.DigestSpec{
			HashAlgorithm:          "testing",
			NormalisationAlgorithm: string(v2.ManifestDigestV1),
			Value:                  string(r.GetIdentityDigest()),
		}, nil
	})
	if err != nil {
		fmt.Printf("ERROR addingDigestsToComponentDescriptor %s", err)
	}

	hasher, err := signatures.HasherForName("sha256")
	if err != nil {
		fmt.Printf("ERROR: %s", err)
	}

	norm, err := signatures.HashForComponentDescriptor(cd, *hasher)
	if err != nil {
		fmt.Printf("ERROR: %s", err)
		return
	}
	fmt.Println(norm.Value)

	signer, err := signatures.CreateRsaSignerFromKeyFile("private")
	if err != nil {
		fmt.Printf("ERROR create signer: %s", err)
		return
	}

	err = signatures.SignComponentDescriptor(&cd, signer, *hasher, "mySignatureName")
	if err != nil {
		fmt.Printf("ERROR sign: %s", err)
		return
	}
	fmt.Println(cd)

	verifier, err := signatures.CreateRsaVerifierFromKeyFile("public")
	if err != nil {
		fmt.Printf("ERROR create verifier: %s", err)
		return
	}
	err = signatures.VerifySignedComponentDescriptor(&cd, verifier, "mySignatureName")
	if err != nil {
		fmt.Printf("ERROR verify signature: %s", err)
		return
	}
	fmt.Println("If not error is printed, successful")
}

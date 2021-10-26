package main

import (
	"encoding/hex"
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
				Name:    "Name",
				Version: "Version",
			},
			ComponentReferences: []v2.ComponentReference{
				v2.ComponentReference{
					Name:          "refName",
					ComponentName: "unimportantName",
					Version:       "v1ref",
					ExtraIdentity: v2.Identity{
						"refKey": "refName",
					},
					Digest: v2.DigestSpec{
						Algorithm: "sha256",
						Value:     "iuhadshksasad",
					},
				},
			},
			Resources: []v2.Resource{
				v2.Resource{
					IdentityObjectMeta: v2.IdentityObjectMeta{
						Name:    "Resource1",
						Version: "v1",
						ExtraIdentity: v2.Identity{
							"key": "value",
						},
					},
				},
				v2.Resource{
					IdentityObjectMeta: v2.IdentityObjectMeta{
						Name:    "Resource2",
						Version: "v2",
						ExtraIdentity: v2.Identity{
							"zzz": "valuezzz",
							"key": "value",
						},
					},
				},
			},
			Sources: []v2.Source{
				v2.Source{
					IdentityObjectMeta: v2.IdentityObjectMeta{
						Name:    "source1",
						Version: "v0.1",
						ExtraIdentity: v2.Identity{
							"key2": "value2",
							"key1": "value1",
						},
					},
				},
			},
		},
		Signatures: []v2.Signature{
			v2.Signature{NormalisationType: v2.NormalisationTypeV1,
				Digest: v2.DigestSpec{
					Algorithm: "sha256",
					Value:     "d782bbae5f6df38c1b7df79319ee6a9625dafcce5df3d730b25aee55db63fcfa",
				},
				Signature: v2.SignatureSpec{
					Algorithm: "RSASSA-PKCS1-V1_5-SIGN",
					Data:      "3d7a81955e1b7cb9556fedb0886229f8b65f1b7ab2cc7be7c1dbe4acff79e2de7415eb1af26402b9d0dd8cdcb90fab212cf122c223bb502900674c90c251b91c044a864b057d1ec4672710ffa79198ab170746100fc2da3b7d78dbfd7a8260bd4adaaf6fccba0d01d7f0a371bdefce9a9792bb1ff4f3e48ad64e33d6a17609a3895d203e17b813ed4f4d4b5f1ef2a803bc2237e8021a92a3bc8780d2617553ece7dfb406af2e3f44000596968476d65b1a4f5e533692c0779823d56fff18c91a9a240a0efd76d012f17f3ecb37152f962590e049274652c634328f70d72ed78e602da182c274f6ded518be560f6846b57cf07227c90b19207b625466a31c0524dcb548433d626200e934a8ae50f987bfe4906fd1db2f1f0337f9ca18793625aedfa4eab2f4fc3b24fdf2c1a6cb2c2f8c4c7082c1b0f037c72d550c0f523349fc6061821d314586939507084ea441437b87152b58b521f1251f4be411e9e96236a385d76968316b4c2eace125417a43730c1bef02db8a0eb8a404352bd8390eda6b4af94681810122917dd59cf11249eb32c8d464f6f34cd7dbd5207efca4275bd20beb06f1afbb112b35980510f50bb4fa958a2426faad74511c48cf1d71abfb6b2bf7291bbb69fcb080a8461b3cd335ca259bd2e54e879d6884c3bf0473922171c018cfb8926ca2d9a0f79fb4618658f22fce7a56d55affd913618124a79e09",
				},
			},
		},
	}

	signatures.AddDigestsToComponentDescriptor(&cd, func(cd v2.ComponentDescriptor, cr v2.ComponentReference) v2.DigestSpec {
		return v2.DigestSpec{
			Algorithm: "testing",
			Value:     string(cr.GetIdentityDigest()),
		}
	}, func(cd v2.ComponentDescriptor, r v2.Resource) v2.DigestSpec {
		return v2.DigestSpec{
			Algorithm: "testing",
			Value:     string(r.GetIdentityDigest()),
		}
	})

	hasher, err := signatures.HasherForName("SHA256")
	if err != nil {
		fmt.Printf("ERROR: %s", err)
	}

	norm, err := signatures.HashForComponentDescriptor(cd, hasher.HashFunction)
	if err != nil {
		fmt.Printf("ERROR: %s", err)
		return
	}
	fmt.Println(hex.EncodeToString(norm))

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

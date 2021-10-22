package signatures

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"

	v2 "github.com/gardener/component-spec/bindings-go/apis/v2"
)

// Entry is used for normalisation and has to contain one key
type Entry map[string]interface{}

// AddDigestsToComponentDescriptor adds digest to componentReferences and resources as returned in the resolver functions
func AddDigestsToComponentDescriptor(cd *v2.ComponentDescriptor, compRefResolver func(v2.ComponentDescriptor, v2.ComponentReference) v2.DigestSpec,
	resResolver func(v2.ComponentDescriptor, v2.Resource) v2.DigestSpec) {

	for i, reference := range cd.ComponentReferences {
		if reference.Digest.Algorithm == "" || reference.Digest.Value == "" {
			cd.ComponentReferences[i].Digest = compRefResolver(*cd, reference)
		}
	}

	for i, res := range cd.Resources {
		if res.Digest.Algorithm == "" || res.Digest.Value == "" {
			cd.Resources[i].Digest = resResolver(*cd, res)
		}
	}
}

// HashForComponentDescriptor return the hash for the component-descriptor, if it is normaliseable
// (= componentReferences and resources contain digest field)
func HashForComponentDescriptor(cd v2.ComponentDescriptor) (string, error) {
	normalisedComponentDescriptor, err := normalizeComponentDescriptor(cd)
	if err != nil {
		return "", fmt.Errorf("failed normalising component descriptor %w", err)
	}
	hash := sha256.Sum256(normalisedComponentDescriptor)

	return hex.EncodeToString(hash[:]), nil
}

func normalizeComponentDescriptor(cd v2.ComponentDescriptor) ([]byte, error) {
	if err := isNormaliseable(cd); err != nil {
		return nil, fmt.Errorf("can not normalise component-descriptor %s:%s: %w", cd.Name, cd.Version, err)
	}

	var normalizedComponentDescriptor []Entry

	meta := []Entry{
		{"schemaVersion": cd.Metadata.Version},
	}
	normalizedComponentDescriptor = append(normalizedComponentDescriptor, Entry{"meta": meta})

	componentReferences := [][]Entry{}
	for _, ref := range cd.ComponentSpec.ComponentReferences {
		extraIdentity := buildExtraIdentity(ref.ExtraIdentity)

		digest := []Entry{
			{"algorithm": ref.Digest.Algorithm},
			{"value": ref.Digest.Value},
		}

		componentReference := []Entry{
			{"name": ref.Name},
			{"version": ref.Version},
			{"extraIdentity": extraIdentity},
			{"digest": digest},
		}
		componentReferences = append(componentReferences, componentReference)
	}

	resources := [][]Entry{}
	for _, res := range cd.ComponentSpec.Resources {
		extraIdentity := buildExtraIdentity(res.ExtraIdentity)

		digest := []Entry{
			{"algorithm": res.Digest.Algorithm},
			{"value": res.Digest.Value},
		}

		resource := []Entry{
			{"name": res.Name},
			{"version": res.Version},
			{"extraIdentity": extraIdentity},
			{"digest": digest},
		}
		resources = append(resources, resource)
	}

	var componentSpec []Entry
	componentSpec = append(componentSpec, Entry{"name": cd.ComponentSpec.Name})
	componentSpec = append(componentSpec, Entry{"version": cd.ComponentSpec.Version})
	componentSpec = append(componentSpec, Entry{"componentReferences": componentReferences})
	componentSpec = append(componentSpec, Entry{"resources": resources})

	normalizedComponentDescriptor = append(normalizedComponentDescriptor, Entry{"component": componentSpec})
	deepSort(normalizedComponentDescriptor)
	normalizedString, err := json.Marshal(normalizedComponentDescriptor)
	if err != nil {
		return nil, err
	}

	return normalizedString, nil
}

func buildExtraIdentity(identity v2.Identity) []Entry {
	var extraIdentities []Entry
	for k, v := range identity {
		extraIdentities = append(extraIdentities, Entry{k: v})
	}
	return extraIdentities
}

// deepSort sorts Entry, []Enry and [][]Entry interfaces recursively, lexicographicly by key(Entry).
func deepSort(in interface{}) {
	switch castIn := in.(type) {
	case []Entry:
		for _, entry := range castIn {
			var val interface{}
			for _, v := range entry {
				val = v
			}
			deepSort(val)

		}
		sort.SliceStable(castIn, func(i, j int) bool {
			var keyI string
			for k := range castIn[i] {
				keyI = k
			}

			var keyJ string
			for k := range castIn[j] {
				keyJ = k
			}

			return keyI < keyJ
		})
	case Entry:
		var val interface{}
		for _, v := range castIn {
			val = v
		}
		deepSort(val)
	case [][]Entry:
		for _, v := range castIn {
			deepSort(v)
		}
	case string:
		break
	default:
		fmt.Println("unknow type")
	}
}

// isNormaliseable checks if componentReferences and resources contain digest
// Does NOT verify if the digest are correct
func isNormaliseable(cd v2.ComponentDescriptor) error {
	// check for digests on component references
	for _, reference := range cd.ComponentReferences {
		if reference.Digest.Algorithm == "" || reference.Digest.Value == "" {
			return fmt.Errorf("missing digest in componentReference for %s:%s", reference.Name, reference.Version)
		}
	}

	// check for digests on resources
	for _, res := range cd.Resources {
		if res.Digest.Algorithm == "" || res.Digest.Value == "" {
			return fmt.Errorf("missing digest in resource for %s:%s", res.Name, res.Version)
		}
	}
	return nil
}

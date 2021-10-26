package signatures

import (
	"encoding/json"
	"fmt"
	"hash"
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
func HashForComponentDescriptor(cd v2.ComponentDescriptor, hashFunction hash.Hash) ([]byte, error) {
	normalisedComponentDescriptor, err := normalizeComponentDescriptor(cd)
	if err != nil {
		return nil, fmt.Errorf("failed normalising component descriptor %w", err)
	}
	hashFunction.Reset()
	_, err = hashFunction.Write(normalisedComponentDescriptor)
	if err != nil {
		return nil, fmt.Errorf("failed hashing the normalisedComponentDescirptorJson: %w", err)
	}
	hash := hashFunction.Sum(nil)
	return hash, nil
}

func normalizeComponentDescriptor(cd v2.ComponentDescriptor) ([]byte, error) {
	if err := isNormaliseable(cd); err != nil {
		return nil, fmt.Errorf("can not normalise component-descriptor %s:%s: %w", cd.Name, cd.Version, err)
	}

	meta := []Entry{
		{"schemaVersion": cd.Metadata.Version},
	}

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

	componentSpec := []Entry{
		{"name": cd.ComponentSpec.Name},
		{"version": cd.ComponentSpec.Version},
		{"componentReferences": componentReferences},
		{"resources": resources},
	}

	normalizedComponentDescriptor := []Entry{
		{"meta": meta},
		{"component": componentSpec},
	}

	deepSort(normalizedComponentDescriptor)

	normalizedJson, err := json.Marshal(normalizedComponentDescriptor)

	if err != nil {
		return nil, err
	}

	return normalizedJson, nil
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
		// sort the values recursively for every entry
		for _, entry := range castIn {
			val := getOnlyValueInEntry(entry)
			deepSort(val)
		}
		// sort the entries based on the key
		sort.SliceStable(castIn, func(i, j int) bool {
			return getOnlyKeyInEntry(castIn[i]) < getOnlyKeyInEntry(castIn[j])
		})
	case Entry:
		val := getOnlyValueInEntry(castIn)
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

func getOnlyKeyInEntry(entry Entry) string {
	var key string
	for k := range entry {
		key = k
	}
	return key
}

func getOnlyValueInEntry(entry Entry) interface{} {
	var value interface{}
	for _, v := range entry {
		value = v
	}
	return value
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

// Copyright 2020 Copyright (c) 2020 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v2

import (
	"encoding/json"
	"fmt"

	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
)

const SchemaVersion = "v2"

// ComponentDescriptor describes the v2 component descriptor containing
// components and their versions.
// Optionally overwrites can be defined.
type ComponentDescriptor struct {
	Metadata Metadata `yaml:"meta"`

	// OverwriteDeclarations contains all overwrites of components
	OverwriteDeclarations []OverwriteDeclaration `json:"overwriteDeclarations"`

	// Components contain all resolvable components with their dependencies
	Components ResolvableComponentList `json:"components"`
}

// Metadata defines the metadata of the component descriptor.
type Metadata struct {
	// Version is the schema version of the component descriptor.
	Version string `json:"schemaVersion"`
}

// ComponentMetadata defines the necessary data of a component
// that uniquely identifies the component.
type ComponentMetadata struct {
	// Scheme defines the specific type of the component
	Type string `json:"type"`
	// Name defines the name of the component.
	Name string `json:"name"`
	// Version defines the semver version of the component.
	// Leading v is optional.
	// Patch number is also optional and will be defaulted to 0.
	Version string `json:"version"`
}

// ComponentAccessor defines the accessor for a component
type ComponentAccessor interface {
	// GetMetadata returns the metadata of the component
	GetMetadata() ComponentMetadata
	// SetMetadata sets the component metadata.
	SetMetadata(metadata ComponentMetadata)
	// GetAdditionalData returns the custom data of a component.
	GetAdditionalData() ([]byte, error)
	// SetAdditionalData sets the custom data of a component.
	SetAdditionalData([]byte) error
	// ApplyOverwrite applies a component overwrite and returns the new overwritten component.
	ApplyOverwrite(overwrite ComponentAccessor) (ComponentAccessor, error)
	// DeepCopy creates a deep copy of the current component
	DeepCopy() ComponentAccessor
}

// ResolvableComponentAccessor defines the accessor for a resolvable component.
type ResolvableComponentAccessor interface {
	// GetMetadata returns the metadata of the component
	GetMetadata() ComponentMetadata
	// SetMetadata sets the component metadata.
	SetMetadata(metadata ComponentMetadata)
	// GetAdditionalData returns the custom data of a component.
	GetAdditionalData() ([]byte, error)
	// SetAdditionalData sets the custom data of a component.
	SetAdditionalData([]byte) error
	// ApplyOverwrite applies a component overwrite and returns the new overwritten component.
	ApplyOverwrite(overwrite ComponentOverwrite) (ResolvableComponentAccessor, error)
	// GetDependencies returns the dependencies of a resolvable component.
	GetDependencies() ComponentList
	// SetDependencies sets the dependencies of a resolvable component.
	SetDependencies(ComponentList)
	// DeepCopy creates a deep copy of the current component
	DeepCopy() ResolvableComponentAccessor
}

type ComponentList []ComponentAccessor

// GetComponentByMetadata returns the component accessor with the given metadata
func (cl *ComponentList) GetComponentByMetadata(meta ComponentMetadata) ComponentAccessor {
	for _, acc := range *cl {
		if ComponentMetadataEquals(acc.GetMetadata(), meta) {
			return acc
		}
	}
	return nil
}

// UnmarshalJSON implements a custom json unmarshal function to parse different types of components
func (cl *ComponentList) UnmarshalJSON(data []byte) error {
	rawComponents := make([]json.RawMessage, 0)
	if err := json.Unmarshal(data, &rawComponents); err != nil {
		return err
	}

	componentList := make(ComponentList, len(rawComponents))
	for i, raw := range rawComponents {
		metadata := ComponentMetadata{}
		if err := json.Unmarshal(raw, &metadata); err != nil {
			return err
		}

		if err := ValidateType(metadata.Type); err != nil {
			return err
		}

		ttype, ok := knownTypes[metadata.Type]
		if !ok {
			// use the generic type for decoding
			comp := &GenericComponent{}
			comp.SetMetadata(metadata)
			if err := comp.SetAdditionalData(raw); err != nil {
				return err
			}
			componentList[i] = comp
			continue
		}

		com, err := ttype.Decode(raw)
		if err != nil {
			return err
		}
		componentList[i] = com
	}

	*cl = componentList
	return nil
}

// MarshalJSON implements a custom json marshal method for a component list
func (cl ComponentList) MarshalJSON() ([]byte, error) {
	rawComponents := make([]json.RawMessage, len(cl))

	for i, component := range cl {
		if err := ValidateType(component.GetMetadata().Type); err != nil {
			return nil, err
		}

		ttype, ok := knownTypes[component.GetMetadata().Type]
		if !ok {
			// encode generic component
			// todo: refactor to own scheme
			data, err := component.GetAdditionalData()
			if err != nil {
				return nil, err
			}
			var rawComponent map[string]interface{}
			if err := json.Unmarshal(data, &rawComponent); err != nil {
				return nil, err
			}

			// metadata should overwrite possible duplicated metadata fields in the generic data
			rawMeta, err := json.Marshal(component.GetMetadata())
			if err != nil {
				return nil, errors.Wrap(err, "unable to marshal component metadata")
			}
			if err := json.Unmarshal(rawMeta, &rawComponent); err != nil {
				return nil, err
			}

			encComp, err := yaml.Marshal(rawComponents)
			if err != nil {
				return nil, err
			}

			rawComponents[i] = encComp
			continue
		}

		encComp, err := ttype.Encode(component)
		if err != nil {
			return nil, err
		}

		rawComponents[i] = encComp
	}

	return json.Marshal(rawComponents)
}

type ResolvableComponentList []ResolvableComponentAccessor

// GetComponentByMetadata returns the component accessor with the given metadata
func (cl *ResolvableComponentList) GetComponentByMetadata(meta ComponentMetadata) ResolvableComponentAccessor {
	for _, acc := range *cl {
		if ComponentMetadataEquals(acc.GetMetadata(), meta) {
			return acc
		}
	}
	return nil
}

// UnmarshalJSON implements a custom json unmarshal function to parse different types of components
func (cl *ResolvableComponentList) UnmarshalJSON(data []byte) error {
	rawComponents := make([]json.RawMessage, 0)
	if err := json.Unmarshal(data, &rawComponents); err != nil {
		return err
	}

	componentList := make(ResolvableComponentList, len(rawComponents))
	for i, raw := range rawComponents {
		metadata := ComponentMetadata{}
		if err := json.Unmarshal(raw, &metadata); err != nil {
			return err
		}

		ttype, ok := knownResolvableTypes[metadata.Type]
		if !ok {
			// generic component are not allowed for resolvable components
			return fmt.Errorf("%s is not a known resolvable component", metadata.Type)
		}

		com, err := ttype.Decode(raw)
		if err != nil {
			return err
		}
		componentList[i] = com
	}

	*cl = componentList
	return nil
}

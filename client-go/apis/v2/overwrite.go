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

import "encoding/json"

// OverwriteDeclaration defines the attribute overwrites for dependencies.
type OverwriteDeclaration struct {
	// DeclaringComponent uniquely identifies the component
	DeclaringComponent ComponentMetadata `json:"declaringComponent"`

	// Overwrites defines the resolvable component overwrites declared
	// by a component
	Overwrites []ComponentOverwrite `json:"overwrites"`
}

// ComponentOverwrite defines an overwrite for a resolvable component.
type ComponentOverwrite struct {
	// ComponentReference references the component to overwrite
	ComponentReference ComponentMetadata `json:"componentReference"`

	// ComponentOverwrites will overwrite overwritable attributes from the
	// component itself.
	ComponentOverwrites json.RawMessage `json:"componentOverwrites"`

	// DependencyOverwrites defines the overwrites for the components dependencies.
	DependencyOverwrites ComponentList `json:"dependencyOverwrites,omitempty"`
}

// ApplyOverwrites applies the overwrites defined by component descriptor
// and returns a new component descriptor with the applied overwrites
func (dc *ComponentDescriptor) ApplyOverwrites() (ComponentDescriptor, error) {
	newDC := ComponentDescriptor{}
	newDC.Metadata = dc.Metadata
	newDC.Components = make(ResolvableComponentList, len(dc.Components))

	for i, comp := range dc.Components {
		overwrites := dc.GetOverwrites(comp.GetMetadata())
		if len(overwrites) == 0 {
			newDC.Components[i] = dc.Components[i]
			continue
		}

		// overwrites are applied in the order they occur in the descriptor
		effectiveComp := comp
		for _, overwrite := range overwrites {
			var err error
			effectiveComp, err = effectiveComp.ApplyOverwrite(overwrite)
			if err != nil {
				return ComponentDescriptor{}, err
			}
		}
		newDC.Components[i] = effectiveComp
	}

	return newDC, nil
}

// GetOverwrites returns the overwrites for component that are defined by all declarations.
func (dc *ComponentDescriptor) GetOverwrites(metadata ComponentMetadata) []ComponentOverwrite {
	overwrites := make([]ComponentOverwrite, 0)
	for _, dec := range dc.OverwriteDeclarations {
		for _, o := range dec.Overwrites {
			if ComponentMetadataEquals(o.ComponentReference, metadata) {
				overwrites = append(overwrites, o)
			}
		}
	}
	return overwrites
}

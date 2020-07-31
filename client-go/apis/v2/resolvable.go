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
	"fmt"
)

// ResolvableComponent describes a top-level component with its dependencies
type ResolvableComponent struct {
	ComponentMetadata `json:",inline"`

	// Dependencies contain all components that the component depends on.
	Dependencies ComponentList `json:"dependencies,omitempty"`
}

// ApplyOverwriteToDependencies applies the given overwrite of the dependencies of a resolvable component
// and returns the new resolvable component.
func ApplyOverwriteToDependencies(rc ResolvableComponentAccessor, overwrite ComponentOverwrite) (ResolvableComponentAccessor, error) {
	deps := rc.GetDependencies()
	newDeps := make(ComponentList, len(rc.GetDependencies()))
	for i, comp := range deps {
		compOverwrite := overwrite.DependencyOverwrites.GetComponentByMetadata(comp.GetMetadata())
		if compOverwrite == nil {
			newDeps[i] = deps[i]
			continue
		}

		newComp, err := comp.ApplyOverwrite(compOverwrite)
		if err != nil {
			return nil, err
		}
		newDeps[i] = newComp
	}
	rc.SetDependencies(newDeps)
	return rc, nil
}

func GetDependenciesByType(comp ResolvableComponentAccessor, ttype string) []ComponentAccessor {
	desc := make([]ComponentAccessor, 0)
	for _, acc := range comp.GetDependencies() {
		if acc.GetMetadata().Type == ttype {
			desc = append(desc, acc)
		}
	}
	return desc
}

func GetOCIImageDependencies(comp ResolvableComponentAccessor) ([]OCIImage, error) {
	desc := make([]OCIImage, 0)
	for _, acc := range comp.GetDependencies() {
		if acc.GetMetadata().Type == OCIImageType {
			// we should be able to cast to the oci type
			img, ok := acc.(*OCIImage)
			if !ok {
				return nil, fmt.Errorf("unabel to cast to internal %s type", OCIImageType)
			}
			desc = append(desc, *img)
		}
	}
	return desc, nil
}

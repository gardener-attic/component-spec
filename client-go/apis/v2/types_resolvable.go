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
	"errors"
	"fmt"

	"github.com/ghodss/yaml"
)

// GardenerComponentType is the type of a gardener component
const GardenerComponentType = "gardenerComponent"

var gardenerComponentScheme = &resolvableScheme{
	ResolvableSchemeDecoder: ResolvableDecoderFunc(func(data []byte) (ResolvableComponentAccessor, error) {
		var comp GardenerComponent
		if err := yaml.Unmarshal(data, &comp); err != nil {
			return nil, err
		}
		return &comp, nil
	}),
	ResolvableSchemeEncoder: ResolvableEncoderFunc(func(accessor ResolvableComponentAccessor) ([]byte, error) {
		comp, ok := accessor.(*GardenerComponent)
		if !ok {
			return nil, fmt.Errorf("accessor is not of type %s", OCIImageType)
		}
		return yaml.Marshal(comp)
	}),
}

// GardenerComponent describes gardener component
type GardenerComponent struct {
	ComponentMetadata `json:",inline"`

	// Dependencies contain all components that the component depends on.
	Dependencies ComponentList `json:"dependencies,omitempty"`
}

var _ ResolvableComponentAccessor = &GardenerComponent{}

func (g GardenerComponent) GetMetadata() ComponentMetadata {
	return g.ComponentMetadata
}

func (g *GardenerComponent) SetMetadata(metadata ComponentMetadata) {
	g.ComponentMetadata = metadata
}

func (g GardenerComponent) GetDependencies() ComponentList {
	return g.Dependencies
}

func (g *GardenerComponent) SetDependencies(list ComponentList) {
	g.Dependencies = list
}

func (g GardenerComponent) GetAdditionalData() ([]byte, error) {
	return nil, nil
}

func (g GardenerComponent) SetAdditionalData(bytes []byte) error {
	return nil
}

func (g *GardenerComponent) ApplyOverwrite(overwrite ComponentOverwrite) (ResolvableComponentAccessor, error) {
	return ApplyOverwriteToDependencies(g.DeepCopy(), overwrite)
}

func (g GardenerComponent) DeepCopy() ResolvableComponentAccessor {
	comp := &GardenerComponent{}
	comp.SetMetadata(g.GetMetadata())
	comp.SetDependencies(g.GetDependencies())
	return comp
}

// GardenerComponentType is the type of a gardener component
const OCIComponentType = "ociComponent"

var ociComponentScheme = &resolvableScheme{
	ResolvableSchemeDecoder: ResolvableDecoderFunc(func(data []byte) (ResolvableComponentAccessor, error) {
		var comp OCIComponent
		if err := yaml.Unmarshal(data, &comp); err != nil {
			return nil, err
		}
		return &comp, nil
	}),
	ResolvableSchemeEncoder: ResolvableEncoderFunc(func(accessor ResolvableComponentAccessor) ([]byte, error) {
		comp, ok := accessor.(*OCIComponent)
		if !ok {
			return nil, fmt.Errorf("accessor is not of type %s", OCIImageType)
		}
		return yaml.Marshal(comp)
	}),
}

// OCIComponent describes oci resolvable component
type OCIComponent struct {
	ComponentMetadata `json:",inline"`

	// Repository is the oci repository where the ociComponent can be accessed
	// the repository is combined with the version to access the specifc component
	// e.g. <repository>:<version>
	Repository string `json:"repository"`

	// Dependencies contain all components that the component depends on.
	Dependencies ComponentList `json:"dependencies,omitempty"`
}

// ComponentOverwrite defines the attributes that can be overwritten in a component overwrite
type OCIComponentOverwrite struct {
	Repository string `json:"repository"`
}

var _ ResolvableComponentAccessor = &OCIComponent{}

func (c OCIComponent) GetMetadata() ComponentMetadata {
	return c.ComponentMetadata
}

func (c *OCIComponent) SetMetadata(metadata ComponentMetadata) {
	c.ComponentMetadata = metadata
}

func (c OCIComponent) GetDependencies() ComponentList {
	return c.Dependencies
}

func (c *OCIComponent) SetDependencies(list ComponentList) {
	c.Dependencies = list
}

func (c OCIComponent) GetAdditionalData() ([]byte, error) {
	return nil, nil
}

func (c OCIComponent) SetAdditionalData(bytes []byte) error {
	return nil
}

func (c *OCIComponent) ApplyOverwrite(overwrite ComponentOverwrite) (ResolvableComponentAccessor, error) {
	effectiveComp, ok := c.DeepCopy().(*OCIComponent)
	if !ok {
		return nil, errors.New("unable to cast ResolvableComponent to OCIComponent")
	}
	if len(overwrite.ComponentOverwrites) != 0 {
		var ociOverwrite OCIComponentOverwrite
		if err := yaml.Unmarshal(overwrite.ComponentOverwrites, &ociOverwrite); err != nil {
			return nil, err
		}

		if len(ociOverwrite.Repository) != 0 {
			effectiveComp.Repository = ociOverwrite.Repository
		}
	}

	return ApplyOverwriteToDependencies(effectiveComp, overwrite)
}

func (c OCIComponent) DeepCopy() ResolvableComponentAccessor {
	comp := &OCIComponent{}
	comp.SetMetadata(c.GetMetadata())
	comp.SetDependencies(c.GetDependencies())
	return comp
}

// SPDX-FileCopyrightText: 2020 SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package v2

import (
	"errors"
)

// ComponentDescriptorList describes the v2 component descriptor containing
// components and their versions.
type ComponentDescriptorList struct {
	Metadata Metadata `json:"meta"`

	// Components contain all resolvable components with their dependencies
	Components []ComponentDescriptor `json:"components"`
}

// GetComponent return the component with a given name and version.
// It returns an error if no component with the name and version is defined.
func (c *ComponentDescriptorList) GetComponent(name, version string) (ComponentDescriptor, error) {
	for _, comp := range c.Components {
		if comp.GetName() == name && comp.GetVersion() == version {
			return comp, nil
		}
	}
	return ComponentDescriptor{}, errors.New("NotFound")
}

// GetComponent returns all components that match the given name.
func (c *ComponentDescriptorList) GetComponentByName(name string) []ComponentDescriptor {
	comps := make([]ComponentDescriptor, 0)
	for _, comp := range c.Components {
		if comp.GetName() == name {
			obj := comp
			comps = append(comps, obj)
		}
	}
	return comps
}

// GetResources returns all resources of a given type, name and version of all components in the list.
func (c *ComponentDescriptorList) GetResources(rtype, name, version string) []Resource {
	resources := make([]Resource, 0)
	for _, comp := range c.Components {
		resource, err := comp.GetResource(rtype, name, version)
		if err == nil {
			resources = append(resources, resource)
		}
	}
	return resources
}

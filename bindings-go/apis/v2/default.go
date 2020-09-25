// SPDX-FileCopyrightText: 2020 SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package v2

// DefaultComponent applies defaults to a component
func DefaultComponent(component *ComponentDescriptor) error {
	if component.Sources == nil {
		component.Sources = make([]Resource, 0)
	}
	if component.ComponentReferences == nil {
		component.ComponentReferences = make([]ComponentReference, 0)
	}
	if component.LocalResources == nil {
		component.LocalResources = make([]Resource, 0)
	}
	if component.ExternalResources == nil {
		component.ExternalResources = make([]Resource, 0)
	}

	for i, res := range component.LocalResources {
		if len(res.Version) == 0 {
			component.LocalResources[i].Version = component.GetVersion()
		}
	}
	return nil
}

func DefaultList(list *ComponentDescriptorList) error {
	for i, comp := range list.Components {
		if len(comp.Metadata.Version) == 0 {
			list.Components[i].Metadata.Version = list.Metadata.Version
		}
	}
	return nil
}

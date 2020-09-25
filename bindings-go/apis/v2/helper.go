// SPDX-FileCopyrightText: 2020 SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package v2

// GetEffectiveRepositoryContext returns the current active repository context.
func (c ComponentDescriptor) GetEffectiveRepositoryContext() RepositoryContext {
	return c.RepositoryContexts[len(c.RepositoryContexts)-1]
}

// GetResource returns a external or local resource with the given type, name and version.
func (c ComponentDescriptor) GetResource(rtype, name, version string) (Resource, error) {
	res, err := c.GetLocalResource(rtype, name, version)
	if err == nil {
		return res, nil
	}

	return c.GetExternalResource(rtype, name, version)
}

// GetExternalResource returns a external resource with the given type, name and version.
func (c ComponentDescriptor) GetExternalResource(rtype, name, version string) (Resource, error) {
	for _, res := range c.ExternalResources {
		if res.GetType() == rtype && res.GetName() == name && res.GetVersion() == version {
			return res, nil
		}
	}
	return Resource{}, NotFound
}

// GetLocalResource returns a local resource with the given type, name and version.
func (c ComponentDescriptor) GetLocalResource(rtype, name, version string) (Resource, error) {
	for _, res := range c.LocalResources {
		if res.GetType() == rtype && res.GetName() == name && res.GetVersion() == version {
			return res, nil
		}
	}
	return Resource{}, NotFound
}

// GetResourcesByType returns all local and external resources of a specific resource type.
func (c ComponentDescriptor) GetResourcesByType(rtype string) []Resource {
	return append(c.GetLocalResourcesByType(rtype), c.GetLocalResourcesByType(rtype)...)
}

// GetLocalResourcesByType returns all local resources of a specific resource type.
func (c ComponentDescriptor) GetLocalResourcesByType(rtype string) []Resource {
	return getResourcesByType(c.LocalResources, rtype)
}

// GetExternalResourcesByType returns all external resources of a specific resource type.
func (c ComponentDescriptor) GetExternalResourcesByType(rtype string) []Resource {
	return getResourcesByType(c.ExternalResources, rtype)
}

func getResourcesByType(list []Resource, rtype string) []Resource {
	resources := make([]Resource, 0)
	for _, obj := range list {
		res := obj
		if res.GetType() == rtype {
			resources = append(resources, res)
		}
	}
	return resources
}

// GetResourcesByType returns all local and external resources of a specific resource type.
func (c ComponentDescriptor) GetResourcesByName(rtype, name string) []Resource {
	return append(c.GetLocalResourcesByName(rtype, name), c.GetExternalResourcesByName(rtype, name)...)
}

// GetLocalResourcesByType returns all local resources of a specific resource type.
func (c ComponentDescriptor) GetLocalResourcesByName(rtype, name string) []Resource {
	return getResourcesByName(c.LocalResources, rtype, name)
}

// GetExternalResourcesByType returns all external resources of a specific resource type.
func (c ComponentDescriptor) GetExternalResourcesByName(rtype, name string) []Resource {
	return getResourcesByName(c.ExternalResources, rtype, name)
}

func getResourcesByName(list []Resource, rtype, name string) []Resource {
	resources := make([]Resource, 0)
	for _, obj := range list {
		res := obj
		if res.GetType() == rtype && res.GetName() == name {
			resources = append(resources, res)
		}
	}
	return resources
}

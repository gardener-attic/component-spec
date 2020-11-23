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
	"bytes"
	"fmt"

	"github.com/gardener/component-spec/bindings-go/utils/selector"
)

type IdentitySelector = selector.Interface

// ResourceSelectorFunc defines a function to filter a resource.
type ResourceSelectorFunc = func(obj Resource) (bool, error)

// MatchResourceSelectorFuncs applies all resource selector against the given resource object.
func MatchResourceSelectorFuncs(obj Resource, resourceSelectors ...ResourceSelectorFunc) (bool, error) {
	for _, sel := range resourceSelectors {
		ok, err := sel(obj)
		if err != nil {
			return false, err
		}
		if !ok {
			return false, nil
		}
	}
	return true, nil
}

// NewTypeResourceSelector creates a new resource selector that
// selects a resource based on its type.
func NewTypeResourceSelector(ttype string) ResourceSelectorFunc {
	return func(obj Resource) (bool, error) {
		return obj.GetType() == ttype, nil
	}
}

// NewVersionResourceSelector creates a new resource selector that
// selects a resource based on its version.
func NewVersionResourceSelector(version string) ResourceSelectorFunc {
	return func(obj Resource) (bool, error) {
		return obj.GetVersion() == version, nil
	}
}

// NewRelationResourceSelector creates a new resource selector that
// selects a resource based on its relation type.
func NewRelationResourceSelector(relation ResourceRelation) ResourceSelectorFunc {
	return func(obj Resource) (bool, error) {
		return obj.Relation == relation, nil
	}
}

// NewNameSelector creates a new selector that matches a resource name.
func NewNameSelector(name string) selector.Interface {
	return selector.DefaultSelector{
		SystemIdentityName: name,
	}
}

// GetEffectiveRepositoryContext returns the current active repository context.
func (c ComponentDescriptor) GetEffectiveRepositoryContext() RepositoryContext {
	return c.RepositoryContexts[len(c.RepositoryContexts)-1]
}

// GetResourceByDefaultSelector returns resources that match the given selectors.
func (c ComponentDescriptor) GetResourceByJSONScheme(sel interface{}) ([]Resource, error) {
	return nil, nil
}

// GetResourceByDefaultSelector returns resources that match the given selectors.
func (c ComponentDescriptor) GetResourceByDefaultSelector(sel interface{}) ([]Resource, error) {
	identitySelector, err := selector.ParseDefaultSelector(sel)
	if err != nil {
		return nil, fmt.Errorf("unable to parse selector: %w", err)
	}
	return c.GetResourceBySelector(identitySelector)
}

// GetResourceBySelector returns resources that match the given selector.
func (c ComponentDescriptor) GetResourceBySelector(selectors ...IdentitySelector) ([]Resource, error) {
	resources := make([]Resource, 0)
	for _, res := range c.Resources {
		ok, err := selector.MatchSelectors(res.GetIdentity(), selectors...)
		if err != nil {
			return nil, fmt.Errorf("unable to match selector for resource %s: %w", res.Name, err)
		}
		if ok {
			resources = append(resources, res)
		}
	}
	if len(resources) == 0 {
		return resources, NotFound
	}
	return resources, nil
}

// GetResourceBySelector returns resources that match the given selector.
func (c ComponentDescriptor) getResourceBySelectors(selectors []IdentitySelector, resourceSelectors []ResourceSelectorFunc) ([]Resource, error) {
	resources := make([]Resource, 0)
	for _, res := range c.Resources {
		ok, err := selector.MatchSelectors(res.GetIdentity(), selectors...)
		if err != nil {
			return nil, fmt.Errorf("unable to match selector for resource %s: %w", res.Name, err)
		}
		if !ok {
			continue
		}
		ok, err = MatchResourceSelectorFuncs(res, resourceSelectors...)
		if err != nil {
			return nil, fmt.Errorf("unable to match selector for resource %s: %w", res.Name, err)
		}
		if !ok {
			continue
		}
		resources = append(resources, res)
	}
	if len(resources) == 0 {
		return resources, NotFound
	}
	return resources, nil
}

// GetResource returns a external or local resource with the given type, name and version.
func (c ComponentDescriptor) GetResource(rtype, name, version string) (Resource, error) {
	for _, res := range c.Resources {
		if res.GetType() == rtype && res.GetName() == name && res.GetVersion() == version {
			return res, nil
		}
	}
	return Resource{}, NotFound
}

// GetExternalResources returns a external resource with the given type, name and version.
func (c ComponentDescriptor) GetExternalResources(rtype, name, version string) ([]Resource, error) {
	return c.getResourceBySelectors(
		[]selector.Interface{NewNameSelector(name)},
		[]ResourceSelectorFunc{
			NewTypeResourceSelector(rtype),
			NewVersionResourceSelector(version),
			NewRelationResourceSelector(ExternalRelation),
		})
}

// GetExternalResource returns a external resource with the given type, name and version.
// If multiple resources match, the first one is returned.
func (c ComponentDescriptor) GetExternalResource(rtype, name, version string) (Resource, error) {
	resources, err := c.GetExternalResources(rtype, name, version)
	if err != nil {
		return Resource{}, err
	}
	// at least one resource must be defined, otherwise the getResourceBySelectors functions returns a NotFound err.
	return resources[0], nil
}

// GetLocalResources returns all local resources with the given type, name and version.
func (c ComponentDescriptor) GetLocalResources(rtype, name, version string) ([]Resource, error) {
	return c.getResourceBySelectors(
		[]selector.Interface{NewNameSelector(name)},
		[]ResourceSelectorFunc{
			NewTypeResourceSelector(rtype),
			NewVersionResourceSelector(version),
			NewRelationResourceSelector(LocalRelation),
		})
}

// GetLocalResource returns a local resource with the given type, name and version.
// If multiple resources match, the first one is returned.
func (c ComponentDescriptor) GetLocalResource(rtype, name, version string) (Resource, error) {
	resources, err := c.GetLocalResources(rtype, name, version)
	if err != nil {
		return Resource{}, err
	}
	// at least one resource must be defined, otherwise the getResourceBySelectors functions returns a NotFound err.
	return resources[0], nil
}

// GetResourcesByType returns all local and external resources of a specific resource type.
func (c ComponentDescriptor) GetResourcesByType(rtype string) []Resource {
	return append(c.GetLocalResourcesByType(rtype), c.GetLocalResourcesByType(rtype)...)
}

// GetLocalResourcesByType returns all local resources of a specific resource type.
func (c ComponentDescriptor) GetLocalResourcesByType(rtype string) []Resource {
	return getResourcesByType(c, rtype)
}

// GetExternalResourcesByType returns all external resources of a specific resource type.
func (c ComponentDescriptor) GetExternalResourcesByType(rtype string) []Resource {
	return getResourcesByType(c, rtype)
}

func getResourcesByType(c ComponentDescriptor, rtype string) []Resource {
	// error can never happen as the type resource selector does not throw an error
	res, _ := c.getResourceBySelectors(
		[]selector.Interface{},
		[]ResourceSelectorFunc{
			NewTypeResourceSelector(rtype),
		})
	return res
}

// GetResourcesByType returns all local and external resources of a specific resource type.
func (c ComponentDescriptor) GetResourcesByName(rtype, name string) ([]Resource, error) {
	return c.getResourceBySelectors(
		[]selector.Interface{NewNameSelector(name)},
		nil)
}

// GetResourceIndex returns the index of a given resource.
// If the index is not found -1 is returned.
func (c ComponentDescriptor) GetResourceIndex(res Resource) int {
	id := res.GetIdentityDigest()
	for i, cur := range c.Resources {
		if bytes.Compare(cur.GetIdentityDigest(), id) == 0 {
			return i
		}
	}
	return -1
}

// ToUnstructuredTypedObject converts a typed object to a unstructured object.
func ToUnstructuredTypedObject(codec TypedObjectCodec, obj TypedObjectAccessor) (*UnstructuredAccessType, error) {
	data, err := codec.Encode(obj)
	if err != nil {
		return nil, err
	}

	uObj := &UnstructuredAccessType{}
	if err := uObj.Decode(data, uObj); err != nil {
		return nil, err
	}
	return uObj, nil
}

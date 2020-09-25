// SPDX-FileCopyrightText: 2020 SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package codec

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/ghodss/yaml"

	"github.com/gardener/component-spec/bindings-go/apis"
	v2 "github.com/gardener/component-spec/bindings-go/apis/v2"
	"github.com/gardener/component-spec/bindings-go/apis/v2/validation"
)

// Decode decodes a component into the given object.
// The obj is expected to be of type v2.ComponentDescriptor or v2.ComponentDescriptorList.
func Decode(data []byte, obj interface{}) error {
	objType := reflect.TypeOf(obj)
	if objType.Kind() != reflect.Ptr {
		return fmt.Errorf("object is expected to be of type pointer but is of type %T", obj)
	}

	raw := make(map[string]json.RawMessage)
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return err
	}

	var metadata apis.Metadata
	if err := yaml.Unmarshal(raw["meta"], &metadata); err != nil {
		return err
	}

	// handle v2
	if metadata.Version == v2.SchemaVersion && objType.Elem() == reflect.TypeOf(v2.ComponentDescriptor{}) {
		// todo: validate against jsonscheme
		if err := yaml.Unmarshal(data, obj); err != nil {
			return err
		}
		comp := obj.(*v2.ComponentDescriptor)
		if err := v2.DefaultComponent(comp); err != nil {
			return err
		}
		return validation.Validate(comp)
	}

	if metadata.Version == v2.SchemaVersion && objType.Elem() == reflect.TypeOf(v2.ComponentDescriptorList{}) {
		if err := yaml.Unmarshal(data, obj); err != nil {
			return err
		}
		list := obj.(*v2.ComponentDescriptorList)
		if err := v2.DefaultList(list); err != nil {
			return err
		}
		return validation.ValidateList(list)
	}

	// todo: implement conversion
	return errors.New("invalid version")
}

// Encode encodes a component or component list into the given object.
// The obj is expected to be of type v2.ComponentDescriptor or v2.ComponentDescriptorList.
func Encode(obj interface{}) ([]byte, error) {
	objType := reflect.TypeOf(obj)
	if objType.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("object is expected to be of type pointer but is of type %T", obj)
	}

	if objType.Elem() == reflect.TypeOf(v2.ComponentDescriptor{}) {
		comp := obj.(*v2.ComponentDescriptor)
		comp.Metadata.Version = v2.SchemaVersion
		if err := v2.DefaultComponent(comp); err != nil {
			return nil, err
		}
		return yaml.Marshal(comp)
	}

	if objType.Elem() == reflect.TypeOf(v2.ComponentDescriptorList{}) {
		list := obj.(*v2.ComponentDescriptorList)
		list.Metadata.Version = v2.SchemaVersion
		if err := v2.DefaultList(list); err != nil {
			return nil, err
		}
		return yaml.Marshal(list)
	}

	// todo: implement conversion
	return nil, errors.New("invalid version")
}

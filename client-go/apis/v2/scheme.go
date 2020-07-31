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
	"strings"
)

// knownTypes contains all known types
var knownTypes = map[string]Scheme{
	OCIImageType: ociScheme,
	WebType:      webScheme,
}

// knownResolvableTypes contains all known resolvable component types
var knownResolvableTypes = map[string]ResolvableScheme{
	GardenerComponentType: gardenerComponentScheme,
	OCIComponentType:      ociComponentScheme,
}

// Scheme describes a known component type and how it is decoded and encoded
type Scheme interface {
	SchemeDecoder
	SchemeEncoder
}

// scheme is a simple struct that implements the Scheme interface
type scheme struct {
	SchemeDecoder
	SchemeEncoder
}

// SchemeDecoder decodes a component dependency.
type SchemeDecoder interface {
	Decode(data []byte) (ComponentAccessor, error)
}

// SchemeEncoder encodes a component dependency.
type SchemeEncoder interface {
	Encode(accessor ComponentAccessor) ([]byte, error)
}

// DecoderFunc is a simple function that implements the SchemeDecoder interface.
type DecoderFunc func(data []byte) (ComponentAccessor, error)

// Decode is the Decode implementation of the SchemeDecoder interface.
func (e DecoderFunc) Decode(data []byte) (ComponentAccessor, error) {
	return e(data)
}

// EncoderFunc is a simple function that implements the SchemeEncoder interface.
type EncoderFunc func(accessor ComponentAccessor) ([]byte, error)

// Encode is the Encode implementation of the SchemeEncoder interface.
func (e EncoderFunc) Encode(accessor ComponentAccessor) ([]byte, error) {
	return e(accessor)
}

// ResolvableScheme describes a known resolvable component type and how it is decoded and encoded
type ResolvableScheme interface {
	ResolvableSchemeDecoder
	ResolvableSchemeEncoder
}

// resolvableScheme is a simple struct that implements the ResolvableScheme interface
type resolvableScheme struct {
	ResolvableSchemeDecoder
	ResolvableSchemeEncoder
}

// SchemeDecoder decodes a resolvable component.
type ResolvableSchemeDecoder interface {
	Decode(data []byte) (ResolvableComponentAccessor, error)
}

// SchemeEncoder encodes a resolvable component.
type ResolvableSchemeEncoder interface {
	Encode(accessor ResolvableComponentAccessor) ([]byte, error)
}

// DecoderFunc is a simple function that implements the ResolvableSchemeDecoder interface.
type ResolvableDecoderFunc func(data []byte) (ResolvableComponentAccessor, error)

// Decode is the Decode implementation of the ResolvableSchemeDecoder interface.
func (e ResolvableDecoderFunc) Decode(data []byte) (ResolvableComponentAccessor, error) {
	return e(data)
}

// ResolvableEncoderFunc is a simple function that implements the SchemeEncoder interface.
type ResolvableEncoderFunc func(accessor ResolvableComponentAccessor) ([]byte, error)

// Encode is the Encode implementation of the SchemeEncoder interface.
func (e ResolvableEncoderFunc) Encode(accessor ResolvableComponentAccessor) ([]byte, error) {
	return e(accessor)
}

// ValidateType validates that a type is known or of a generic type.
// todo: revisit currently "x-" specifies a generic type
func ValidateType(ttype string) error {
	if _, ok := knownTypes[ttype]; ok {
		return nil
	}

	if !strings.HasPrefix(ttype, "x-") {
		return fmt.Errorf("unknown non generic types %s", ttype)
	}
	return nil
}

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

package validation

import (
	"fmt"

	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/validation/field"

	v2 "github.com/gardener/component-spec/bindings-go/apis/v2"
)

// Validate validates a parsed v2 component descriptor
func Validate(cd *v2.ComponentDescriptor) error {
	if cd == nil {
		return nil
	}
	allErrs := field.ErrorList{}

	allErrs = append(allErrs, validateResolvableComponentsList(field.NewPath("components"), cd.Components)...)
	allErrs = append(allErrs, validateOverwriteDeclarations(field.NewPath("overwriteDeclarations"), cd.OverwriteDeclarations)...)

	return allErrs.ToAggregate()
}

func validateOverwriteDeclarations(fldPath *field.Path, declarations []v2.OverwriteDeclaration) field.ErrorList {
	allErrs := field.ErrorList{}
	for i, dec := range declarations {
		iPath := fldPath.Index(i)
		allErrs = append(allErrs, validateComponentMetadata(iPath.Child("declaringComponent"), dec.DeclaringComponent)...)
		allErrs = append(allErrs, validateOverwrite(iPath.Child("overwrites"), dec.Overwrites)...)
	}
	return allErrs
}

func validateOverwrite(fldPath *field.Path, overwrites []v2.ComponentOverwrite) field.ErrorList {
	allErrs := field.ErrorList{}
	for i, ow := range overwrites {
		iPath := fldPath.Index(i)

		allErrs = append(allErrs, validateComponentMetadata(iPath.Child("componentReference"), ow.ComponentReference)...)
		allErrs = append(allErrs, ValidateDependencies(iPath.Child("dependencyOverwrites"), ow.DependencyOverwrites)...)
	}
	return allErrs
}

func validateResolvableComponentsList(fldPath *field.Path, list v2.ResolvableComponentList) field.ErrorList {
	allErrs := field.ErrorList{}
	definedComps := sets.NewString()
	for i, comp := range list {
		cPath := fldPath.Index(i)
		allErrs = append(allErrs, validateComponentMetadata(cPath, comp.GetMetadata())...)

		key := metadataKey(comp.GetMetadata())
		if definedComps.Has(key) {
			allErrs = append(allErrs, field.Duplicate(cPath, fmt.Sprintf("%s: component is already defined", key)))
		}
		definedComps.Insert(key)

		allErrs = append(allErrs, ValidateDependencies(cPath.Child("dependencies"), comp.GetDependencies())...)
	}

	return allErrs
}

// ValidateDependencies validates a list of components.
func ValidateDependencies(fldPath *field.Path, deps v2.ComponentList) field.ErrorList {
	allErrs := field.ErrorList{}
	definedComps := sets.NewString()
	for i, comp := range deps {
		cPath := fldPath.Index(i)
		allErrs = append(allErrs, validateComponentMetadata(cPath, comp.GetMetadata())...)

		key := metadataKey(comp.GetMetadata())
		if definedComps.Has(key) {
			allErrs = append(allErrs, field.Duplicate(cPath, fmt.Sprintf("%s: component is already defined", key)))
		}
		definedComps.Insert(key)
	}
	return allErrs
}

func validateComponentMetadata(fldPath *field.Path, metadata v2.ComponentMetadata) field.ErrorList {
	allErrs := field.ErrorList{}
	if len(metadata.Type) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("name"), "must specify a type"))
	}
	if len(metadata.Name) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("name"), "must specify a name"))
	}
	if len(metadata.Version) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("version"), "must specify a version"))
	}
	return allErrs
}

func metadataKey(metadata v2.ComponentMetadata) string {
	return fmt.Sprintf("%s/%s/%s", metadata.Type, metadata.Name, metadata.Version)
}

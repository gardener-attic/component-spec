// SPDX-FileCopyrightText: 2020 SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package validation

import (
	"k8s.io/apimachinery/pkg/util/validation/field"

	v2 "github.com/gardener/component-spec/bindings-go/apis/v2"
)

// Validate validates a parsed v2 component descriptor
func ValidateList(list *v2.ComponentDescriptorList) error {
	if err := validateList(list); err != nil {
		return err.ToAggregate()
	}
	return nil
}

func validateList(list *v2.ComponentDescriptorList) field.ErrorList {
	if list == nil {
		return nil
	}
	allErrs := field.ErrorList{}

	if len(list.Metadata.Version) == 0 {
		allErrs = append(allErrs, field.Required(field.NewPath("meta").Child("schemaVersion"), "must specify a version"))
	}

	compsPath := field.NewPath("components")
	for i, comp := range list.Components {
		compPath := compsPath.Index(i)
		allErrs = append(allErrs, validate(compPath, &comp)...)
	}

	return allErrs
}

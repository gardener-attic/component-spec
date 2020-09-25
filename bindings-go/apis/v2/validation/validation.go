// SPDX-FileCopyrightText: 2020 SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package validation

import (
	"k8s.io/apimachinery/pkg/util/validation/field"

	v2 "github.com/gardener/component-spec/bindings-go/apis/v2"
)

// Validate validates a parsed v2 component descriptor
func Validate(component *v2.ComponentDescriptor) error {
	if err := validate(nil, component); err != nil {
		return err.ToAggregate()
	}
	return nil
}

func validate(fldPath *field.Path, component *v2.ComponentDescriptor) field.ErrorList {
	if component == nil {
		return nil
	}
	allErrs := field.ErrorList{}

	if len(component.Metadata.Version) == 0 {
		metaPath := field.NewPath("meta").Child("schemaVersion")
		if fldPath != nil {
			metaPath = fldPath.Child("meta").Child("schemaVersion")
		}
		allErrs = append(allErrs, field.Required(metaPath, "must specify a version"))
	}

	compPath := field.NewPath("component")
	if fldPath != nil {
		compPath = fldPath.Child("component")
	}

	if err := validateProvider(compPath.Child("provider"), component.Provider); err != nil {
		allErrs = append(allErrs, err)
	}

	allErrs = append(allErrs, validateObjectMeta(compPath, component.ObjectMeta)...)

	srcPath := compPath.Child("sources")
	allErrs = append(allErrs, validateResources(srcPath, component.Sources)...)

	refPath := compPath.Child("componentReferences")
	allErrs = append(allErrs, validateComponentReferences(refPath, component.ComponentReferences)...)

	localPath := compPath.Child("localResources")
	allErrs = append(allErrs, validateLocalResources(localPath, component.GetVersion(), component.LocalResources)...)

	extPath := compPath.Child("externalResources")
	allErrs = append(allErrs, validateResources(extPath, component.ExternalResources)...)

	return allErrs
}

func validateObjectMeta(fldPath *field.Path, om v2.ObjectMeta) field.ErrorList {
	allErrs := field.ErrorList{}
	if len(om.Name) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("name"), "must specify a name"))
	}
	if len(om.Version) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("version"), "must specify a version"))
	}
	allErrs = append(allErrs, validateLabels(fldPath.Child("labels"), om.Labels)...)
	return allErrs
}

func validateLabels(fldPath *field.Path, labels []v2.Label) field.ErrorList {
	allErrs := field.ErrorList{}
	labelNames := make(map[string]struct{})
	for i, label := range labels {
		labelPath := fldPath.Index(i)
		if len(label.Name) == 0 {
			allErrs = append(allErrs, field.Required(labelPath.Child("name"), "must specify a name"))
			continue
		}

		if _, ok := labelNames[label.Name]; ok {
			allErrs = append(allErrs, field.Duplicate(labelPath, "duplicate label name"))
			continue
		}
		labelNames[label.Name] = struct{}{}
	}
	return allErrs
}

func validateComponentReferences(fldPath *field.Path, refs []v2.ComponentReference) field.ErrorList {
	allErrs := field.ErrorList{}
	refNames := make(map[string]struct{})
	for i, ref := range refs {
		refPath := fldPath.Index(i)
		allErrs = append(allErrs, validateComponentReference(refPath, ref)...)

		if _, ok := refNames[ref.Name]; ok {
			allErrs = append(allErrs, field.Duplicate(refPath, "duplicate component reference name"))
			continue
		}
		refNames[ref.Name] = struct{}{}
	}
	return allErrs
}

func validateComponentReference(fldPath *field.Path, cr v2.ComponentReference) field.ErrorList {
	allErrs := field.ErrorList{}
	if len(cr.ComponentName) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("componentName"), "must specify a component name"))
	}
	allErrs = append(allErrs, validateObjectMeta(fldPath, cr.ObjectMeta)...)
	return allErrs
}

func validateResources(fldPath *field.Path, resources []v2.Resource) field.ErrorList {
	allErrs := field.ErrorList{}
	resourceNames := make(map[string]struct{})
	for i, res := range resources {
		resPath := fldPath.Index(i)
		allErrs = append(allErrs, validateResource(resPath, res)...)

		if _, ok := resourceNames[res.Name]; ok {
			allErrs = append(allErrs, field.Duplicate(resPath, "duplicate resource"))
			continue
		}
		resourceNames[res.Name] = struct{}{}
	}
	return allErrs
}

func validateResource(fldPath *field.Path, res v2.Resource) field.ErrorList {
	allErrs := field.ErrorList{}
	allErrs = append(allErrs, validateObjectMeta(fldPath, res.ObjectMeta)...)

	if len(res.Type) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("type"), "must specify a type"))
	}

	return allErrs
}

func validateLocalResources(fldPath *field.Path, componentVersion string, resources []v2.Resource) field.ErrorList {
	allErrs := field.ErrorList{}
	resourceNames := make(map[string]struct{})
	for i, res := range resources {
		localPath := fldPath.Index(i)
		allErrs = append(allErrs, validateResource(localPath, res)...)
		if res.GetVersion() != componentVersion {
			allErrs = append(allErrs, field.Invalid(localPath.Child("version"), "invalid version",
				"version of local resources must match the component version"))
		}
		if _, ok := resourceNames[res.Name]; ok {
			allErrs = append(allErrs, field.Duplicate(localPath, "duplicated local resource"))
			continue
		}
		resourceNames[res.Name] = struct{}{}
	}
	return allErrs
}

func validateProvider(fldPath *field.Path, provider v2.ProviderType) *field.Error {
	if len(provider) == 0 {
		return field.Required(fldPath, "provider must be set and one of (internal, external)")
	}
	if provider != v2.InternalProvider && provider != v2.ExternalProvider {
		return field.Invalid(fldPath, "unknown provider type", "provider must be one of (internal, external)")
	}
	return nil
}

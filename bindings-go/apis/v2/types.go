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
	"encoding/json"
	"fmt"

	"github.com/ghodss/yaml"
)

// OCIImageType is the type of an oci image component
const OCIImageType = "ociImage"

var ociScheme = &scheme{
	SchemeDecoder: DecoderFunc(func(data []byte) (ComponentAccessor, error) {
		var ociImage OCIImage
		if err := yaml.Unmarshal(data, &ociImage); err != nil {
			return nil, err
		}
		return &ociImage, nil
	}),
	SchemeEncoder: EncoderFunc(func(accessor ComponentAccessor) ([]byte, error) {
		ociImage, ok := accessor.(*OCIImage)
		if !ok {
			return nil, fmt.Errorf("accessor is not of type %s", OCIImageType)
		}
		return yaml.Marshal(ociImage)
	}),
}

// OCIImage is a component describing a OCI container image.
type OCIImage struct {
	ComponentMetadata `json:",inline"`

	// ImageReference is the actual reference to the oci image repository and tag.
	// The format is expected to be "repository:tag".
	ImageReference string `json:"imageReference"`
}

var _ ComponentAccessor = &OCIImage{}

func (O *OCIImage) GetMetadata() ComponentMetadata {
	return O.ComponentMetadata
}

func (O *OCIImage) SetMetadata(metadata ComponentMetadata) {
	O.ComponentMetadata = metadata
}

func (O OCIImage) GetAdditionalData() ([]byte, error) {
	return yaml.Marshal(O)
}

func (O *OCIImage) SetAdditionalData(bytes []byte) error {
	var newOCIImage OCIImage
	if err := yaml.Unmarshal(bytes, &newOCIImage); err != nil {
		return err
	}

	O.ImageReference = newOCIImage.ImageReference
	return nil
}

func (O OCIImage) ApplyOverwrite(overwrite ComponentAccessor) (ComponentAccessor, error) {
	if overwrite.GetMetadata().Type != OCIImageType {
		return nil, fmt.Errorf("unable to overwrite component of type %s with overwrite of type %s", OCIImageType, overwrite.GetMetadata().Type)
	}

	ociOverwrite := overwrite.(*OCIImage)
	O.ImageReference = ociOverwrite.ImageReference
	return &O, nil
}

func (O OCIImage) DeepCopy() ComponentAccessor {
	img := &OCIImage{}
	img.SetMetadata(O.GetMetadata())
	img.ImageReference = O.ImageReference
	return img
}

// WebType is the type of a web component
const WebType = "web"

var webScheme = &scheme{
	SchemeDecoder: DecoderFunc(func(data []byte) (ComponentAccessor, error) {
		var web Web
		if err := yaml.Unmarshal(data, &web); err != nil {
			return nil, err
		}
		return &web, nil
	}),
	SchemeEncoder: EncoderFunc(func(accessor ComponentAccessor) ([]byte, error) {
		web, ok := accessor.(*Web)
		if !ok {
			return nil, fmt.Errorf("accessor is not of type %s", OCIImageType)
		}
		return yaml.Marshal(web)
	}),
}

// Web describes a web resource that can be fetched via http GET request.
type Web struct {
	ComponentMetadata `json:",inline"`

	// URL is the http get accessible url resource.
	URL string `json:"url"`
}

var _ ComponentAccessor = &Web{}

func (w *Web) GetMetadata() ComponentMetadata {
	return w.ComponentMetadata
}

func (w *Web) SetMetadata(metadata ComponentMetadata) {
	w.ComponentMetadata = metadata
}

func (w Web) GetAdditionalData() ([]byte, error) {
	return yaml.Marshal(w)
}

func (w *Web) SetAdditionalData(bytes []byte) error {
	var newWeb Web
	if err := yaml.Unmarshal(bytes, &newWeb); err != nil {
		return err
	}

	w.URL = newWeb.URL
	return nil
}

func (w Web) ApplyOverwrite(overwrite ComponentAccessor) (ComponentAccessor, error) {
	if overwrite.GetMetadata().Type != OCIImageType {
		return nil, fmt.Errorf("unable to overwrite component of type %s with overwrite of type %s", OCIImageType, overwrite.GetMetadata().Type)
	}

	webOverwrite := overwrite.(*Web)
	w.URL = webOverwrite.URL
	return &w, nil
}

func (w Web) DeepCopy() ComponentAccessor {
	web := &Web{}
	web.SetMetadata(w.GetMetadata())
	web.URL = w.URL
	return web
}

// GenericComponent describes a generic dependency of a resolvable component.
type GenericComponent struct {
	ComponentMetadata `json:",inline"`
	AdditionalData    map[string]interface{} `json:"additionalData,omitempty"`
}

var _ ComponentAccessor = &GenericComponent{}

func (c *GenericComponent) GetMetadata() ComponentMetadata {
	return c.ComponentMetadata
}

func (c *GenericComponent) SetMetadata(metadata ComponentMetadata) {
	c.ComponentMetadata = metadata
}

func (c GenericComponent) GetAdditionalData() ([]byte, error) {
	return json.Marshal(c.AdditionalData)
}

func (c *GenericComponent) SetAdditionalData(data []byte) error {
	var values map[string]interface{}
	if err := yaml.Unmarshal(data, &values); err != nil {
		return err
	}
	c.AdditionalData = values
	return nil
}

func (c *GenericComponent) ApplyOverwrite(overwrite ComponentAccessor) (ComponentAccessor, error) {
	newAccessor := &GenericComponent{}
	newAccessor.SetMetadata(overwrite.GetMetadata())
	newData, err := overwrite.GetAdditionalData()
	if err != nil {
		return nil, err
	}
	if err := newAccessor.SetAdditionalData(newData); err != nil {
		return nil, err
	}
	return newAccessor, nil
}

func (c GenericComponent) DeepCopy() ComponentAccessor {
	comp := &GenericComponent{}
	comp.SetMetadata(c.GetMetadata())
	comp.AdditionalData = c.AdditionalData
	return comp
}

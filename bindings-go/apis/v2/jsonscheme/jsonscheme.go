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

//go:generate go-bindata -pkg jsonscheme ../../../../component-descriptor-v2-schema.yaml

package jsonscheme

import (
	"context"
	"errors"
	"fmt"

	"github.com/ghodss/yaml"
	"github.com/qri-io/jsonschema"
)

var Schema *jsonschema.Schema

func init() {
	data, err := ComponentDescriptorV2SchemaYamlBytes()
	if err != nil {
		panic(err)
	}

	Schema = &jsonschema.Schema{}
	if err := yaml.Unmarshal(data, Schema); err != nil {
		panic(err)
	}
}

// Validate validates the given data against the component descriptor v2 jsonscheme.
func Validate(data []byte) error {
	ctx := context.Background()
	defer ctx.Done()
	var doc interface{}
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return err
	}
	state := Schema.Validate(ctx, doc)
	if state == nil {
		return nil
	}

	if state.Errs == nil || len(*state.Errs) == 0 {
		return nil
	}
	errs := *state.Errs
	errMsg := errs[0].Error()
	for i := 1; i < len(errs); i++ {
		errMsg = fmt.Sprintf("%s;%s", errMsg, errs[i].Error())
	}
	return errors.New(errMsg)
}

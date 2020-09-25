// SPDX-FileCopyrightText: 2020 SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

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

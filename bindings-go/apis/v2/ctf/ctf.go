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

package ctf

import (
	"context"
	"errors"
	"io"

	v2 "github.com/gardener/component-spec/bindings-go/apis/v2"
)

// ComponentDescriptorFileName is the name of the component-descriptor file.
const ComponentDescriptorFileName = "component-descriptor.yaml"

// ArtefactDescriptorFileName is the name of the artefact-descriptor file.
const ArtefactDescriptorFileName = "artefact-descriptor.yaml"

// ManifestFileName is the name of the manifest json file.
const ManifestFileName = "manifest.json"

// BlobsDirectoryName is the name of the blob directory in the tar.
const BlobsDirectoryName = "blobs"

var UnsupportedResolveType = errors.New("UnsupportedResolveType")

// BlobResolver defines a resolver that can fetch
// blobs in a specific context defined in a component descriptor.
type BlobResolver interface {
	Info(ctx context.Context, res v2.Resource) (*BlobInfo, error)
	Resolve(ctx context.Context, res v2.Resource, writer io.Writer) (*BlobInfo, error)
}

// BlobInfo describes a blob.
type BlobInfo struct {
	// MediaType is the media type of the object this schema refers to.
	MediaType string `json:"mediaType,omitempty"`

	// Digest is the digest of the targeted content.
	Digest string `json:"digest"`

	// Size specifies the size in bytes of the blob.
	Size int64 `json:"size"`
}

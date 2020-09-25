// SPDX-FileCopyrightText: 2020 SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package apis

// Metadata defines the metadata of the component descriptor.
type Metadata struct {
	// Version is the schema version of the component descriptor.
	Version string `json:"schemaVersion"`
}

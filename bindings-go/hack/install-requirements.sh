#!/usr/bin/env bash

# SPDX-FileCopyrightText: 2020 SAP SE or an SAP affiliate company and Gardener contributors
#
# SPDX-License-Identifier: Apache-2.0

set -e

echo "> Install requirements"

GO111MODULE=off go get golang.org/x/tools/cmd/goimports
GO111MODULE=on curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.27.0
go get -u github.com/go-bindata/go-bindata/...

#!/usr/bin/env bash

# SPDX-FileCopyrightText: 2020 SAP SE or an SAP affiliate company and Gardener contributors
#
# SPDX-License-Identifier: Apache-2.0

set -e

echo "> Check"

echo "Check generated files"
unformatted_files="$(goimports -l --local github.com/gardener/component-spec $@)"
if [[ "$unformatted_files" ]]; then
  echo "Unformatted files:"
  echo "$unformatted_files"
  exit 1
fi

echo "Checks succeeded"
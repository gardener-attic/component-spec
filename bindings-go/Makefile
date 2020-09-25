# SPDX-FileCopyrightText: 2020 SAP SE or an SAP affiliate company and Gardener contributors
#
# SPDX-License-Identifier: Apache-2.0

REPO_ROOT := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

.PHONY: install-requirements
install-requirements:
	@$(REPO_ROOT)/hack/install-requirements.sh

.PHONY: format
format:
	@$(REPO_ROOT)/hack/format.sh $(REPO_ROOT)/apis $(REPO_ROOT)/codec $(REPO_ROOT)/examples

.PHONY: test
test:
	@go test $(REPO_ROOT)/...

.PHONY: check
check:
	@echo "Run lint"; golangci-lint run --timeout 10m $(REPO_ROOT)/...
	@$(REPO_ROOT)/hack/check.sh $(REPO_ROOT)/apis $(REPO_ROOT)/codec $(REPO_ROOT)/examples

.PHONY: revendor
revendor:
	@go mod tidy
	@go mod vendor
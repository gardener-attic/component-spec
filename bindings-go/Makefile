# Copyright (c) 2020 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

REPO_ROOT := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

.PHONY: install-requirements
install-requirements:
	@$(REPO_ROOT)/hack/install-requirements.sh

.PHONY: format
format:
	@$(REPO_ROOT)/hack/format.sh $(REPO_ROOT)/apis $(REPO_ROOT)/codec $(REPO_ROOT)/examples

.PHONY: generate
generate:
	@$(REPO_ROOT)/hack/generate-code.sh

.PHONY: test
test:
	@go test $(REPO_ROOT)/...

.PHONY: check
check:
	@echo "Run lint"; golangci-lint run --timeout 10m $(REPO_ROOT)/...
	@$(REPO_ROOT)/hack/check.sh $(REPO_ROOT)/apis $(REPO_ROOT)/codec $(REPO_ROOT)/examples

.PHONY: verify
verify: check test

.PHONY: revendor
revendor:
	@go mod tidy
	@go mod vendor
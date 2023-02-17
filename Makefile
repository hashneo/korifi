# Setting SHELL to bash allows bash commands to be executed by recipes.
# This is a requirement for 'setup-envtest.sh' in the test target.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk commands is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

CONTROLLERS=controllers job-task-runner kpack-image-builder statefulset-runner
COMPONENTS=api $(CONTROLLERS)

manifests:
	@for comp in $(COMPONENTS); do make -C $$comp manifests; done

generate:
	@for comp in $(CONTROLLERS); do make -C $$comp generate; done
	go run ./scripts/helmdoc/main.go > README.helm.md

generate-fakes:
	go generate ./...

fmt: install-gofumpt install-shfmt
	$(GOFUMPT) -w .
	$(SHFMT) -f . | grep -v '^tests/vendor' | xargs $(SHFMT) -w -i 2 -ci

vet: ## Run go vet against code.
	go vet ./...

lint: fmt vet
	golangci-lint run -v

test: lint
	@for comp in $(COMPONENTS); do make -C $$comp test; done
	make test-e2e

test-e2e: install-ginkgo
	./scripts/run-tests.sh tests/e2e

GOFUMPT = $(shell go env GOPATH)/bin/gofumpt
install-gofumpt:
	go install mvdan.cc/gofumpt@latest

SHFMT = $(shell go env GOPATH)/bin/shfmt
install-shfmt:
	go install mvdan.cc/sh/v3/cmd/shfmt@latest

install-ginkgo:
	go install github.com/onsi/ginkgo/v2/ginkgo

VENDIR = $(shell go env GOPATH)/bin/vendir
install-vendir:
	go install github.com/vmware-tanzu/carvel-vendir/cmd/vendir@latest

vendir-update-dependencies: install-vendir
	$(VENDIR) sync --chdir tests

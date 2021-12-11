PROJECT_NAME:=template-golang
ORG_NAME?=bartoszmajsak
PACKAGE_NAME:=github.com/$(ORG_NAME)/$(PROJECT_NAME)

PROJECT_DIR:=$(shell pwd)
BUILD_DIR:=$(PROJECT_DIR)/build
BINARY_DIR:=$(PROJECT_DIR)/dist
BINARY_NAME:=binary

GOPATH_1:=$(shell echo ${GOPATH} | cut -d':' -f 1)
GOBIN=$(GOPATH_1)/bin
PATH:=${GOBIN}/bin:$(PROJECT_DIR)/bin:$(PATH)

# Call this function with $(call header,"Your message") to see underscored green text
define header =
@echo -e "\n\e[92m\e[4m\e[1m$(1)\e[0m\n"
endef

##@ Default target (all you need - just run "make")
.DEFAULT_GOAL:=all
.PHONY: all
all: deps tools format lint test compile ## Runs 'deps format lint test compile' targets

##@ Build

.PHONY: compile
compile: $(BINARY_DIR)/$(BINARY_NAME) ## Compiles binaries

.PHONY: test
test: ## Runs tests
	$(call header,"Running tests")
	ginkgo -r -v -progress -vet=off -trace --junit-report=ginkgo-test-results.xml ${args}

.PHONY: clean
clean: ## Lets you start from clean state
	rm -rf $(BINARY_DIR) $(PROJECT_DIR)/bin/ vendor/

.PHONY: deps
deps:  ## Fetches all dependencies
	$(call header,"Fetching dependencies")
	@go mod download
	@go mod vendor
	@go mod tidy

.PHONY: format
format: ## Removes unneeded imports and formats source code
	$(call header,"Formatting code")
	goimports -l -w ./pkg/ ./cmd/

.PHONY: lint-prepare
lint-prepare: deps

.PHONY: lint
lint: lint-prepare ## Concurrently runs a whole bunch of static analysis tools
	$(call header,"Running a whole bunch of static analysis tools")
	golangci-lint run --fix --sort-results

# ##########################################################################
# Build configuration
# ##########################################################################

OS:=$(shell uname -s)
GOOS?=$(shell echo $(OS) | awk '{print tolower($$0)}')
GOARCH:=amd64

BUILD_TIME=$(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
GITUNTRACKEDCHANGES:=$(shell git status --porcelain --untracked-files=no)
COMMIT:=$(shell git rev-parse --short HEAD)
ifneq ($(GITUNTRACKEDCHANGES),)
	COMMIT:=$(COMMIT)-dirty
endif

BINARY_VERSION?=$(shell git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")
GIT_TAG:=$(shell git describe --tags --abbrev=0 --exact-match > /dev/null 2>&1; echo $$?)
ifneq ($(GIT_TAG),0)
	BINARY_VERSION:=$(BINARY_VERSION)-next-$(COMMIT)
else ifneq ($(GITUNTRACKEDCHANGES),)
	BINARY_VERSION:=$(BINARY_VERSION)-dirty
endif

GOBUILD:=GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=0
RELEASE?=false
LDFLAGS="-w -X ${PACKAGE_NAME}/version.Release=${RELEASE} -X ${PACKAGE_NAME}/version.Version=${BINARY_VERSION} -X ${PACKAGE_NAME}/version.Commit=${COMMIT} -X ${PACKAGE_NAME}/version.BuildTime=${BUILD_TIME}"
SRCS=$(shell find ./pkg -name "*.go") $(shell find ./cmd -name "*.go") $(shell find ./version -name "*.go")

$(BINARY_DIR):
	[ -d $@ ] || mkdir -p $@

$(BINARY_DIR)/$(BINARY_NAME): $(BINARY_DIR) $(SRCS)
	$(call header,"Compiling... carry on!")
	${GOBUILD} go build -ldflags ${LDFLAGS} -o $@ ./cmd

##@ Setup

.PHONY: tools
tools: $(PROJECT_DIR)/bin/ginkgo $(PROJECT_DIR)/bin/goimports  ## Installs required go tools
tools: $(PROJECT_DIR)/bin/golangci-lint
	$(call header,"Installing required tools")


$(PROJECT_DIR)/bin/ginkgo:
	$(call header,"    Installing ginkgo")
	GOBIN=$(PROJECT_DIR)/bin go install -mod=readonly github.com/onsi/ginkgo/v2/ginkgo

$(PROJECT_DIR)/bin/goimports:
	$(call header,"    Installing goimports")
	GOBIN=$(PROJECT_DIR)/bin go install -mod=readonly golang.org/x/tools/cmd/goimports

$(PROJECT_DIR)/bin/golangci-lint:
	$(call header,"    Installing golangci-lint")
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(PROJECT_DIR)/bin v1.43.0

##@ Helpers

.PHONY: help
help:  ## Displays this help \o/
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-25s\033[0m\033[2m %s\033[0m\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
	@cat $(MAKEFILE_LIST) | grep "^[A-Za-z_]*.?=" | sort | awk 'BEGIN {FS="?="; printf "\n\n\033[1mEnvironment variables\033[0m\n"} {printf "  \033[36m%-25s\033[0m\033[2m %s\033[0m\n", $$1, $$2}'

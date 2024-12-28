# Current Operator version
VERSION ?= $(shell git describe --tags | sed 's/^v//')
VERSION_DATE ?= $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
VERSION_PKG ?= github.com/otel-warez/whitegloves-operator/internal/version
COMMON_LDFLAGS ?= -s -w
OPERATOR_LDFLAGS ?= -X ${VERSION_PKG}.version=${VERSION} -X ${VERSION_PKG}.buildDate=${VERSION_DATE}
ARCH ?= $(shell go env GOARCH)
ifeq ($(shell uname), Darwin)
  SED_INPLACE := sed -i ''
else
  SED_INPLACE := sed -i
endif

# Image URL to use all building/pushing image targets
DOCKER_USER ?= otel-warez
IMG_PREFIX ?= ghcr.io/${DOCKER_USER}/whitegloves-operator
IMG_REPO ?= whitegloves-operator
IMG ?= ${IMG_PREFIX}/${IMG_REPO}:${VERSION}
BUNDLE_IMG ?= ${IMG_PREFIX}/${IMG_REPO}-bundle:${VERSION}

BRIDGETESTSERVER_IMG_REPO ?= e2e-test-app-bridge-server
BRIDGETESTSERVER_IMG ?= ${IMG_PREFIX}/${BRIDGETESTSERVER_IMG_REPO}:ve2e

MUSTGATHER_IMG ?= ${IMG_PREFIX}/must-gather

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# by default, do not run the manager with webhooks enabled. This only affects local runs, not the build or in-cluster deployments.
ENABLE_WEBHOOKS ?= false

# If we are running in CI, run go test in verbose mode
ifeq (,$(CI))
GOTEST_OPTS=-race
else
GOTEST_OPTS=-race -v
endif

START_KIND_CLUSTER ?= true

KUBE_VERSION ?= 1.31
KIND_CONFIG ?= kind-$(KUBE_VERSION).yaml
KIND_CLUSTER_NAME ?= "otel-operator"

OPERATOR_SDK_VERSION ?= 1.29.0

CERTMANAGER_VERSION ?= 1.10.0

ifndef ignore-not-found
  ignore-not-found = false
endif

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

# CHANNELS define the bundle channels used in the bundle.
# Add a new line here if you would like to change its default config. (E.g CHANNELS = "candidate,fast,stable")
# To re-generate a bundle for other specific channels without changing the standard setup, you can:
# - use the CHANNELS as arg of the bundle target (e.g make bundle CHANNELS=candidate,fast,stable)
# - use environment variables to overwrite this value (e.g export CHANNELS="candidate,fast,stable")
ifneq ($(origin CHANNELS), undefined)
BUNDLE_CHANNELS := --channels=$(CHANNELS)
endif

# DEFAULT_CHANNEL defines the default channel used in the bundle.
# Add a new line here if you would like to change its default config. (E.g DEFAULT_CHANNEL = "stable")
# To re-generate a bundle for any other default channel without changing the default setup, you can:
# - use the DEFAULT_CHANNEL as arg of the bundle target (e.g make bundle DEFAULT_CHANNEL=stable)
# - use environment variables to overwrite this value (e.g export DEFAULT_CHANNEL="stable")
ifneq ($(origin DEFAULT_CHANNEL), undefined)
BUNDLE_DEFAULT_CHANNEL := --default-channel=$(DEFAULT_CHANNEL)
endif
BUNDLE_METADATA_OPTS ?= $(BUNDLE_CHANNELS) $(BUNDLE_DEFAULT_CHANNEL)

# BUNDLE_GEN_FLAGS are the flags passed to the operator-sdk generate bundle command
BUNDLE_GEN_FLAGS ?= -q --overwrite --version $(OPERATOR_VERSION) $(BUNDLE_METADATA_OPTS)

# USE_IMAGE_DIGESTS defines if images are resolved via tags or digests
# You can enable this value if you would like to use SHA Based Digests
# To enable set flag to true
USE_IMAGE_DIGESTS ?= false
ifeq ($(USE_IMAGE_DIGESTS), true)
	BUNDLE_GEN_FLAGS += --use-image-digests
endif
MANIFEST_DIR ?= config/crd/bases

# kubectl apply does not work on large CRDs.
CRD_OPTIONS ?= "crd:generateEmbeddedObjectMeta=true,maxDescLen=0"

# Choose wich version to generate
BUNDLE_VARIANT ?= community
BUNDLE_DIR = ./bundle/$(BUNDLE_VARIANT)
MANIFESTS_DIR = config/manifests/$(BUNDLE_VARIANT)
BUNDLE_BUILD_GEN_FLAGS ?= $(BUNDLE_GEN_FLAGS) --output-dir . --kustomize-dir ../../$(MANIFESTS_DIR)

MIN_KUBERNETES_VERSION ?= 1.23.0
MIN_OPENSHIFT_VERSION ?= 4.12

## On MacOS, use gsed instead of sed, to make sed behavior
## consistent with Linux.
SED ?= $(shell which gsed 2>/dev/null || which sed)

.PHONY: ensure-generate-is-noop
ensure-generate-is-noop: VERSION=$(OPERATOR_VERSION)
ensure-generate-is-noop: DOCKER_USER=open-telemetry
ensure-generate-is-noop: set-image-controller generate bundle
	@# on make bundle config/manager/kustomization.yaml includes changes, which should be ignored for the below check
	@git restore config/manager/kustomization.yaml
	@git diff -s --exit-code apis/v1alpha1/zz_generated.*.go || (echo "Build failed: a model has been changed but the generated resources aren't up to date. Run 'make generate' and update your PR." && exit 1)
	@git diff -s --exit-code bundle config || (echo "Build failed: the bundle, config files has been changed but the generated bundle, config files aren't up to date. Run 'make bundle' and update your PR." && git diff && exit 1)
	@git diff -s --exit-code docs/api.md || (echo "Build failed: the api.md file has been changed but the generated api.md file isn't up to date. Run 'make api-docs' and update your PR." && git diff && exit 1)

.PHONY: all
all: manager

# No lint here, as CI runs it separately
.PHONY: ci
ci: generate fmt vet test ensure-generate-is-noop

# Build manager binary
.PHONY: manager
manager: generate
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(ARCH) go build -o bin/manager_${ARCH} -ldflags "${COMMON_LDFLAGS} ${OPERATOR_LDFLAGS}" main.go

.PHONY: must-gather
must-gather:
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(ARCH) go build -o bin/must-gather_${ARCH} -ldflags "${COMMON_LDFLAGS}" ./cmd/gather/main.go

# Run against the configured Kubernetes cluster in ~/.kube/config
.PHONY: run
run: generate fmt vet manifests
	ENABLE_WEBHOOKS=$(ENABLE_WEBHOOKS) go run -ldflags "${OPERATOR_LDFLAGS}" ./main.go --zap-devel

.PHONY: add-operator-arg
add-operator-arg: PATCH = [{"op":"add","path":"/spec/template/spec/containers/0/args/-","value":"$(OPERATOR_ARG)"}]
add-operator-arg: manifests kustomize
	cd config/manager && $(KUSTOMIZE) edit add patch --kind Deployment --patch '$(PATCH)'

.PHONY: add-instrumentation-params
add-instrumentation-params:
	@$(MAKE) add-operator-arg OPERATOR_ARG=--enable-go-instrumentation=true

.PHONY: add-multi-instrumentation-params
add-multi-instrumentation-params:
	@$(MAKE) add-operator-arg OPERATOR_ARG=--enable-multi-instrumentation

# Run go fmt against code
.PHONY: fmt
fmt:
	go fmt ./...

# Run go vet against code
.PHONY: vet
vet:
	go vet ./...

# Run go lint against code
.PHONY: lint
lint: golangci-lint
	$(GOLANGCI_LINT) run

# Generate code
.PHONY: generate
generate: controller-gen

# Build the container image, used only for local dev purposes
# buildx is used to ensure same results for arm based systems (m1/2 chips)
.PHONY: container
container: GOOS = linux
container: manager
	docker build --load -t ${IMG} .
	docker tag ${IMG} wgo:dev

.PHONY: install-metrics-server
install-metrics-server:
	./hack/install-metrics-server.sh

.PHONY: load-image-all
load-image-all: load-image-operator

PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
CMCTL = $(shell pwd)/bin/cmctl
.PHONY: cmctl
cmctl:
	@{ \
	set -e ;\
	if (`pwd`/bin/cmctl version | grep ${CERTMANAGER_VERSION}) > /dev/null 2>&1 ; then \
		exit 0; \
	fi ;\
	TMP_DIR=$$(mktemp -d) ;\
	curl -L -o $$TMP_DIR/cmctl.tar.gz https://github.com/jetstack/cert-manager/releases/download/v$(CERTMANAGER_VERSION)/cmctl-`go env GOOS`-`go env GOARCH`.tar.gz ;\
	tar xzf $$TMP_DIR/cmctl.tar.gz -C $$TMP_DIR ;\
	[ -d bin ] || mkdir bin ;\
	mv $$TMP_DIR/cmctl $(CMCTL) ;\
	rm -rf $$TMP_DIR ;\
	}

KUSTOMIZE ?= $(LOCALBIN)/kustomize
KIND ?= $(LOCALBIN)/kind
CONTROLLER_GEN ?= $(LOCALBIN)/controller-gen
ENVTEST ?= $(LOCALBIN)/setup-envtest
CHLOGGEN ?= $(LOCALBIN)/chloggen
GOLANGCI_LINT ?= $(LOCALBIN)/golangci-lint
CHAINSAW ?= $(LOCALBIN)/chainsaw

# renovate: datasource=go depName=sigs.k8s.io/kustomize/kustomize/v5
KUSTOMIZE_VERSION ?= v5.5.0
# renovate: datasource=go depName=sigs.k8s.io/controller-tools/cmd/controller-gen
CONTROLLER_TOOLS_VERSION ?= v0.16.5
# renovate: datasource=go depName=github.com/golangci/golangci-lint/cmd/golangci-lint
GOLANGCI_LINT_VERSION ?= v1.57.2
# renovate: datasource=go depName=sigs.k8s.io/kind
KIND_VERSION ?= v0.26.0
# renovate: datasource=go depName=github.com/kyverno/chainsaw
CHAINSAW_VERSION ?= v0.2.12

.PHONY: install-tools
install-tools: kustomize golangci-lint kind controller-gen envtest crdoc kind operator-sdk chainsaw

.PHONY: kustomize
kustomize: ## Download kustomize locally if necessary.
	$(call go-get-tool,$(KUSTOMIZE),sigs.k8s.io/kustomize/kustomize/v5,$(KUSTOMIZE_VERSION))

.PHONY: golangci-lint
golangci-lint: ## Download golangci-lint locally if necessary.
	$(call go-get-tool,$(GOLANGCI_LINT),github.com/golangci/golangci-lint/cmd/golangci-lint,$(GOLANGCI_LINT_VERSION))

.PHONY: kind
kind: ## Download kind locally if necessary.
	$(call go-get-tool,$(KIND),sigs.k8s.io/kind,$(KIND_VERSION))

.PHONY: controller-gen
controller-gen: $(CONTROLLER_GEN) ## Download controller-gen locally if necessary.
$(CONTROLLER_GEN): $(LOCALBIN)
	$(call go-get-tool,$(CONTROLLER_GEN), sigs.k8s.io/controller-tools/cmd/controller-gen,$(CONTROLLER_TOOLS_VERSION))

.PHONY: envtest
envtest: $(ENVTEST) ## Download envtest-setup locally if necessary.
$(ENVTEST): $(LOCALBIN)
	$(call go-get-tool,$(ENVTEST), sigs.k8s.io/controller-runtime/tools/setup-envtest,latest)

CRDOC = $(shell pwd)/bin/crdoc
.PHONY: crdoc
crdoc: ## Download crdoc locally if necessary.
	$(call go-get-tool,$(CRDOC), fybrik.io/crdoc,v0.5.2)

.PHONY: chainsaw
chainsaw: ## Find or download chainsaw
	$(call go-get-tool,$(CHAINSAW), github.com/kyverno/chainsaw,$(CHAINSAW_VERSION))

# go-get-tool will 'go get' any package $2 and install it to $1.
PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
define go-get-tool
@[ -f $(1) ] || { \
set -e ;\
TMP_DIR=$$(mktemp -d) ;\
cd $$TMP_DIR ;\
go mod init tmp ;\
echo "Downloading $(2)" ;\
go get -d $(2)@$(3) ;\
GOBIN=$(PROJECT_DIR)/bin go install $(2) ;\
rm -rf $$TMP_DIR ;\
}
endef

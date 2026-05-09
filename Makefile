## HILIOS Operator Makefile.
## Targets are grouped by category; run `make help` for the full list.

IMG ?= ghcr.io/kitio-tek/hilios-operator:dev
ENVTEST_K8S_VERSION ?= 1.30.0
CONTROLLER_TOOLS_VERSION ?= v0.18.0
ENVTEST_VERSION ?= release-0.18
GOLANGCI_LINT_VERSION ?= v2.12.2
KUSTOMIZE_VERSION ?= v5.4.2
PLATFORMS ?= linux/amd64,linux/arm64

GOBIN ?= $(shell go env GOPATH)/bin
LOCALBIN ?= $(shell pwd)/bin
KUBECTL ?= kubectl
KUSTOMIZE ?= $(LOCALBIN)/kustomize
CONTROLLER_GEN ?= $(LOCALBIN)/controller-gen
ENVTEST ?= $(LOCALBIN)/setup-envtest
GOLANGCI_LINT ?= $(LOCALBIN)/golangci-lint
CONTAINER_TOOL ?= docker

SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

.PHONY: all
all: build

##@ General

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Code generation

.PHONY: manifests
manifests: controller-gen ## Generate CRDs, RBAC and webhook manifests.
	$(CONTROLLER_GEN) rbac:roleName=manager-role crd webhook paths="./..." output:crd:artifacts:config=config/crd/bases

.PHONY: generate
generate: controller-gen ## Generate DeepCopy methods.
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./..."

.PHONY: chart-sync
chart-sync: manifests ## Sync generated CRDs into the Helm chart templates/crds directory.
	cp config/crd/bases/*.yaml charts/hilios-operator/templates/crds/

##@ Verification

.PHONY: fmt
fmt: ## Run gofmt.
	go fmt ./...

.PHONY: verify-licenses
verify-licenses: ## Fail when a Go source file is missing the license header.
	bash hack/check-licenses.sh

.PHONY: verify-yaml
verify-yaml: ## Fail when YAML files cannot be parsed.
	bash hack/check-yaml.sh

.PHONY: verify-fmt
verify-fmt: ## Fail when go files are not gofmt'd.
	@if [ -n "$$(gofmt -l . | grep -v '^vendor/')" ]; then \
		echo "gofmt produced changes:"; \
		gofmt -d $$(gofmt -l . | grep -v '^vendor/'); \
		exit 1; \
	fi

.PHONY: vet
vet: ## Run go vet.
	go vet ./...

.PHONY: lint
lint: golangci-lint ## Run golangci-lint.
	$(GOLANGCI_LINT) run ./...

.PHONY: lint-fix
lint-fix: golangci-lint ## Run golangci-lint with --fix.
	$(GOLANGCI_LINT) run --fix ./...

.PHONY: helm-lint
helm-lint: ## Lint the Helm chart.
	helm lint charts/hilios-operator/

.PHONY: helm-template
helm-template: ## Render the Helm chart for inspection.
	helm template hilios charts/hilios-operator/ --namespace hilios-system

##@ Tests

.PHONY: test
test: manifests generate fmt vet ## Run unit tests.
	go test -race -coverprofile=cover.out ./...

.PHONY: test-envtest
test-envtest: envtest manifests generate ## Run controller integration tests against envtest.
	KUBEBUILDER_ASSETS="$(shell $(ENVTEST) use $(ENVTEST_K8S_VERSION) --bin-dir $(LOCALBIN) -p path)" \
		go test ./internal/controller/... -coverprofile=cover-envtest.out

.PHONY: test-e2e
test-e2e: ## Run KUTTL end-to-end tests against the cluster pointed at by KUBECONFIG.
	$(KUBECTL) kuttl test tests/e2e/kuttl/ --config tests/e2e/kuttl/kuttl-test.yaml

##@ Build

.PHONY: build
build: manifests generate fmt vet ## Build the manager binary.
	go build -trimpath -ldflags "-s -w" -o bin/manager ./cmd/main.go

.PHONY: run
run: manifests generate fmt vet ## Run the manager from the host.
	go run ./cmd/main.go

.PHONY: docker-build
docker-build: ## Build the manager container image.
	$(CONTAINER_TOOL) build -t ${IMG} .

.PHONY: docker-push
docker-push: ## Push the manager container image.
	$(CONTAINER_TOOL) push ${IMG}

.PHONY: docker-buildx
docker-buildx: ## Build and push the multi-arch manager image.
	$(CONTAINER_TOOL) buildx build --push --platform=$(PLATFORMS) --tag $(IMG) -f Dockerfile .

##@ Deployment

.PHONY: install-crds
install-crds: manifests kustomize ## Install CRDs into the cluster pointed at by KUBECONFIG.
	$(KUSTOMIZE) build config/crd | $(KUBECTL) apply -f -

.PHONY: uninstall-crds
uninstall-crds: manifests kustomize ## Uninstall CRDs from the cluster pointed at by KUBECONFIG.
	$(KUSTOMIZE) build config/crd | $(KUBECTL) delete --ignore-not-found=true -f -

.PHONY: deploy
deploy: ## Install the operator into the cluster via Helm.
	helm upgrade --install hilios charts/hilios-operator/ \
		--namespace hilios-system --create-namespace \
		--set image.repository=$(word 1,$(subst :, ,$(IMG))) \
		--set image.tag=$(word 2,$(subst :, ,$(IMG)))

.PHONY: undeploy
undeploy: ## Remove the operator release from the cluster.
	-helm uninstall hilios -n hilios-system

##@ Tools

$(LOCALBIN):
	mkdir -p $(LOCALBIN)

.PHONY: controller-gen
controller-gen: $(CONTROLLER_GEN) ## Install controller-gen into ./bin if missing.
$(CONTROLLER_GEN): $(LOCALBIN)
	$(call go-install-tool,$(CONTROLLER_GEN),sigs.k8s.io/controller-tools/cmd/controller-gen,$(CONTROLLER_TOOLS_VERSION))

.PHONY: kustomize
kustomize: $(KUSTOMIZE) ## Install kustomize into ./bin if missing.
$(KUSTOMIZE): $(LOCALBIN)
	$(call go-install-tool,$(KUSTOMIZE),sigs.k8s.io/kustomize/kustomize/v5,$(KUSTOMIZE_VERSION))

.PHONY: envtest
envtest: $(ENVTEST) ## Install setup-envtest into ./bin if missing.
$(ENVTEST): $(LOCALBIN)
	$(call go-install-tool,$(ENVTEST),sigs.k8s.io/controller-runtime/tools/setup-envtest,$(ENVTEST_VERSION))

.PHONY: golangci-lint
golangci-lint: $(GOLANGCI_LINT) ## Install golangci-lint into ./bin if missing.
$(GOLANGCI_LINT): $(LOCALBIN)
	$(call go-install-tool,$(GOLANGCI_LINT),github.com/golangci/golangci-lint/cmd/golangci-lint,$(GOLANGCI_LINT_VERSION))

# go-install-tool installs a versioned binary into LOCALBIN under a name suffixed
# with the version, then symlinks the requested name to it. This avoids
# accidentally re-using a stale binary when the version variable changes.
define go-install-tool
@[ -f "$(1)-$(3)" ] || { \
set -e; \
package=$(2)@$(3) ;\
echo "Downloading $${package}" ;\
rm -f $(1) || true ;\
GOBIN=$(LOCALBIN) go install $${package} ;\
mv $(1) $(1)-$(3) ;\
} ;\
ln -sf $(1)-$(3) $(1)
endef

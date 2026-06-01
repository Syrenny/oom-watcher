##@ General
PACKAGE_NAME := oom-watcher
VERSION ?= $(shell cat VERSION)
DEB_DIR := deb
DEB_OUTPUT := $(DEB_DIR)/$(PACKAGE_NAME)_$(VERSION)_amd64.deb
RELEASE_DIR := dist
RELEASE_ASSET := $(RELEASE_DIR)/$(PACKAGE_NAME)_amd64.deb
BIN_DIR := bin
BINARY := $(BIN_DIR)/$(PACKAGE_NAME)
BUILD_DIR := .build
DEB_STAGE_DIR := $(BUILD_DIR)/deb
GO ?= $(shell command -v go 2>/dev/null || echo /usr/local/go/bin/go)
GOFMT ?= $(shell command -v gofmt 2>/dev/null || echo /usr/local/go/bin/gofmt)
GOLANGCI_LINT ?= $(shell command -v golangci-lint 2>/dev/null)

.PHONY: help all build deb release-deb lint fmt clean

help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

all: build deb ## Build binary and versioned .deb package

##@ Build
build: fmt lint ## Build binary
	@mkdir -p $(BIN_DIR)
	$(GO) build -ldflags="-X github.com/Syrenny/oom-watcher/internal/version.Version=$(VERSION)" -o $(BINARY) ./cmd/$(PACKAGE_NAME)

##@ Debian package
deb: build ## Build versioned .deb package
	rm -rf $(DEB_STAGE_DIR)
	mkdir -p $(BUILD_DIR)
	cp -a $(DEB_DIR) $(DEB_STAGE_DIR)
	chmod 755 $(DEB_STAGE_DIR)/DEBIAN/postinst $(DEB_STAGE_DIR)/DEBIAN/postrm
	find $(DEB_STAGE_DIR)/etc -type d -exec chmod 755 {} +
	find $(DEB_STAGE_DIR)/etc -type f -exec chmod 644 {} +
	find $(DEB_STAGE_DIR)/usr -type d -exec chmod 755 {} +
	find $(DEB_STAGE_DIR)/usr -type f -exec chmod 755 {} +
	sed -i "s/^Version:.*/Version: $(VERSION)/" $(DEB_STAGE_DIR)/DEBIAN/control
	rm -f $(DEB_STAGE_DIR)/usr/local/bin/$(PACKAGE_NAME)
	cp $(BINARY) $(DEB_STAGE_DIR)/usr/local/bin/
	dpkg-deb --root-owner-group --build $(DEB_STAGE_DIR) $(DEB_OUTPUT)
	@echo "Built package: $(DEB_OUTPUT)"

release-deb: deb ## Build stable release asset for GitHub Releases
	mkdir -p $(RELEASE_DIR)
	cp $(DEB_OUTPUT) $(RELEASE_ASSET)
	@echo "Built release asset: $(RELEASE_ASSET)"

##@ Lint & fmt
lint: ## Run golangci-lint
ifneq ($(strip $(GOLANGCI_LINT)),)
	$(GOLANGCI_LINT) run ./...
else
	@echo "golangci-lint not found; skipping lint"
endif

fmt: ## Format Go code
	$(GOFMT) -w $$(find . -name '*.go' -not -path './vendor/*')

##@ Clean
clean: ## Remove build artifacts
	rm -rf $(BIN_DIR)
	rm -rf $(BUILD_DIR)
	rm -rf $(RELEASE_DIR)
	rm -f $(DEB_OUTPUT)

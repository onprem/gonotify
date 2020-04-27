PROJECTNA   ?= $(shell basename "$(PWD)")
BASE         = $(shell pwd)
BUILD_DIR   ?= $(BASE)/build
VETARGS     ?= -all
GOFMT_FILES ?= $$(find . -name '*.go' | grep -v vendor)

# Ensure GOPATH is set
GOPATH            ?= $(shell go env GOPATH)

TMP_GOPATH        ?= /tmp/gonotify-go
GOBIN             ?= $(firstword $(subst :, ,${GOPATH}))/bin
GO111MODULE       ?= on
export GO111MODULE
GOPROXY           ?= https://proxy.golang.org
export GOPROXY

# Tools
GOBINDATA         ?= $(GOBIN)/go-bindata


.PHONY: help
help: ## Display usage and help message.
help:
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:\n"} /^[a-z0-9A-Z_-]+:.*?##/ { printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)


$(GOBINDATA):
	@echo ">> installing go-bindata"
	@GO111MODULE='off' go get -u github.com/go-bindata/go-bindata/...


.PHONY: run
run: ## Runs the program.
run:
	@echo ">> starting $(PROJECTNAME)"
	@go run cmd/gonotify/main.go


.PHONY: run-prod
run-prod: ## runs the program in release mode.
run-prod:
	@echo ">> starting $(PROJECTNAME) in release mode"
	@GIN_MODE=release go run cmd/gonotify/main.go


.PHONY: build
build: ## Builds the GoNotify binary.
build:
	@echo ">> building $(PROJECTNAME) binary"
	@go build -o $(BUILD_DIR)/gonotify ./cmd/gonotify


.PHONY: assets
assets: ## Repacks all static assets into go file for easier deploy.
assets: $(GOBINDATA)
	@echo ">> deleting asset file"
	@rm pkg/ui/bindata.go || true
	@echo ">> writing assets"
	@go-bindata -pkg ui -o pkg/ui/bindata.go -prefix "/pkg/ui/react-app/build" pkg/ui/react-app/build/...
	@go fmt ./pkg/ui


.PHONY: vet
vet: ## Runs go vet against all packages.
vet:
	@echo ">> running go vet on packages"
	@go vet $(VETARGS) ./pkg/... ./cmd/... ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi


.PHONY: fmt
fmt: ## Format all go files using go fmt.
fmt:
	@echo ">> running go fmt on all go files"
	@gofmt -w $(GOFMT_FILES)

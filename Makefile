PROJECTNAME   ?= $(shell basename "$(PWD)")
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

# React
REACT_APP_PATH = pkg/ui/react-app
REACT_APP_SOURCE_FILES = $(wildcard $(REACT_APP_PATH)/public/* $(REACT_APP_PATH)/src/* $(REACT_APP_PATH)/package.json)
REACT_APP_OUTPUT_DIR = $(REACT_APP_PATH)/build
REACT_APP_NODE_MODULES_PATH = $(REACT_APP_PATH)/node_modules


.PHONY: help
help: ## Display usage and help message.
help:
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:\n"} /^[a-z0-9A-Z_-]+:.*?##/ { printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)


$(GOBINDATA):
	@echo ">> installing go-bindata"
	@GO111MODULE='off' go get -u github.com/go-bindata/go-bindata/...


$(REACT_APP_NODE_MODULES_PATH): $(REACT_APP_PATH)/package.json $(REACT_APP_PATH)/yarn.lock
	@echo ">> installing npm dependencies for React UI"
	@cd $(REACT_APP_PATH) && yarn --frozen-lockfile

$(REACT_APP_OUTPUT_DIR): $(REACT_APP_NODE_MODULES_PATH) $(REACT_APP_SOURCE_FILES)
	@echo ">> building React app"
	@cd $(REACT_APP_PATH) && yarn build


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
build: assets
	@echo ">> building $(PROJECTNAME) binary"
	@go build -o $(BUILD_DIR)/gonotify ./cmd/gonotify

.PHONY: build-static
build-static: ## Builds a statically linked binary for easy deployment
build-static: assets
	@echo ">> building statically linked $(PROJECTNAME) binary"
	@CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o $(BUILD_DIR)/gonotify ./cmd/gonotify

.PHONY: build-cli
build-cli: ## Builds gncli - the GoNotify CLI client.
build-cli:
	@echo ">> building gncli binary"
	@go build -o $(BUILD_DIR)/gncli ./cmd/gncli

.PHONY: assets
assets: ## Repacks all static assets into go file for easier deploy.
assets: $(GOBINDATA) $(REACT_APP_OUTPUT_DIR)
	@echo ">> deleting asset file"
	@rm pkg/ui/bindata.go || true
	@echo ">> writing assets"
	@go-bindata -pkg ui -o pkg/ui/bindata.go -prefix "pkg/ui/react-app/build" pkg/ui/react-app/build/...
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


.PHONY: react-app-start
react-app-start: ## Start React app for local development
react-app-start: $(REACT_APP_NODE_MODULES_PATH)
	@echo ">> running React app"
	@cd $(REACT_APP_PATH) && yarn start

.PHONY: react-app-lint
react-app-lint: $(REACT_APP_NODE_MODULES_PATH)
	   @echo ">> running React app linting"
	   cd $(REACT_APP_PATH) && yarn lint:ci

.PHONY: react-app-lint-fix
react-app-lint-fix:
	@echo ">> running React app linting and fixing errors where possible"
	cd $(REACT_APP_PATH) && yarn lint

.PHONY: react-app-test
react-app-test: | $(REACT_APP_NODE_MODULES_PATH) react-app-lint
	@echo ">> running React app tests"
	cd $(REACT_APP_PATH) && export CI=true && yarn test --no-watch --coverage

.PHONY: build test clean changelog tools help

S3HUB       = s3hub
VERSION     = $(shell git describe --tags --abbrev=0)
GO          = go
GO_BUILD    = $(GO) build
GO_INSTALL  = $(GO) install
GO_TEST     = $(GO) test -v
GO_TOOL     = $(GO) tool
GO_DEP      = $(GO) mod
GOOS        = ""
GOARCH      = ""
GO_PKGROOT  = ./...
GO_PACKAGES = $(shell $(GO_LIST) $(GO_PKGROOT))
GO_LDFLAGS  = -ldflags '-X github.com/nao1215/rainbow/version.Version=${VERSION}'

build:  ## Build binary
	env GO111MODULE=on GOOS=$(GOOS) GOARCH=$(GOARCH) $(GO_BUILD) $(GO_LDFLAGS) -o $(S3HUB) cmd/s3hub/main.go

clean: ## Clean project
	-rm -rf $(S3HUB) cover.out cover.html

test: ## Start test
	env GOOS=$(GOOS) $(GO_TEST) -cover $(GO_PKGROOT) -coverprofile=cover.out
	$(GO_TOOL) cover -html=cover.out -o cover.html

changelog: ## Generate changelog
	ghch --format markdown > CHANGELOG.md

tools: ## Install dependency tools 
	$(GO_INSTALL) github.com/Songmu/ghch/cmd/ghch@latest

.DEFAULT_GOAL := help
help:  
	@grep -E '^[0-9a-zA-Z_-]+[[:blank:]]*:.*?## .*$$' $(MAKEFILE_LIST) | sort \
	| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[1;32m%-15s\033[0m %s\n", $$1, $$2}'
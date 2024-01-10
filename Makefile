.PHONY: build test clean changelog tools help docker generate gif coverage-tree

S3HUB       = s3hub
SPARE       = spare
CFN         = cfn
VERSION     = $(shell git describe --tags --abbrev=0)
GO          = go
GO_BUILD    = $(GO) build
GO_INSTALL  = $(GO) install
GO_TEST     = hottest -v
GO_TOOL     = $(GO) tool
GO_DEP      = $(GO) mod
GOOS        = ""
GOARCH      = ""
GO_PKGROOT  = ./...
GO_PACKAGES = $(shell $(GO_LIST) $(GO_PKGROOT))
GO_LDFLAGS  = -ldflags '-X github.com/nao1215/rainbow/version.Version=${VERSION}'

build:  ## Build binary
	env GO111MODULE=on GOOS=$(GOOS) GOARCH=$(GOARCH) $(GO_BUILD) $(GO_LDFLAGS) -o $(S3HUB) cmd/s3hub/main.go
	env GO111MODULE=on GOOS=$(GOOS) GOARCH=$(GOARCH) $(GO_BUILD) $(GO_LDFLAGS) -o $(SPARE) cmd/spare/main.go
	env GO111MODULE=on GOOS=$(GOOS) GOARCH=$(GOARCH) $(GO_BUILD) $(GO_LDFLAGS) -o $(CFN) cmd/cfn/main.go

clean: ## Clean project
	-rm -rf $(S3HUB) $(SPARE) $(CFN) cover.out cover.html

test: ## Start unit test
	env GOOS=$(GOOS) $(GO_TEST) -coverpkg=./... -coverprofile=cover.out.tmp -cover ./...
	cat cover.out.tmp | grep -v "_gen.go" | grep -v "main.go" > cover.out
	$(GO_TOOL) cover -html=cover.out -o cover.html

coverage-tree: test ## Generate coverage tree
	go-cover-treemap -statements -coverprofile cover.out > doc/img/cover.svg

changelog: ## Generate changelog
	ghch --format markdown > CHANGELOG.md

tools: ## Install dependency tools 
	$(GO_INSTALL) github.com/Songmu/ghch/cmd/ghch@latest
	$(GO_INSTALL) github.com/nao1215/hottest@latest
	$(GO_INSTALL) github.com/google/wire/cmd/wire@latest
	$(GO_INSTALL) github.com/charmbracelet/vhs@latest
	$(GO_INSTALL) github.com/nikolaydubina/go-cover-treemap@latest

generate: ## Generate code from templates
	$(GO) generate ./...

gif: docker ## Generate gif image
	vhs < doc/img/vhs/s3hub-mb.tape
	vhs < doc/img/vhs/s3hub-ls.tape
	vhs < doc/img/vhs/s3hub-ls.tape
	vhs < doc/img/vhs/s3hub-rm-all.tape

docker:  ## Start docker (localstack)
	docker compose up -d

.DEFAULT_GOAL := help
help:  
	@grep -E '^[0-9a-zA-Z_-]+[[:blank:]]*:.*?## .*$$' $(MAKEFILE_LIST) | sort \
	| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[1;32m%-15s\033[0m %s\n", $$1, $$2}'
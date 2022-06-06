INSTALL_DIR ?= ./atomci
GO_TAGS ?=
CGO_ENABLED ?= 0
GO_PACKAGES ?= ./...
GO_BUILD = go build -v -tags="${GO_TAGS}"
NAME = atomci

export CGO_ENABLED := ${CGO_ENABLED}
export PROJECT_ROOT := $(CURDIR)


VERSION := v1.5.1
BuildTS = $(shell date -u '+%Y-%m-%d %I:%M:%S')
COMMIT_ID=$(shell git rev-parse --short HEAD)
BRANCH_NAME=$(shell git rev-parse --abbrev-ref HEAD)

project=$(shell go list -m)
LDFLAGS += -X "$(project)/version.BuildTS=$(BuildTS)"
LDFLAGS += -X "$(project)/version.GitHash=$(COMMIT_ID)"
LDFLAGS += -X "$(project)/version.Version=$(VERSION)"
LDFLAGS += -X "$(project)/version.GitBranch=$(BRANCH_NAME)"

## linux-amd64: Compile linux-amd64 package
linux-amd64:
	@env GOOS=linux GOARCH=amd64 go build -o deploy/binary/$(NAME)-linux-amd64 cmd/atomci/main.go

## linux-arm64: Compile linux-amd64 package
linux-arm64:
	@env GOOS=linux GOARCH=arm64 go build -o deploy/binary/$(NAME)-linux-arm64 cmd/atomci/main.go

.PHONY: build
## build: Compile the packages.
build:
	@go build -ldflags '$(LDFLAGS)' -o $(NAME) cmd/atomci/main.go

.PHONY: run
## run: Build and Run in local mode.
run: build
	@./$(NAME)

.PHONY: web
web:
	cd web; yarn run dev
	

.PHONY: lint
lint:
	go vet -structtag=false -tags="${GO_TAGS}" ${GO_PACKAGES}


# NOTE: cgo is required by go test
.PHONY: test
test:
	CGO_ENABLED=1 go test -v -race -cover -failfast -vet=off -tags="${GO_TAGS}" ${GO_PACKAGES}

.PHONY: clean
clean:
	rm -f $(NAME)
	go clean -cache ${GO_PACKAGES}


.PHONY: help
all: help
# help: show this help message
help: Makefile
	@echo
	@echo " Choose a command to run in "$(NAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
 
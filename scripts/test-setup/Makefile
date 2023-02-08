SHELL := /bin/bash 
GO111MODULE=on

# Build ldflags
VERSION ?= "v0.0.0"
GITCOMMIT=$(shell git rev-parse HEAD)
BUILDDATE=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
PKG_PATH=github.com/aws-controllers-k8s/test-infra/scripts/test-setup
GO_LDFLAGS=-ldflags "-X $(PKG_PATH)/version.GitVersion=$(VERSION) \
			-X $(PKG_PATH)/version.GitCommit=$(GITCOMMIT) \
			-X $(PKG_PATH)/version.BuildDate=$(BUILDDATE)"

GO_CMD_FLAGS=-gcflags="all=-N -l" 
GO_CMD_LOCAL_FLAGS=-modfile=go.local.mod $(GO_CMD_FLAGS)

.PHONY: build run test

build:
	CGO_ENABLED=0 GOOS=$(shell go env GOOS) GOARCH=$(shell go env GOARCH) go build ${GO_CMD_FLAGS} ${GO_LDFLAGS} -o ./bin/test-setup ./main.go

run: build
	./bin/test-setup setup --kubeconfig ${KUBECONFIG} --config $(shell go env GOPATH)/src/github.com/aws-controllers-k8s/test-infra/test_config.yaml

test:
	go test -tags $(shell go env GOOS) -v ./...
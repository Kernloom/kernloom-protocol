GO ?= go
GOBIN ?= $(HOME)/go/bin
PATH := $(GOBIN):$(PATH)
BUF ?= buf
BUF_CACHE_DIR ?= /tmp/kernloom-bufcache

.PHONY: fmt vet test build generate lint proto-check

fmt:
	$(GO) fmt ./...

vet:
	$(GO) vet ./...

test:
	$(GO) test ./...

build:
	$(GO) test ./...

generate:
	BUF_CACHE_DIR=$(BUF_CACHE_DIR) $(BUF) generate
	sh scripts/add-license-headers.sh sdk/go/adapter/v1/adapter.pb.go sdk/go/adapter/v1/adapter_grpc.pb.go
	$(GO) fmt ./sdk/go/adapter/v1

lint:
	BUF_CACHE_DIR=$(BUF_CACHE_DIR) $(BUF) lint

proto-check:
	$(GO) test ./internal/protoschema

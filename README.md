# kernloom-protocol

`kernloom-protocol` contains Kernloom v1 protocol definitions, the initial Go adapter SDK and the adapter contract test kit.

## Build

```sh
make build
```

## Generate

```sh
make generate
```

This uses `buf`, `protoc-gen-go` and `protoc-gen-go-grpc`. The Makefile prepends `$(HOME)/go/bin` to `PATH`, so locally installed Go plugins are discovered without changing your shell profile.

## Test

```sh
make test
```

## Release

Protocol releases must version the protobuf surface, Go SDK and contract tests together. Breaking-change checks are required before a tagged release.

## Dependencies

Slice 0 uses Go 1.26.4, Buf, `protoc-gen-go`, `protoc-gen-go-grpc` and `github.com/bufbuild/protocompile`. Generated Go and gRPC stubs are committed under `sdk/go/adapter/v1`.

## Related Repos

Adapters import the SDK and contract tests from this repo. `kernloom-core` consumes protocol descriptors through out-of-process gRPC adapters.

GO ?= go
GOBIN ?= $(HOME)/go/bin
PATH := $(GOBIN):$(PATH)
BUF ?= buf
BUF_CACHE_DIR ?= /tmp/kernloom-bufcache
TRIVY ?= trivy
COSIGN ?= cosign
GOVULNCHECK ?= govulncheck
DIST ?= dist

.PHONY: fmt vet test build generate lint proto-check checksums sbom vuln-scan govulncheck release-provenance release-sign release-promote-check release-check

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

checksums:
	mkdir -p $(DIST)
	find proto sdk contract -type f | sort | xargs sha256sum > $(DIST)/checksums.txt

release-provenance: checksums
	{ \
		echo "{"; \
		echo "  \"kind\": \"KernloomProtocolReleaseProvenance\","; \
		echo "  \"source_commit\": \"$$(git rev-parse HEAD)\","; \
		echo "  \"go_version\": \"$$($(GO) version)\","; \
		echo "  \"checksums\": \"$(DIST)/checksums.txt\""; \
		echo "}"; \
	} > $(DIST)/provenance.json

sbom:
	@command -v $(TRIVY) >/dev/null 2>&1 || { echo "trivy is required for SBOM generation"; exit 127; }
	mkdir -p $(DIST)
	$(TRIVY) fs --format cyclonedx --output $(DIST)/sbom.cdx.json .

vuln-scan:
	@command -v $(TRIVY) >/dev/null 2>&1 || { echo "trivy is required for vulnerability scanning"; exit 127; }
	$(TRIVY) fs --exit-code 1 --severity HIGH,CRITICAL .

govulncheck:
	@command -v $(GOVULNCHECK) >/dev/null 2>&1 || { echo "govulncheck is required"; exit 127; }
	$(GOVULNCHECK) ./...

release-sign: checksums
	@command -v $(COSIGN) >/dev/null 2>&1 || { echo "cosign is required for release signing"; exit 127; }
	$(COSIGN) sign-blob --yes --output-signature $(DIST)/checksums.txt.sig $(DIST)/checksums.txt

release-promote-check: checksums sbom release-provenance
	test -s $(DIST)/checksums.txt
	test -s $(DIST)/sbom.cdx.json
	test -s $(DIST)/provenance.json

release-check: generate lint fmt vet test proto-check checksums sbom vuln-scan govulncheck release-provenance release-promote-check

.PHONY: check gen gen-check build release clean-dist

GO := GOWORK=off go

BIN     := country-map-svg-generator
PKG     := ./cmd/$(BIN)
DIST    := dist
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)

# buildVersion is the only stamped identity. The corpus and algorithm versions
# are read from the embedded artifacts, so a binary cannot claim a corpus it does
# not carry. -trimpath keeps the build reproducible across checkout paths, which
# REQ-8's byte-identical determinism is measured against.
LDFLAGS := -s -w -X main.buildVersion=$(VERSION)
BUILDFLAGS := -trimpath -ldflags '$(LDFLAGS)'

# Every target is pure Go with an embedded corpus: no CGO, no network, no
# auxiliary files at runtime, so cross-compilation is a plain GOOS/GOARCH switch
# and the result runs on a clean machine.
PLATFORMS := darwin/arm64 darwin/amd64 linux/arm64 linux/amd64

# -p 1 is load-bearing, not tidiness. Several tests assert wall-clock brakes —
# "this build must not take longer than 30s", "this generation must not take
# longer than 10s" — which exist to catch a runaway search, the failure mode that
# produced twenty interrupted runs. Run in parallel with the full-catalog ladder
# gates those brakes measure CPU contention between packages instead of the code
# under test, and flake: lodbuild's post-build micro-brake was observed at 34.5s
# alongside the 381s geometry package and at well under 30s on its own.
#
# The cost is that check now takes roughly the sum of the packages rather than
# the slowest one. The alternative — raising the thresholds until they stop
# flaking — would leave the brakes shaped like guards while guarding nothing.
# Restoring parallelism means expressing the brakes in work done (candidate
# evaluations, phase attempts) rather than seconds; tracked in the backlog.
check:
	$(GO) test ./... -count=1 -p 1

# One binary for this machine, stamped the same way a release is.
build:
	CGO_ENABLED=0 $(GO) build $(BUILDFLAGS) -o $(DIST)/$(BIN) $(PKG)
	@$(DIST)/$(BIN) version

# Cross-compiled archives plus checksums. VERSION defaults to git describe, so
# a tagged commit releases itself: make release, or make release VERSION=v0.1.0.
release: clean-dist
	@set -eu; for platform in $(PLATFORMS); do \
		os=$${platform%/*}; arch=$${platform#*/}; \
		stage="$(DIST)/$(BIN)_$(VERSION)_$${os}_$${arch}"; \
		mkdir -p "$$stage"; \
		echo "building $$os/$$arch"; \
		CGO_ENABLED=0 GOOS=$$os GOARCH=$$arch $(GO) build $(BUILDFLAGS) -o "$$stage/$(BIN)" $(PKG); \
		cp README.md "$$stage/"; \
		tar -czf "$$stage.tar.gz" -C "$(DIST)" "$$(basename "$$stage")"; \
		rm -rf "$$stage"; \
	done; \
	cd $(DIST) && shasum -a 256 *.tar.gz > SHA256SUMS
	@echo; ls -lh $(DIST)

clean-dist:
	@rm -rf $(DIST)

gen:
	$(GO) run ./internal/catalog/cmd/corpuscompile -data data -out data/corpus/v1 -report data/reports/v1/compare.json -embedded internal/catalog/embedded_gen.go

gen-check:
	@tmp=$$(mktemp -d); \
	trap 'rm -rf "$$tmp"' EXIT; \
	$(GO) run ./internal/catalog/cmd/corpuscompile -data data -out "$$tmp/corpus" -report "$$tmp/compare.json" -embedded "$$tmp/embedded_gen.go"; \
	diff -ru data/corpus/v1 "$$tmp/corpus"; \
	diff -u data/reports/v1/compare.json "$$tmp/compare.json"; \
	diff -u internal/catalog/embedded_gen.go "$$tmp/embedded_gen.go"

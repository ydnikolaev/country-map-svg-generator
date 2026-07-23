.PHONY: check gen gen-check

GO := GOWORK=off go

check:
	$(GO) test ./... -count=1

gen:
	$(GO) run ./internal/catalog/cmd/corpuscompile -data data -out data/corpus/v1 -report data/reports/v1/compare.json -embedded internal/catalog/embedded_gen.go

gen-check:
	@tmp=$$(mktemp -d); \
	trap 'rm -rf "$$tmp"' EXIT; \
	$(GO) run ./internal/catalog/cmd/corpuscompile -data data -out "$$tmp/corpus" -report "$$tmp/compare.json" -embedded "$$tmp/embedded_gen.go"; \
	diff -ru data/corpus/v1 "$$tmp/corpus"; \
	diff -u data/reports/v1/compare.json "$$tmp/compare.json"; \
	diff -u internal/catalog/embedded_gen.go "$$tmp/embedded_gen.go"

.PHONY: check gen gen-check

GO := GOWORK=off go

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

gen:
	$(GO) run ./internal/catalog/cmd/corpuscompile -data data -out data/corpus/v1 -report data/reports/v1/compare.json -embedded internal/catalog/embedded_gen.go

gen-check:
	@tmp=$$(mktemp -d); \
	trap 'rm -rf "$$tmp"' EXIT; \
	$(GO) run ./internal/catalog/cmd/corpuscompile -data data -out "$$tmp/corpus" -report "$$tmp/compare.json" -embedded "$$tmp/embedded_gen.go"; \
	diff -ru data/corpus/v1 "$$tmp/corpus"; \
	diff -u data/reports/v1/compare.json "$$tmp/compare.json"; \
	diff -u internal/catalog/embedded_gen.go "$$tmp/embedded_gen.go"

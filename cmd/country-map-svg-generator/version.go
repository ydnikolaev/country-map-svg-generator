package main

import (
	"fmt"
	"runtime/debug"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yuranikolaev/country-map-svg-generator/internal/catalog"
	"github.com/yuranikolaev/country-map-svg-generator/internal/geometry"
)

// buildVersion is overridden at link time by the release build. It stays "dev"
// in a source build so a test can assert an exact string; the corpus and
// algorithm identities below are the ones that actually determine output, and
// they are read from the embedded artifacts rather than stamped.
var buildVersion = "dev"

// VersionInfo is CTR-004's version payload. Every field here is an input to
// reproducibility: ARCH-INV-4 says equal inputs, corpus version and generator
// version produce byte-identical output, so a caller comparing two runs needs
// all three identities from one command.
type VersionInfo struct {
	Version          string `json:"version"`
	Schema           string `json:"schema"`
	CorpusVersion    int    `json:"corpus_version"`
	CorpusIdentity   string `json:"corpus_identity"`
	CorpusEntities   int    `json:"corpus_entities"`
	GeometrySchema   int    `json:"geometry_schema_version"`
	AlgorithmVersion string `json:"algorithm_version"`
	GoVersion        string `json:"go_version"`
}

func newVersionCommand(flags *globalFlags) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Report generator, corpus and algorithm identities",
		Long: strings.TrimSpace(`
Report every identity that determines output: the binary version, the embedded
corpus version and content identity, and the geometry algorithm version.

Two runs that agree on all of these produce byte-identical assets.`),
		Args: exactArgs(0, CLIName+" version [--json]"),
		RunE: func(cmd *cobra.Command, args []string) error {
			info, err := versionInfo()
			if err != nil {
				return err
			}
			return emit(cmd, flags, info, humanVersion(info))
		},
	}
}

func versionInfo() (VersionInfo, error) {
	corpus, err := catalog.Embedded()
	if err != nil {
		return VersionInfo{}, failf(ExitData, "corpus_unreadable", "embedded corpus could not be decoded: %v", err)
	}
	info := VersionInfo{
		Version:          buildVersion,
		Schema:           EnvelopeSchema,
		CorpusVersion:    corpus.Manifest.CorpusVersion,
		CorpusIdentity:   corpus.Manifest.Identity,
		CorpusEntities:   corpus.Manifest.EntityCount,
		GeometrySchema:   geometry.SchemaVersion,
		AlgorithmVersion: geometry.AlgorithmVersion,
	}
	if build, ok := debug.ReadBuildInfo(); ok {
		info.GoVersion = build.GoVersion
	}
	return info, nil
}

func humanVersion(info VersionInfo) string {
	var out strings.Builder
	fmt.Fprintf(&out, "%s %s\n", CLIName, info.Version)
	fmt.Fprintf(&out, "corpus     v%d (%s, %d entities)\n", info.CorpusVersion, short(info.CorpusIdentity), info.CorpusEntities)
	fmt.Fprintf(&out, "geometry   schema %d, %s\n", info.GeometrySchema, info.AlgorithmVersion)
	if info.GoVersion != "" {
		fmt.Fprintf(&out, "built with %s", info.GoVersion)
	}
	return out.String()
}

// short truncates a content digest for human output only. The JSON envelope
// always carries the full identity — a digest a caller might compare must never
// be the truncated one.
//
// The truncation counts hex characters, not bytes of the whole string: the
// corpus identity is stored as "sha256:<64 hex>", and taking the first twelve
// bytes of that leaves five hex digits, which is short enough to collide and to
// be useless for telling two corpus versions apart at a glance.
func short(digest string) string {
	prefix, hex, found := strings.Cut(digest, ":")
	if !found {
		prefix, hex = "", digest
	}
	if len(hex) > 12 {
		hex = hex[:12]
	}
	if prefix == "" {
		return hex
	}
	return prefix + ":" + hex
}

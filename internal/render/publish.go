package render

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// This file is INV-1: "a successful batch is deterministic and atomically
// published; a failed batch leaves the previous successful output unchanged."
//
// The whole batch is built and validated in memory first, and only then does
// anything touch the output directory. That ordering is the point: a batch that
// wrote as it went would leave a half-regenerated catalog behind on the first
// failure, and the manifest would agree with it, so nothing downstream could
// tell the difference between a complete catalog and a truncated one.
//
// The replace-with-backup-and-rollback shape follows the corpus publisher in
// internal/catalog, which solved the same problem for P1. It is re-implemented
// rather than shared because those helpers are unexported and the boundary
// between BND-001 and BND-006 is worth more than the duplication.

// ManifestVersion is CTR-006's schema version. The manifest is consumed by site
// and build pipelines, so its shape is a public compatibility surface.
const ManifestVersion = "country-map.manifest/v1"

// ManifestName is the file the manifest is always written to, beside the assets
// it describes.
const ManifestName = "country-map.manifest.json"

// Asset is one generated file as CTR-006 records it.
type Asset struct {
	ISO      string `json:"iso"`
	Name     string `json:"name"`
	File     string `json:"file"`
	Profile  string `json:"profile"`
	Boundary string `json:"boundary"`
	Style    string `json:"style"`
	Delivery string `json:"delivery"`
	Bytes    int    `json:"bytes"`
	Width    string `json:"width"`
	Height   string `json:"height"`
	// Markers is CTR-006's pin state: how many pins the asset actually carries,
	// which is what a consumer checks rather than what was requested.
	Markers int    `json:"markers"`
	SHA256  string `json:"sha256"`
	// PathBytes is the geometry's contribution, separate from the file total, so
	// a budget conversation can distinguish the silhouette from its wrapper.
	PathBytes int `json:"path_bytes"`
}

// Skipped is an entity that produced no asset for a reason that is a decision
// rather than a failure — today only DEC-009's typed no-artifact outcome.
//
// It is recorded rather than omitted: a consumer comparing the manifest against
// the corpus needs to tell "this entity was refused for a stated reason" from
// "this entity is missing and nobody noticed".
type Skipped struct {
	ISO    string `json:"iso"`
	Name   string `json:"name"`
	Reason string `json:"reason"`
}

// Manifest is CTR-006: one machine-readable index of a generation.
type Manifest struct {
	Schema string `json:"schema"`
	// Generator, Corpus and Config are the three identities ARCH-INV-4 names.
	// Two runs that agree on all three produce byte-identical output, so a
	// consumer comparing manifests can tell a real change from a re-run.
	Generator      string    `json:"generator_version"`
	CorpusVersion  int       `json:"corpus_version"`
	CorpusIdentity string    `json:"corpus_identity"`
	ConfigDigest   string    `json:"config_digest"`
	Assets         []Asset   `json:"assets"`
	Skipped        []Skipped `json:"skipped"`
	TotalBytes     int       `json:"total_bytes"`
}

// StagedAsset is one asset held in memory before publication.
type StagedAsset struct {
	Asset
	Content []byte
}

// Batch is a complete, validated generation waiting to be published.
type Batch struct {
	Assets   []StagedAsset
	Skipped  []Skipped
	Manifest Manifest
}

// Finalize sorts a batch and builds its manifest. Sorting is what makes the
// manifest byte-identical across runs: map iteration order would otherwise vary
// and ARCH-INV-4 would hold for every asset but not for the index of them.
func (b *Batch) Finalize(generator string, corpusVersion int, corpusIdentity, configDigest string) {
	sort.Slice(b.Assets, func(i, j int) bool { return b.Assets[i].ISO < b.Assets[j].ISO })
	sort.Slice(b.Skipped, func(i, j int) bool { return b.Skipped[i].ISO < b.Skipped[j].ISO })

	manifest := Manifest{
		Schema: ManifestVersion, Generator: generator,
		CorpusVersion: corpusVersion, CorpusIdentity: corpusIdentity, ConfigDigest: configDigest,
		Assets: make([]Asset, 0, len(b.Assets)), Skipped: b.Skipped,
	}
	if manifest.Skipped == nil {
		manifest.Skipped = []Skipped{}
	}
	for _, staged := range b.Assets {
		manifest.Assets = append(manifest.Assets, staged.Asset)
		manifest.TotalBytes += staged.Bytes
	}
	b.Manifest = manifest
}

// ManifestBytes renders the manifest. Indented and newline-terminated because it
// is a file a human diffs as often as a pipeline reads it.
func (b *Batch) ManifestBytes() ([]byte, error) {
	encoded, err := json.MarshalIndent(b.Manifest, "", "  ")
	if err != nil {
		return nil, err
	}
	return append(encoded, '\n'), nil
}

// PublishError names what went wrong and, crucially, whether the previous output
// survived. A caller that cannot tell those apart has to assume the worst and
// regenerate from scratch.
type PublishError struct {
	Err            error
	PreviousIntact bool
}

func (e *PublishError) Error() string {
	if e.PreviousIntact {
		return e.Err.Error() + " (the previous output is unchanged)"
	}
	return e.Err.Error() + " (the output directory may be inconsistent; regenerate)"
}

func (e *PublishError) Unwrap() error { return e.Err }

// Publish writes a finalized batch to the output directory transactionally.
//
// Every file is staged beside its destination first, so the rename that commits
// it is within one filesystem and cannot fail half-way for the usual reason. If
// any rename fails, everything already replaced is rolled back from its backup.
func Publish(batch *Batch, outputDir string) error {
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return &PublishError{Err: err, PreviousIntact: true}
	}
	manifestBytes, err := batch.ManifestBytes()
	if err != nil {
		return &PublishError{Err: err, PreviousIntact: true}
	}

	staging, err := os.MkdirTemp(outputDir, ".country-map-staging-")
	if err != nil {
		return &PublishError{Err: err, PreviousIntact: true}
	}
	defer os.RemoveAll(staging)

	type pending struct{ target, staged, backup string }
	items := make([]pending, 0, len(batch.Assets)+1)

	stage := func(relative string, content []byte) error {
		target := filepath.Join(outputDir, relative)
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		staged := filepath.Join(staging, strings.ReplaceAll(relative, string(filepath.Separator), "__"))
		if err := os.WriteFile(staged, content, 0o644); err != nil {
			return err
		}
		items = append(items, pending{target: target, staged: staged})
		return nil
	}

	for _, asset := range batch.Assets {
		if err := stage(asset.File, asset.Content); err != nil {
			// Nothing has been replaced yet, so the previous output is untouched.
			return &PublishError{Err: err, PreviousIntact: true}
		}
	}
	// The manifest is staged last and published last, so a consumer that reads it
	// never sees an index describing files that are not there yet.
	if err := stage(ManifestName, manifestBytes); err != nil {
		return &PublishError{Err: err, PreviousIntact: true}
	}

	rollback := func(last int) {
		for i := last; i >= 0; i-- {
			_ = os.RemoveAll(items[i].target)
			if items[i].backup != "" {
				_ = os.Rename(items[i].backup, items[i].target)
			}
		}
	}

	for i := range items {
		if _, err := os.Lstat(items[i].target); err == nil {
			backup, backupErr := reserveBackup(items[i].target)
			if backupErr != nil {
				rollback(i - 1)
				return &PublishError{Err: backupErr, PreviousIntact: true}
			}
			if err := os.Rename(items[i].target, backup); err != nil {
				rollback(i - 1)
				return &PublishError{Err: err, PreviousIntact: true}
			}
			items[i].backup = backup
		} else if !errors.Is(err, os.ErrNotExist) {
			rollback(i - 1)
			return &PublishError{Err: err, PreviousIntact: true}
		}
		if err := os.Rename(items[i].staged, items[i].target); err != nil {
			if items[i].backup != "" {
				_ = os.Rename(items[i].backup, items[i].target)
				items[i].backup = ""
			}
			rollback(i - 1)
			return &PublishError{Err: err, PreviousIntact: true}
		}
	}
	for _, item := range items {
		if item.backup != "" {
			_ = os.RemoveAll(item.backup)
		}
	}
	return nil
}

// reserveBackup claims a unique name beside the target, so the backup lands on
// the same filesystem and the restore is a rename rather than a copy.
func reserveBackup(target string) (string, error) {
	file, err := os.CreateTemp(filepath.Dir(target), "."+filepath.Base(target)+".backup-")
	if err != nil {
		return "", err
	}
	name := file.Name()
	if err := file.Close(); err != nil {
		os.Remove(name)
		return "", err
	}
	return name, os.Remove(name)
}

// Digest is the content hash CTR-006 records per asset, spelled the way the
// corpus manifest spells its own so the two can be compared by eye.
func Digest(content []byte) string {
	sum := sha256.Sum256(content)
	return "sha256:" + hex.EncodeToString(sum[:])
}

// ConfigDigest identifies the resolved configuration a generation ran under.
// It hashes the canonical JSON of the settings rather than the source file, so
// two spellings of the same configuration produce the same digest and a
// reformatted file does not look like a change.
func ConfigDigest(settings any) (string, error) {
	encoded, err := json.Marshal(settings)
	if err != nil {
		return "", fmt.Errorf("configuration could not be digested: %w", err)
	}
	return Digest(encoded), nil
}

package catalog

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
)

type isoSnapshot struct {
	ObservedAt string      `json:"observed_at"`
	Authority  string      `json:"authority"`
	Entities   []isoEntity `json:"entities"`
}

type isoEntity struct {
	Alpha2 string `json:"alpha2"`
	Alpha3 string `json:"alpha3"`
	Name   string `json:"name"`
}

type isoReconciliation struct {
	Authority      string                 `json:"authority"`
	ObservedAt     string                 `json:"observed_at"`
	SnapshotPath   string                 `json:"snapshot_path"`
	SnapshotSHA256 string                 `json:"snapshot_sha256"`
	Review         string                 `json:"review"`
	Rows           []isoReconciliationRow `json:"rows"`
}

type isoReconciliationRow struct {
	Alpha2      string `json:"alpha2"`
	Alpha3      string `json:"alpha3"`
	Name        string `json:"name"`
	Disposition string `json:"disposition"`
}

type receiptFile struct {
	Receipts []Receipt `json:"receipts"`
}

var (
	alpha2Pattern = regexp.MustCompile(`^[A-Z]{2}$`)
	alpha3Pattern = regexp.MustCompile(`^[A-Z]{3}$`)
	digestPattern = regexp.MustCompile(`^[0-9a-f]{64}$`)
)

func decodeStrict(path string, dst any) error {
	f, err := os.Open(path)
	if err != nil {
		return diagnostic(path, "-", "-", "input-readable", "%v", err)
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return diagnostic(path, "-", "-", "strict-json", "%v", err)
	}
	var trailing any
	if err := dec.Decode(&trailing); err != io.EOF {
		if err == nil {
			return diagnostic(path, "-", "-", "strict-json", "trailing JSON value")
		}
		return diagnostic(path, "-", "-", "strict-json", "trailing data: %v", err)
	}
	return nil
}

func loadISO(path string) ([]isoEntity, error) {
	var snapshot isoSnapshot
	if err := decodeStrict(path, &snapshot); err != nil {
		return nil, err
	}
	if snapshot.ObservedAt == "" || snapshot.Authority == "" {
		return nil, diagnostic(path, "-", "metadata", "iso-authority", "observation date and authority are required")
	}
	if len(snapshot.Entities) != 249 {
		return nil, diagnostic(path, "-", "entities", "iso-count-249", "got %d entities", len(snapshot.Entities))
	}
	seen2, seen3 := map[string]bool{}, map[string]bool{}
	for i, entity := range snapshot.Entities {
		if !alpha2Pattern.MatchString(entity.Alpha2) {
			return nil, diagnostic(path, entity.Alpha2, "alpha2", "iso-code", "invalid uppercase alpha-2 at index %d", i)
		}
		if !alpha3Pattern.MatchString(entity.Alpha3) {
			return nil, diagnostic(path, entity.Alpha2, "alpha3", "iso-code", "invalid uppercase alpha-3")
		}
		if entity.Name == "" {
			return nil, diagnostic(path, entity.Alpha2, "name", "iso-name", "empty display name")
		}
		if seen2[entity.Alpha2] || seen3[entity.Alpha3] {
			return nil, diagnostic(path, entity.Alpha2, "code", "iso-unique", "duplicate alpha code")
		}
		if i > 0 && snapshot.Entities[i-1].Alpha2 >= entity.Alpha2 {
			return nil, diagnostic(path, entity.Alpha2, "entities", "iso-sorted", "records are not strictly alpha-2 sorted")
		}
		seen2[entity.Alpha2], seen3[entity.Alpha3] = true, true
	}
	return snapshot.Entities, nil
}

func loadISOReconciliation(dataRoot, path string, entities []isoEntity) error {
	var attestation isoReconciliation
	if err := decodeStrict(path, &attestation); err != nil {
		return err
	}
	if attestation.Authority == "" || attestation.ObservedAt == "" || attestation.Review == "" ||
		filepath.ToSlash(attestation.SnapshotPath) != isoPath || !digestPattern.MatchString(attestation.SnapshotSHA256) {
		return diagnostic(path, "-", "metadata", "iso-reconciliation-bound", "complete authority, review, snapshot path and digest are required")
	}
	snapshotBytes, err := os.ReadFile(filepath.Join(dataRoot, filepath.FromSlash(attestation.SnapshotPath)))
	if err != nil {
		return diagnostic(path, "-", "snapshot_path", "iso-reconciliation-bound", "%v", err)
	}
	sum := sha256.Sum256(snapshotBytes)
	if hex.EncodeToString(sum[:]) != attestation.SnapshotSHA256 {
		return diagnostic(path, "-", "snapshot_sha256", "iso-reconciliation-bound", "snapshot digest mismatch")
	}
	if len(attestation.Rows) != len(entities) || len(attestation.Rows) != 249 {
		return diagnostic(path, "-", "rows", "iso-reconciliation-249", "got %d rows", len(attestation.Rows))
	}
	seen := map[string]bool{}
	for i, row := range attestation.Rows {
		if seen[row.Alpha2] {
			return diagnostic(path, row.Alpha2, "rows", "iso-reconciliation-unique", "duplicate row")
		}
		seen[row.Alpha2] = true
		entity := entities[i]
		if row.Alpha2 != entity.Alpha2 || row.Alpha3 != entity.Alpha3 || row.Name != entity.Name || row.Disposition != "matched" {
			return diagnostic(path, row.Alpha2, "rows", "iso-reconciliation-match", "row %d does not attest the bound snapshot", i)
		}
	}
	return nil
}

func loadAndVerifyReceipts(dataRoot string) ([]Receipt, error) {
	path := filepath.Join(dataRoot, "sources", "receipts.json")
	var file receiptFile
	if err := decodeStrict(path, &file); err != nil {
		return nil, err
	}
	if len(file.Receipts) == 0 {
		return nil, diagnostic(path, "-", "receipts", "source-receipts", "no receipts")
	}
	seen := map[string]bool{}
	for _, receipt := range file.Receipts {
		if receipt.ID == "" || seen[receipt.ID] {
			return nil, diagnostic(path, "-", "id", "receipt-unique", "empty or duplicate receipt id %q", receipt.ID)
		}
		seen[receipt.ID] = true
		if receipt.Path == "" || receipt.Origin == "" || receipt.UpstreamVersion == "" || receipt.License == "" || receipt.Transformation == "" || !digestPattern.MatchString(receipt.SHA256) {
			return nil, diagnostic(path, receipt.ID, "receipt", "receipt-complete", "origin, version, license, digest, path and transformation are required")
		}
		full := filepath.Join(dataRoot, filepath.FromSlash(receipt.Path))
		rel, err := filepath.Rel(dataRoot, full)
		if err != nil || rel == ".." || filepath.IsAbs(rel) || (len(rel) > 3 && rel[:3] == ".."+string(filepath.Separator)) {
			return nil, diagnostic(path, receipt.ID, "path", "source-contained", "path escapes data root")
		}
		b, err := os.ReadFile(full)
		if err != nil {
			return nil, diagnostic(full, receipt.ID, "path", "source-readable", "%v", err)
		}
		sum := sha256.Sum256(b)
		if hex.EncodeToString(sum[:]) != receipt.SHA256 {
			return nil, diagnostic(full, receipt.ID, "sha256", "source-digest", "acquisition digest mismatch")
		}
	}
	sort.Slice(file.Receipts, func(i, j int) bool { return file.Receipts[i].ID < file.Receipts[j].ID })
	return file.Receipts, nil
}

func requireReceipt(receipts []Receipt, path string) error {
	for _, receipt := range receipts {
		if filepath.ToSlash(receipt.Path) == filepath.ToSlash(path) {
			return nil
		}
	}
	return fmt.Errorf("source=%s entity=- field=receipt invariant=source-documented: missing receipt", path)
}

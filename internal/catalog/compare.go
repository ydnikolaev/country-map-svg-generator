package catalog

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

type Comparison struct {
	PreviousIdentity string   `json:"previous_identity,omitempty"`
	CurrentIdentity  string   `json:"current_identity"`
	Added            []string `json:"added"`
	Removed          []string `json:"removed"`
	Changed          []string `json:"changed"`
}

func Compare(previous, current *Corpus) Comparison {
	result := Comparison{CurrentIdentity: current.Manifest.Identity, Added: []string{}, Removed: []string{}, Changed: []string{}}
	if previous == nil {
		for _, entity := range current.Manifest.Entities {
			result.Added = append(result.Added, entity.Alpha2)
		}
		return result
	}
	result.PreviousIdentity = previous.Manifest.Identity
	old := map[string]Entity{}
	for _, entity := range previous.Manifest.Entities {
		old[entity.Alpha2] = entity
	}
	now := map[string]Entity{}
	for _, entity := range current.Manifest.Entities {
		now[entity.Alpha2] = entity
		if prior, ok := old[entity.Alpha2]; !ok {
			result.Added = append(result.Added, entity.Alpha2)
		} else {
			a, _ := CanonicalJSON(prior)
			b, _ := CanonicalJSON(entity)
			if !bytes.Equal(a, b) {
				result.Changed = append(result.Changed, entity.Alpha2)
			}
		}
	}
	for code := range old {
		if _, ok := now[code]; !ok {
			result.Removed = append(result.Removed, code)
		}
	}
	sort.Strings(result.Removed)
	return result
}

func Publish(directory string, corpus *Corpus) error {
	files, err := MarshalCorpus(corpus)
	if err != nil {
		return err
	}
	if existing, err := DecodeDirectory(directory); err == nil {
		existingFiles, _ := MarshalCorpus(existing)
		if equalFiles(existingFiles, files) {
			return nil
		}
		return fmt.Errorf("immutable corpus %s already exists with identity %s", directory, existing.Manifest.Identity)
	} else if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("existing corpus is unreadable: %w", err)
	}
	parent := filepath.Dir(directory)
	if err := os.MkdirAll(parent, 0o755); err != nil {
		return err
	}
	temp, err := os.MkdirTemp(parent, ".corpus-stage-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(temp)
	for name, data := range files {
		if err := os.WriteFile(filepath.Join(temp, name), data, 0o644); err != nil {
			return err
		}
	}
	decoded, err := DecodeDirectory(temp)
	if err != nil {
		return err
	}
	if decoded.Manifest.Identity != corpus.Manifest.Identity {
		return fmt.Errorf("staged corpus identity changed")
	}
	return os.Rename(temp, directory)
}

func DecodeDirectory(directory string) (*Corpus, error) {
	files := map[string][]byte{}
	expected := []string{"coverage.json", "geometries.json", "manifest.json", "receipts.json"}
	entries, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}
	if len(entries) != len(expected) {
		return nil, fmt.Errorf("corpus directory must contain exactly %d authoritative files", len(expected))
	}
	for i, entry := range entries {
		if entry.IsDir() || entry.Name() != expected[i] {
			return nil, fmt.Errorf("unexpected corpus directory entry %q", entry.Name())
		}
	}
	for _, name := range expected {
		b, err := os.ReadFile(filepath.Join(directory, name))
		if err != nil {
			return nil, err
		}
		files[name] = b
	}
	return DecodeCorpus(files)
}

func WriteComparison(path string, comparison Comparison) error {
	b, err := CanonicalJSON(comparison)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	if existing, err := os.ReadFile(path); err == nil {
		if bytes.Equal(existing, b) {
			return nil
		}
		return fmt.Errorf("immutable comparison report %s already exists with different bytes", path)
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return os.WriteFile(path, b, 0o644)
}

type generatedStages struct {
	Corpus   string
	Report   string
	Embedded string
}

// PublishGeneratedSet stages and validates the complete generated output set
// before replacing any destination. If a replacement fails, prior destinations
// are restored before the error is returned.
func PublishGeneratedSet(directory, reportPath, embeddedPath string, corpus *Corpus) error {
	return publishGeneratedSet(directory, reportPath, embeddedPath, corpus, nil)
}

func publishGeneratedSet(directory, reportPath, embeddedPath string, corpus *Corpus, mutateStages func(generatedStages) error) error {
	files, err := MarshalCorpus(corpus)
	if err != nil {
		return err
	}
	report, err := CanonicalJSON(Compare(nil, corpus))
	if err != nil {
		return err
	}
	embedded, err := EmbeddedSource(corpus)
	if err != nil {
		return err
	}
	corpusStage, err := stageCorpus(directory, files)
	if err != nil {
		return err
	}
	defer os.RemoveAll(corpusStage)
	reportStage, err := stageFile(reportPath, report)
	if err != nil {
		return err
	}
	defer os.Remove(reportStage)
	embeddedStage, err := stageFile(embeddedPath, embedded)
	if err != nil {
		return err
	}
	defer os.Remove(embeddedStage)
	stages := generatedStages{Corpus: corpusStage, Report: reportStage, Embedded: embeddedStage}
	if mutateStages != nil {
		if err := mutateStages(stages); err != nil {
			return err
		}
	}
	if err := validateGeneratedStages(stages, corpus, report, embedded); err != nil {
		return err
	}
	return replaceGeneratedSet([]replacement{
		{target: directory, staged: corpusStage},
		{target: reportPath, staged: reportStage},
		{target: embeddedPath, staged: embeddedStage},
	})
}

func stageCorpus(target string, files map[string][]byte) (string, error) {
	parent := filepath.Dir(target)
	if err := os.MkdirAll(parent, 0o755); err != nil {
		return "", err
	}
	stage, err := os.MkdirTemp(parent, ".corpus-set-stage-")
	if err != nil {
		return "", err
	}
	for name, data := range files {
		if err := os.WriteFile(filepath.Join(stage, name), data, 0o644); err != nil {
			os.RemoveAll(stage)
			return "", err
		}
	}
	return stage, nil
}

func stageFile(target string, data []byte) (string, error) {
	parent := filepath.Dir(target)
	if err := os.MkdirAll(parent, 0o755); err != nil {
		return "", err
	}
	file, err := os.CreateTemp(parent, "."+filepath.Base(target)+".set-stage-")
	if err != nil {
		return "", err
	}
	name := file.Name()
	if _, err := file.Write(data); err != nil {
		file.Close()
		os.Remove(name)
		return "", err
	}
	if err := file.Chmod(0o644); err != nil {
		file.Close()
		os.Remove(name)
		return "", err
	}
	if err := file.Close(); err != nil {
		os.Remove(name)
		return "", err
	}
	return name, nil
}

func validateGeneratedStages(stages generatedStages, corpus *Corpus, report, embedded []byte) error {
	decoded, err := DecodeDirectory(stages.Corpus)
	if err != nil {
		return fmt.Errorf("staged corpus: %w", err)
	}
	if decoded.Manifest.Identity != corpus.Manifest.Identity {
		return fmt.Errorf("staged corpus identity changed")
	}
	for path, want := range map[string][]byte{stages.Report: report, stages.Embedded: embedded} {
		got, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if !bytes.Equal(got, want) {
			return fmt.Errorf("staged generated artifact %s changed", path)
		}
	}
	return nil
}

type replacement struct {
	target string
	staged string
	backup string
}

func replaceGeneratedSet(items []replacement) error {
	for i := range items {
		if _, err := os.Lstat(items[i].target); err == nil {
			backup, err := reserveBackup(items[i].target)
			if err != nil {
				rollbackGeneratedSet(items, i-1)
				return err
			}
			if err := os.Rename(items[i].target, backup); err != nil {
				rollbackGeneratedSet(items, i-1)
				return err
			}
			items[i].backup = backup
		} else if !errors.Is(err, os.ErrNotExist) {
			rollbackGeneratedSet(items, i-1)
			return err
		}
		if err := os.Rename(items[i].staged, items[i].target); err != nil {
			if items[i].backup != "" {
				_ = os.Rename(items[i].backup, items[i].target)
				items[i].backup = ""
			}
			rollbackGeneratedSet(items, i-1)
			return err
		}
	}
	for _, item := range items {
		if item.backup != "" {
			_ = os.RemoveAll(item.backup)
		}
	}
	return nil
}

func reserveBackup(target string) (string, error) {
	file, err := os.CreateTemp(filepath.Dir(target), "."+filepath.Base(target)+".set-backup-")
	if err != nil {
		return "", err
	}
	name := file.Name()
	if err := file.Close(); err != nil {
		os.Remove(name)
		return "", err
	}
	if err := os.Remove(name); err != nil {
		return "", err
	}
	return name, nil
}

func rollbackGeneratedSet(items []replacement, last int) {
	for i := last; i >= 0; i-- {
		_ = os.RemoveAll(items[i].target)
		if items[i].backup != "" {
			_ = os.Rename(items[i].backup, items[i].target)
		}
	}
}

func equalFiles(a, b map[string][]byte) bool {
	if len(a) != len(b) {
		return false
	}
	for name, data := range a {
		if !bytes.Equal(data, b[name]) {
			return false
		}
	}
	return true
}

func CloneCorpus(corpus *Corpus) *Corpus {
	b, _ := json.Marshal(corpus)
	var clone Corpus
	_ = json.Unmarshal(b, &clone)
	return &clone
}

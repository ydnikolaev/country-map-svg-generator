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
	for _, name := range []string{"manifest.json", "geometries.json", "receipts.json", "coverage.json"} {
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

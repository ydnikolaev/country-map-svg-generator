package catalog

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
)

func CanonicalJSON(value any) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(value); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func geometryIdentity(coordinates MultiPolygon) (string, error) {
	b, err := CanonicalJSON(coordinates)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(b)
	return "geo-" + hex.EncodeToString(sum[:]), nil
}

func CorpusIdentity(manifest Manifest, geometries []Geometry, receipts []Receipt, coverage Coverage) (string, error) {
	manifest.Identity = ""
	envelope := struct {
		Manifest   Manifest   `json:"manifest"`
		Geometries []Geometry `json:"geometries"`
		Receipts   []Receipt  `json:"receipts"`
		Coverage   Coverage   `json:"coverage"`
	}{manifest, geometries, receipts, coverage}
	m, err := CanonicalJSON(envelope)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(m)
	return "sha256:" + hex.EncodeToString(sum[:]), nil
}

func MarshalCorpus(corpus *Corpus) (map[string][]byte, error) {
	files := map[string][]byte{}
	values := map[string]any{
		"manifest.json":   corpus.Manifest,
		"geometries.json": corpus.Geometries,
		"receipts.json":   corpus.Receipts,
		"coverage.json":   corpus.Coverage,
	}
	for name, value := range values {
		b, err := CanonicalJSON(value)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		files[name] = b
	}
	return files, nil
}

func DecodeCorpus(files map[string][]byte) (*Corpus, error) {
	var corpus Corpus
	targets := map[string]any{
		"manifest.json":   &corpus.Manifest,
		"geometries.json": &corpus.Geometries,
		"receipts.json":   &corpus.Receipts,
		"coverage.json":   &corpus.Coverage,
	}
	if len(files) != len(targets) {
		return nil, fmt.Errorf("corpus file set must contain exactly %d authoritative files", len(targets))
	}
	for name, target := range targets {
		b, ok := files[name]
		if !ok {
			return nil, fmt.Errorf("corpus file %s is missing", name)
		}
		dec := json.NewDecoder(bytes.NewReader(b))
		dec.DisallowUnknownFields()
		if err := dec.Decode(target); err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		var trailing any
		if err := dec.Decode(&trailing); err != io.EOF {
			return nil, fmt.Errorf("%s: non-canonical trailing data", name)
		}
		canonical, err := CanonicalJSON(target)
		if err != nil {
			return nil, err
		}
		if !bytes.Equal(canonical, b) {
			return nil, fmt.Errorf("%s: published bytes are not canonical", name)
		}
	}
	if err := Validate(&corpus); err != nil {
		return nil, err
	}
	return &corpus, nil
}

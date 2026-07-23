package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/yuranikolaev/country-map-svg-generator/internal/catalog"
)

func main() {
	data := flag.String("data", "data", "root containing sources and policy")
	output := flag.String("out", "", "immutable corpus output directory")
	report := flag.String("report", "", "comparison report path")
	embedded := flag.String("embedded", "", "generated Go embedded bundle path")
	flag.Parse()
	if *output == "" || *report == "" || *embedded == "" {
		fmt.Fprintln(os.Stderr, "-out, -report and -embedded are required")
		os.Exit(2)
	}
	corpus, err := catalog.Compile(*data)
	if err != nil {
		fail(err)
	}
	if err := catalog.Publish(*output, corpus); err != nil {
		fail(err)
	}
	if err := catalog.WriteComparison(*report, catalog.Compare(nil, corpus)); err != nil {
		fail(err)
	}
	source, err := catalog.EmbeddedSource(corpus)
	if err != nil {
		fail(err)
	}
	if err := writeCreateOrEqual(*embedded, source); err != nil {
		fail(err)
	}
	fmt.Printf("corpus=%s entities=%d geometries=%d capitals=%d protected=%d\n",
		corpus.Manifest.Identity, corpus.Manifest.EntityCount, len(corpus.Geometries),
		corpus.Coverage.CapitalRecords, len(corpus.Coverage.Protected))
}

func writeCreateOrEqual(path string, data []byte) error {
	if existing, err := os.ReadFile(path); err == nil {
		if string(existing) == string(data) {
			return nil
		}
		if filepath.Base(path) != "embedded_gen.go" || filepath.Dir(path) != filepath.FromSlash("internal/catalog") {
			return fmt.Errorf("refusing to overwrite different generated file %s", path)
		}
	} else if !os.IsNotExist(err) {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	temp, err := os.CreateTemp(filepath.Dir(path), ".embedded-stage-")
	if err != nil {
		return err
	}
	name := temp.Name()
	defer os.Remove(name)
	if _, err := temp.Write(data); err != nil {
		temp.Close()
		return err
	}
	if err := temp.Chmod(0o644); err != nil {
		temp.Close()
		return err
	}
	if err := temp.Close(); err != nil {
		return err
	}
	return os.Rename(name, path)
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

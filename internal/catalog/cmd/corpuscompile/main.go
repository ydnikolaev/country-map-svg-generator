package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ydnikolaev/country-map-svg-generator/internal/catalog"
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
	if err := catalog.PublishGeneratedSet(*output, *report, *embedded, corpus); err != nil {
		fail(err)
	}
	fmt.Printf("corpus=%s entities=%d geometries=%d capitals=%d protected=%d\n",
		corpus.Manifest.Identity, corpus.Manifest.EntityCount, len(corpus.Geometries),
		corpus.Coverage.CapitalRecords, len(corpus.Coverage.Protected))
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

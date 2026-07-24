// svgproof renders the shipped Generate() path over the whole catalog and
// writes real SVG documents plus a contact sheet. It is the empirical answer to
// "the tests pass" — a maintainer runs it, opens the files, and sees whether the
// product actually draws countries.
//
// It is deliberately not part of `make check`: writing ~2000 files is not a
// gate. The gate is TestShippedCatalogServesTheLadder, which asserts the same
// invariants without writing anything.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/yuranikolaev/country-map-svg-generator/internal/catalog"
	geometry "github.com/yuranikolaev/country-map-svg-generator/internal/geometry"
)

type proofRow struct {
	Alpha2, Profile, Band, Preset string
	Outcome                       string
	SelectedTier                  string
	Selection                     string
	PathBytes                     int
	CommittedBytes                int
	BandPathCap                   int
	Parts                         int
	ViewBoxW, ViewBoxH            float64
	Note                          string
}

func main() {
	out := flag.String("out", "", "directory to write SVG files and the contact sheet into")
	only := flag.String("only", "", "comma-separated alpha2 filter; empty means the whole catalog")
	flag.Parse()

	if err := run(*out, *only); err != nil {
		fmt.Fprintln(os.Stderr, "svgproof:", err)
		os.Exit(1)
	}
}

func run(out, only string) error {
	c, err := catalog.Embedded()
	if err != nil {
		return err
	}
	oracle, err := geometry.EmbeddedSilhouetteOracle()
	if err != nil {
		return err
	}
	ladder, err := geometry.EmbeddedLadderTable()
	if err != nil {
		return err
	}
	filter := map[string]bool{}
	for _, code := range strings.Split(only, ",") {
		if code = strings.TrimSpace(code); code != "" {
			filter[code] = true
		}
	}
	if out != "" {
		if err := os.MkdirAll(out, 0o755); err != nil {
			return err
		}
	}

	bands := map[string]geometry.SilhouetteBand{}
	for _, band := range oracle.Bands {
		bands[band.ID] = band
	}
	presets := map[string]string{"compact": "card", "standard": "hero"}

	rows := make([]proofRow, 0, len(ladder.Rows))
	for _, committed := range ladder.Rows {
		if len(filter) > 0 && !filter[committed.Alpha2] {
			continue
		}
		band := bands[committed.Band]
		preset := presets[committed.Band]
		row := proofRow{
			Alpha2: committed.Alpha2, Profile: committed.Profile, Band: committed.Band, Preset: preset,
			Selection: committed.Selection, CommittedBytes: committed.PathBytes, BandPathCap: band.PathCap,
		}
		in, err := geometry.InputFromCatalog(c, committed.Alpha2, committed.Profile, preset)
		if err != nil {
			return fmt.Errorf("%s/%s: %w", committed.Alpha2, committed.Profile, err)
		}
		result, genErr := geometry.Generate(in)
		switch {
		case geometry.IsNoArtifact(genErr):
			row.Outcome = "no_artifact"
			row.Note = committed.NoArtifactReason
		case genErr != nil:
			row.Outcome = "error"
			row.Note = genErr.Error()
		default:
			row.Outcome = "rendered"
			row.SelectedTier = result.LOD.SelectedTier
			row.PathBytes = len(result.Path)
			row.Parts = strings.Count(result.Path, "M")
			row.ViewBoxW, row.ViewBoxH = result.ViewBox.Width(), result.ViewBox.Height()
			if out != "" {
				name := fmt.Sprintf("%s.%s.%s.svg", committed.Alpha2, committed.Profile, committed.Band)
				if err := os.WriteFile(filepath.Join(out, name), []byte(svgDocument(result)), 0o644); err != nil {
					return err
				}
			}
		}
		rows = append(rows, row)
	}

	sort.Slice(rows, func(i, j int) bool {
		if rows[i].Alpha2 != rows[j].Alpha2 {
			return rows[i].Alpha2 < rows[j].Alpha2
		}
		if rows[i].Profile != rows[j].Profile {
			return rows[i].Profile < rows[j].Profile
		}
		return rows[i].Band < rows[j].Band
	})

	rendered, noArtifact, errored := 0, 0, 0
	sourceFallback, overCap, byteMatch := 0, 0, 0
	for _, row := range rows {
		switch row.Outcome {
		case "rendered":
			rendered++
			if row.SelectedTier == "source" {
				sourceFallback++
				fmt.Printf("[source-fallback] %s/%s/%s bytes=%d cap=%d\n", row.Alpha2, row.Profile, row.Band, row.PathBytes, row.BandPathCap)
			}
			if row.PathBytes > row.BandPathCap {
				overCap++
				fmt.Printf("[over-band-cap] %s/%s/%s bytes=%d cap=%d committed=%d tier=%s\n",
					row.Alpha2, row.Profile, row.Band, row.PathBytes, row.BandPathCap, row.CommittedBytes, row.SelectedTier)
			}
			if row.PathBytes == row.CommittedBytes {
				byteMatch++
			}
		case "no_artifact":
			noArtifact++
		default:
			errored++
			fmt.Printf("[error] %s/%s/%s %s\n", row.Alpha2, row.Profile, row.Band, row.Note)
		}
	}
	fmt.Printf("\nrows=%d rendered=%d no_artifact=%d error=%d source_fallback=%d over_band_cap=%d byte_identical_to_committed=%d\n",
		len(rows), rendered, noArtifact, errored, sourceFallback, overCap, byteMatch)

	if out != "" {
		raw, err := json.MarshalIndent(rows, "", " ")
		if err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(out, "contact-sheet.json"), append(raw, '\n'), 0o644); err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(out, "index.html"), []byte(contactSheet(rows)), 0o644); err != nil {
			return err
		}
		fmt.Printf("wrote %s\n", filepath.Join(out, "index.html"))
	}
	return nil
}

// svgDocument is a minimal, presentation-free wrapper (INV-1): geometry only,
// no fill, stroke or colour. Styling is P3's concern; this exists so a human
// can open the file.
func svgDocument(r geometry.Result) string {
	var b strings.Builder
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" viewBox="%g %g %g %g" role="img" aria-label="%s %s">`,
		r.ViewBox.MinX, r.ViewBox.MinY, r.ViewBox.Width(), r.ViewBox.Height(), r.Entity, r.Profile)
	fmt.Fprintf(&b, `<path d="%s"/>`, r.Path)
	b.WriteString(`</svg>`)
	return b.String()
}

func contactSheet(rows []proofRow) string {
	var b strings.Builder
	b.WriteString("<!doctype html><meta charset=\"utf-8\"><title>country-map-svg-generator proof</title>")
	b.WriteString("<style>body{font:13px system-ui;margin:24px;background:#111;color:#eee}" +
		"h1{font-size:16px}.grid{display:grid;grid-template-columns:repeat(auto-fill,minmax(150px,1fr));gap:14px}" +
		".c{background:#1b1b1b;border-radius:8px;padding:8px;text-align:center}" +
		".c svg{width:100%;height:110px;fill:#e8e8e8}.m{color:#888;font-size:11px;margin-top:4px}" +
		".bad{outline:2px solid #c33}.na{color:#c96}</style>")
	fmt.Fprintf(&b, "<h1>%d cells — shipped Generate() output</h1><div class=\"grid\">", len(rows))
	for _, row := range rows {
		cls := "c"
		if row.SelectedTier == "source" || row.PathBytes > row.BandPathCap {
			cls += " bad"
		}
		fmt.Fprintf(&b, "<div class=\"%s\">", cls)
		switch row.Outcome {
		case "rendered":
			fmt.Fprintf(&b, "<object type=\"image/svg+xml\" data=\"%s.%s.%s.svg\" style=\"width:100%%;height:110px\"></object>",
				row.Alpha2, row.Profile, row.Band)
			fmt.Fprintf(&b, "<div><b>%s</b> %s/%s</div><div class=\"m\">rung %s · %d B / %d · tier %s</div>",
				row.Alpha2, row.Profile, row.Band, row.Selection, row.PathBytes, row.BandPathCap, row.SelectedTier)
		case "no_artifact":
			fmt.Fprintf(&b, "<div style=\"height:110px;display:flex;align-items:center;justify-content:center\" class=\"na\">no artifact</div>")
			fmt.Fprintf(&b, "<div><b>%s</b> %s/%s</div><div class=\"m\">%s</div>", row.Alpha2, row.Profile, row.Band, row.Note)
		default:
			fmt.Fprintf(&b, "<div style=\"height:110px;color:#c33\">error</div><div><b>%s</b> %s/%s</div><div class=\"m\">%s</div>",
				row.Alpha2, row.Profile, row.Band, row.Note)
		}
		b.WriteString("</div>")
	}
	b.WriteString("</div>")
	return b.String()
}

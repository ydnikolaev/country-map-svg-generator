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

	// svg is the inlined document, kept out of the JSON record so the machine
	// evidence stays readable while the human page stays self-contained.
	svg string `json:"-"`
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
			row.svg = svgDocument(result)
			if out != "" {
				name := fmt.Sprintf("%s.%s.%s.svg", committed.Alpha2, committed.Profile, committed.Band)
				if err := os.WriteFile(filepath.Join(out, name), []byte(row.svg), 0o644); err != nil {
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

// contactSheet writes one self-contained page: every silhouette is inlined, so
// the file can be opened or sent on its own with no directory around it. Filters
// are plain CSS radio state rather than script, so it also survives being
// rendered somewhere that will not run JavaScript.
func contactSheet(rows []proofRow) string {
	rendered, absent, over, fallback := 0, 0, 0, 0
	budget := map[string][]float64{}
	for _, row := range rows {
		switch row.Outcome {
		case "rendered":
			rendered++
			if row.PathBytes > row.BandPathCap {
				over++
			}
			if row.SelectedTier == "source" {
				fallback++
			}
			if row.BandPathCap > 0 {
				budget[row.Band] = append(budget[row.Band], float64(row.PathBytes)/float64(row.BandPathCap))
			}
		case "no_artifact":
			absent++
		}
	}

	var b strings.Builder
	b.WriteString(`<!doctype html><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1">`)
	b.WriteString(`<title>Country silhouettes — full catalog</title><style>
:root{--bg:#0e0f11;--card:#191b1f;--ink:#e9eaec;--dim:#8b9099;--line:#2a2d33;--warn:#d8a657;--bad:#e06c75}
*{box-sizing:border-box}
body{font:14px/1.5 ui-sans-serif,system-ui,sans-serif;margin:0;padding:28px;background:var(--bg);color:var(--ink)}
h1{font-size:19px;margin:0 0 4px}
.sub{color:var(--dim);font-size:13px;margin-bottom:20px}
.stats{display:flex;flex-wrap:wrap;gap:10px;margin-bottom:22px}
.stat{background:var(--card);border:1px solid var(--line);border-radius:10px;padding:10px 14px;min-width:118px}
.stat b{display:block;font-size:20px;font-variant-numeric:tabular-nums}
.stat span{color:var(--dim);font-size:12px}
.filters{display:flex;gap:8px;margin-bottom:18px;flex-wrap:wrap}
.filters label{background:var(--card);border:1px solid var(--line);border-radius:999px;padding:6px 14px;cursor:pointer;font-size:13px;color:var(--dim)}
.filters input{position:absolute;opacity:0;pointer-events:none}
.filters input:checked+label{color:var(--ink);border-color:#4a5568;background:#22252b}
.grid{display:grid;grid-template-columns:repeat(auto-fill,minmax(158px,1fr));gap:12px}
.c{background:var(--card);border:1px solid var(--line);border-radius:10px;padding:10px;text-align:center;overflow:hidden}
.c svg{width:100%;height:112px;display:block;fill:var(--ink)}
.c .t{font-weight:600;margin-top:8px;font-size:13px}
.m{color:var(--dim);font-size:11px;font-variant-numeric:tabular-nums;word-break:break-word}
.flag{outline:2px solid var(--bad)}
.na{height:112px;display:flex;align-items:center;justify-content:center;color:var(--warn);font-size:12px;
    border:1px dashed #3a3d44;border-radius:6px}
#fa:checked~.grid .band-standard,#fb:checked~.grid .band-compact,
#fc:checked~.grid .ok{display:none}
@media (prefers-color-scheme:light){
:root{--bg:#fafafa;--card:#fff;--ink:#16181d;--dim:#6b7280;--line:#e3e5e9}
.filters input:checked+label{background:#eef0f4;border-color:#c3c8d0}}
</style>`)
	fmt.Fprintf(&b, `<h1>Country silhouettes — full catalog</h1>
<div class="sub">Every cell is the shipped <code>Generate()</code> output, inlined. Geometry only — no fill, stroke or colour is baked in.</div>`)
	fmt.Fprintf(&b, `<div class="stats">
<div class="stat"><b>%d</b><span>rendered</span></div>
<div class="stat"><b>%d</b><span>no artifact</span></div>
<div class="stat"><b>%d</b><span>over band cap</span></div>
<div class="stat"><b>%d</b><span>source fallback</span></div>`, rendered, absent, over, fallback)
	for _, band := range []string{"compact", "standard"} {
		used := append([]float64(nil), budget[band]...)
		if len(used) == 0 {
			continue
		}
		sort.Float64s(used)
		fmt.Fprintf(&b, `<div class="stat"><b>%.0f%%</b><span>median budget, %s</span></div>`,
			100*used[len(used)/2], band)
	}
	b.WriteString(`</div>`)
	b.WriteString(`<input type="radio" name="f" id="f0" checked><label for="f0">all</label>
<input type="radio" name="f" id="fa"><label for="fa">cards only</label>
<input type="radio" name="f" id="fb"><label for="fb">heroes only</label>
<input type="radio" name="f" id="fc"><label for="fc">needs a look</label>`)
	b.WriteString(`<div class="filters"></div><div class="grid">`)
	for _, row := range rows {
		flagged := row.Outcome != "rendered" || row.SelectedTier == "source" || row.PathBytes > row.BandPathCap
		cls := "c band-" + row.Band
		if flagged {
			cls += " flag"
		} else {
			cls += " ok"
		}
		fmt.Fprintf(&b, `<div class="%s">`, cls)
		switch row.Outcome {
		case "rendered":
			b.WriteString(row.svg)
			fmt.Fprintf(&b, `<div class="t">%s <span class="m">%s</span></div><div class="m">%s · rung %s · %d/%d B</div>`,
				row.Alpha2, row.Profile, row.Band, row.Selection, row.PathBytes, row.BandPathCap)
		case "no_artifact":
			fmt.Fprintf(&b, `<div class="na">no artifact</div><div class="t">%s <span class="m">%s</span></div><div class="m">%s · %s</div>`,
				row.Alpha2, row.Profile, row.Band, row.Note)
		default:
			fmt.Fprintf(&b, `<div class="na" style="color:var(--bad)">error</div><div class="t">%s <span class="m">%s</span></div><div class="m">%s</div>`,
				row.Alpha2, row.Profile, row.Note)
		}
		b.WriteString(`</div>`)
	}
	b.WriteString(`</div>`)
	return b.String()
}

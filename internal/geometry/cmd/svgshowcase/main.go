// svgshowcase renders one page that answers "what can this thing actually do".
//
// Every silhouette on the page is the same shipped Generate() output the catalog
// sheet uses — the only thing that changes between treatments is CSS. That is
// the point rather than a shortcut: INV-1 keeps the geometry presentation-free,
// so fill, stroke, opacity, layering and capital markers are all decisions a
// consumer makes later, on one path.
package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/yuranikolaev/country-map-svg-generator/internal/catalog"
	geometry "github.com/yuranikolaev/country-map-svg-generator/internal/geometry"
)

type cell struct {
	Code, Profile, Band, Preset string
	Result                      geometry.Result
	Cropped                     bool
}

func main() {
	out := flag.String("out", "", "file to write the showcase page to")
	flag.Parse()
	if *out == "" {
		fmt.Fprintln(os.Stderr, "svgshowcase: -out is required")
		os.Exit(2)
	}
	if err := run(*out); err != nil {
		fmt.Fprintln(os.Stderr, "svgshowcase:", err)
		os.Exit(1)
	}
}

func run(out string) error {
	c, err := catalog.Embedded()
	if err != nil {
		return err
	}
	load := func(code, band string) (cell, error) {
		preset := "card"
		if band == "standard" {
			preset = "hero"
		}
		in, err := geometry.InputFromCatalog(c, code, "un", preset)
		if err != nil {
			return cell{}, err
		}
		result, err := geometry.Generate(in)
		if err != nil {
			return cell{}, fmt.Errorf("%s/%s: %w", code, band, err)
		}
		return cell{Code: code, Profile: "un", Band: band, Preset: preset, Result: result}, nil
	}

	hero, err := load("FR", "standard")
	if err != nil {
		return err
	}
	gallery := []cell{}
	for _, code := range []string{"IT", "JP", "GR", "NZ", "BR", "IN", "PT", "NO", "ID", "CL"} {
		cell, err := load(code, "standard")
		if err != nil {
			return err
		}
		gallery = append(gallery, cell)
	}
	usFull, err := load("US", "standard")
	if err != nil {
		return err
	}
	page := render(hero, gallery, usFull)
	if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(out, []byte(page), 0o644); err != nil {
		return err
	}
	fmt.Printf("wrote %s (%d cells)\n", out, len(gallery)+6)
	return nil
}

// croppedViewBox frames the largest subpath of an already-generated result.
// Nothing about the geometry changes — the path is emitted untouched and the
// viewBox does the work, so this is a consumer decision available today with no
// pipeline change at all. It is the honest version of the demo: the frame is
// presentation, and presentation is exactly what INV-1 hands over.
func croppedViewBox(r geometry.Result) (minX, minY, w, h float64, ok bool) {
	type box struct{ minX, minY, maxX, maxY float64 }
	var boxes []box
	var cur *box
	for _, command := range r.Commands {
		if command.Op == "M" || command.Op == "m" {
			boxes = append(boxes, box{math.Inf(1), math.Inf(1), math.Inf(-1), math.Inf(-1)})
			cur = &boxes[len(boxes)-1]
		}
		if cur == nil {
			continue
		}
		for i := 0; i+1 < len(command.Values); i += 2 {
			x, y := command.Values[i], command.Values[i+1]
			cur.minX, cur.maxX = math.Min(cur.minX, x), math.Max(cur.maxX, x)
			cur.minY, cur.maxY = math.Min(cur.minY, y), math.Max(cur.maxY, y)
		}
	}
	best, bestArea := -1, 0.0
	for i, b := range boxes {
		if b.maxX <= b.minX || b.maxY <= b.minY {
			continue
		}
		if a := (b.maxX - b.minX) * (b.maxY - b.minY); a > bestArea {
			best, bestArea = i, a
		}
	}
	if best < 0 {
		return 0, 0, 0, 0, false
	}
	b := boxes[best]
	pad := math.Max(b.maxX-b.minX, b.maxY-b.minY) * 0.04
	return b.minX - pad, b.minY - pad, (b.maxX - b.minX) + 2*pad, (b.maxY - b.minY) + 2*pad, true
}

// svgFramed emits the same untouched path inside an explicit frame.
func svgFramed(r geometry.Result, minX, minY, w, h float64) string {
	return fmt.Sprintf(`<svg viewBox="%g %g %g %g" role="img" aria-label="%s"><path class="land" d="%s"/></svg>`,
		minX, minY, w, h, r.Entity, r.Path)
}

// svg emits one silhouette. The path carries no presentation of its own; the
// classes are hooks the stylesheet decides on.
func svg(r geometry.Result, markers bool) string {
	var b strings.Builder
	fmt.Fprintf(&b, `<svg viewBox="%g %g %g %g" role="img" aria-label="%s">`,
		r.ViewBox.MinX, r.ViewBox.MinY, r.ViewBox.Width(), r.ViewBox.Height(), r.Entity)
	fmt.Fprintf(&b, `<path class="land" d="%s"/>`, r.Path)
	if markers {
		radius := r.ViewBox.Width() / 55
		for _, marker := range r.Markers {
			if marker.Anomaly != "" {
				continue
			}
			fmt.Fprintf(&b, `<circle class="capital" cx="%g" cy="%g" r="%g"/>`, marker.X, marker.Y, radius)
		}
	}
	b.WriteString(`</svg>`)
	return b.String()
}

func render(hero cell, gallery []cell, usFull cell) string {
	var b strings.Builder
	b.WriteString(`<!doctype html><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1">`)
	b.WriteString(`<title>One path, many treatments</title>`)
	b.WriteString(style)

	fmt.Fprintf(&b, `<header>
<p class="eyebrow">country-map-svg-generator</p>
<h1>One path,<br>many treatments</h1>
<p class="lede">Every silhouette below is the same shipped output — a single <code>&lt;path&gt;</code> with no fill, stroke or colour baked into it. Everything you see is a stylesheet decision made afterwards. That separation is the product: the generator commits to geometry and budgets, and hands presentation to whoever is drawing.</p>
</header>`)

	// The hero, one country through five treatments.
	b.WriteString(`<section><h2>The same France, five ways</h2>
<p class="note">Identical path data in all five. Only CSS differs.</p>
<div class="row five">`)
	for _, t := range []struct{ cls, name, note string }{
		{"t-solid", "Solid", "fill only"},
		{"t-outline", "Outline", "stroke, no fill"},
		{"t-ghost", "Ghost", "stroke + translucent fill"},
		{"t-capital", "Capitals", "markers from the corpus"},
		{"t-offset", "Offset", "one path, drawn twice"},
	} {
		fmt.Fprintf(&b, `<figure class="cell %s"><div class="art">%s</div><figcaption><b>%s</b><span>%s</span></figcaption></figure>`,
			t.cls, svg(hero.Result, t.cls == "t-capital"), t.name, t.note)
	}
	fmt.Fprintf(&b, `</div><p class="meta">FR · un · hero · %d bytes of a 7500 budget · viewBox %.0f×%.0f</p></section>`,
		len(hero.Result.Path), hero.Result.ViewBox.Width(), hero.Result.ViewBox.Height())

	// Breadth.
	b.WriteString(`<section><h2>Ten more, one treatment each</h2>
<p class="note">Cycling the same five treatments across different geographies — an archipelago, a long thin coast, an island chain, a subcontinent.</p>
<div class="row grid">`)
	treatments := []string{"t-solid", "t-ghost", "t-outline", "t-capital", "t-offset"}
	for i, cell := range gallery {
		t := treatments[i%len(treatments)]
		fmt.Fprintf(&b, `<figure class="cell %s"><div class="art">%s</div><figcaption><b>%s</b><span>%d B · %d parts</span></figcaption></figure>`,
			t, svg(cell.Result, t == "t-capital"), cell.Code, len(cell.Result.Path), partCount(cell.Result))
	}
	b.WriteString(`</div></section>`)

	// The crop capability.
	fmt.Fprintf(&b, `<section><h2>Framing follows the geometry</h2>
<p class="note">Both cells hold the byte-for-byte identical path. The right one only changes the <code>viewBox</code>, and the outlying components fall outside it — so a consumer can reframe to the mainland today, with no pipeline change and no extra bytes.</p>
<div class="row two">
<figure class="cell t-solid"><div class="art">%s</div><figcaption><b>US</b><span>everything · %d B</span></figcaption></figure>
<figure class="cell t-solid"><div class="art">%s</div><figcaption><b>US</b><span>framed to the mainland · same %d B</span></figcaption></figure>
</div></section>`,
		svg(usFull.Result, false), len(usFull.Result.Path),
		croppedSVG(usFull.Result), len(usFull.Result.Path))

	b.WriteString(`<footer><p>Geometry is presentation-free by design (INV-1). Colour, stroke, opacity, layering and markers are all applied here in CSS over unmodified output.</p></footer>`)
	return b.String()
}

func croppedSVG(r geometry.Result) string {
	minX, minY, w, h, ok := croppedViewBox(r)
	if !ok {
		return svg(r, false)
	}
	return svgFramed(r, minX, minY, w, h)
}

func partCount(r geometry.Result) int {
	n := 0
	for _, command := range r.Commands {
		if command.Op == "M" || command.Op == "m" {
			n++
		}
	}
	return n
}

const style = `<style>
/* An atlas plate rather than a dashboard: paper ground, indigo ink, and a
   brick accent reserved for the one thing that is not geography — the capitals. */
:root{
  --paper:#f2f1ec; --plate:#fbfaf7; --ink:#161a24; --ink-soft:#5c6474;
  --rule:#ddd9cf; --sea:#e6e9e4; --accent:#a8402c; --accent-soft:#d99a8a; --land:#243044;
}
@media (prefers-color-scheme:dark){
  :root{--paper:#0f1116; --plate:#171a21; --ink:#e8e9ec; --ink-soft:#949bab;
        --rule:#272b35; --sea:#1d2129; --accent:#e0715a; --accent-soft:#7a3a2e; --land:#c9d2e2;}
}
:root[data-theme="dark"]{--paper:#0f1116; --plate:#171a21; --ink:#e8e9ec; --ink-soft:#949bab;
  --rule:#272b35; --sea:#1d2129; --accent:#e0715a; --accent-soft:#7a3a2e; --land:#c9d2e2;}
:root[data-theme="light"]{--paper:#f2f1ec; --plate:#fbfaf7; --ink:#161a24; --ink-soft:#5c6474;
  --rule:#ddd9cf; --sea:#e6e9e4; --accent:#a8402c; --accent-soft:#d99a8a; --land:#243044;}
*{box-sizing:border-box}
body{margin:0;padding:clamp(24px,5vw,64px);background:var(--paper);color:var(--ink);
  font:16px/1.6 ui-sans-serif,system-ui,-apple-system,sans-serif;
  display:flex;flex-direction:column;gap:clamp(40px,7vw,84px)}
header{max-width:62ch;display:flex;flex-direction:column;gap:14px}
.eyebrow{margin:0;font:600 11px/1 ui-monospace,SFMono-Regular,Menlo,monospace;
  letter-spacing:.16em;text-transform:uppercase;color:var(--ink-soft)}
h1{margin:0;font:400 clamp(34px,6vw,60px)/1.04 ui-serif,Georgia,"Times New Roman",serif;
  letter-spacing:-.015em;text-wrap:balance}
.lede{margin:0;color:var(--ink-soft);font-size:17px;text-wrap:pretty}
code{font:14px ui-monospace,SFMono-Regular,Menlo,monospace;background:var(--sea);
  padding:1px 5px;border-radius:4px}
section{display:flex;flex-direction:column;gap:14px}
h2{margin:0;font:400 clamp(21px,2.6vw,27px)/1.2 ui-serif,Georgia,serif;letter-spacing:-.01em;
  padding-bottom:10px;border-bottom:1px solid var(--rule)}
.note{margin:0;color:var(--ink-soft);max-width:70ch;font-size:15px}
.meta{margin:0;color:var(--ink-soft);font:12px ui-monospace,SFMono-Regular,Menlo,monospace;
  font-variant-numeric:tabular-nums}
.row{display:grid;gap:14px}
.five{grid-template-columns:repeat(auto-fit,minmax(180px,1fr))}
.two{grid-template-columns:repeat(auto-fit,minmax(280px,1fr))}
.grid{grid-template-columns:repeat(auto-fill,minmax(158px,1fr))}
.cell{margin:0;background:var(--plate);border:1px solid var(--rule);border-radius:3px;
  padding:16px;display:flex;flex-direction:column;gap:12px}
.art{display:flex;align-items:center;justify-content:center;min-height:0}
.cell svg{width:100%;height:132px;overflow:visible}
figcaption{display:flex;flex-direction:column;gap:2px;border-top:1px solid var(--rule);padding-top:9px}
figcaption b{font-size:14px;font-weight:600}
figcaption span{color:var(--ink-soft);font:11px ui-monospace,SFMono-Regular,Menlo,monospace;
  font-variant-numeric:tabular-nums}
footer{border-top:1px solid var(--rule);padding-top:18px;color:var(--ink-soft);max-width:70ch;font-size:14px}
footer p{margin:0}

/* The treatments. Every rule below acts on the same untouched path. */
.land{fill:none;stroke:none;vector-effect:non-scaling-stroke}
.t-solid .land{fill:var(--land)}
.t-outline .land{stroke:var(--land);stroke-width:1.4;stroke-linejoin:round}
/* fill-opacity rather than color-mix: a translucent body has to survive
   everywhere the page might be opened, and a silently unsupported colour
   function would fail as a solid shape rather than as an error. */
.t-ghost .land{fill:var(--land);fill-opacity:.16;
  stroke:var(--land);stroke-width:1.2;stroke-linejoin:round}
.t-capital .land{fill:var(--land);fill-opacity:.22;stroke:var(--land);stroke-width:1}
.capital{fill:var(--accent)}
.t-offset .land{fill:var(--land);filter:drop-shadow(3px 3px 0 var(--accent-soft))}
@media (prefers-reduced-motion:reduce){*{animation:none!important;transition:none!important}}
</style>`

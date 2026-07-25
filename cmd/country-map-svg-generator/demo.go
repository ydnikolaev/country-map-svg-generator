package main

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ydnikolaev/country-map-svg-generator/internal/catalog"
	"github.com/ydnikolaev/country-map-svg-generator/internal/config"
	"github.com/ydnikolaev/country-map-svg-generator/internal/geometry"
	"github.com/ydnikolaev/country-map-svg-generator/internal/render"
)

//go:embed demo/page.html
var demoTemplate string

// demoEntities is an editorial selection, and the only one in this repository.
// It is chosen for SHAPE DIVERSITY rather than for size or importance: a demo
// that showed ten compact European countries would prove nothing about the
// pipeline. A ranking would also invite an argument the page does not need.
//
// The no-country-literal invariant is `internal/geometry`'s, where a per-entity
// branch would mean the pipeline treats one place specially. Naming entities in
// a demo page is the opposite: it changes no behaviour, and every one of them
// goes through the identical path.
var demoEntities = []struct{ ISO, Why string }{
	{"FR", "compact, with an offshore component"},
	{"BR", "large and compact, long coastline"},
	{"KZ", "elongated east to west"},
	{"IN", "compact with a peninsula"},
	{"AU", "island continent, isolated silhouette"},
	{"ZA", "encloses another state, so the fill rule shows"},
	{"IT", "distinctive outline plus islands"},
	{"JP", "archipelago, arc-shaped"},
	{"NO", "extreme coastline complexity"},
	{"CA", "vast, deeply indented, arctic islands"},
	{"DE", "compact, the densest silhouette in the set"},
	{"ES", "near-square landmass with island groups"},
	{"GB", "island, heavily indented"},
	{"US", "spans a continent plus two detached states"},
	{"CN", "large and compact, complex southern edge"},
	{"RU", "widest span in the corpus"},
	{"MX", "long curve, tapering"},
	{"AR", "north-south wedge"},
	{"EG", "near-rectangular, an unusual outline"},
	{"TR", "wide and low, two seas"},
	{"SE", "elongated north-south, Baltic islands"},
	{"PL", "compact, almost no coastline detail"},
	{"UA", "compact with a peninsula"},
	{"TH", "narrow tail, awkward to fit"},
	{"VN", "extreme S-curve"},
	{"ID", "archipelago spanning the widest arc"},
	{"NG", "compact, simple boundary"},
	{"GR", "mainland plus a dense island field"},
	{"CH", "small and compact, intricate border"},
	{"IS", "single island, rounded"},
}

// DemoReport says what was built and where it went. The path is reported whether
// or not a browser was opened, because the headless path has to be useful rather
// than merely non-broken.
type DemoReport struct {
	Path     string   `json:"path"`
	Entities []string `json:"entities"`
	Styles   []string `json:"styles"`
	Bytes    int      `json:"bytes"`
	Opened   bool     `json:"opened"`
}

func newDemoCommand(flags *globalFlags) *cobra.Command {
	var out string
	var open bool
	cmd := &cobra.Command{
		Use:   "demo",
		Short: "Build an interactive showcase page and open it",
		Long: strings.TrimSpace(`
Build one self-contained page showing ten entities chosen for shape diversity,
with live selectors for style, palette, pattern and theme, and open it in the
default browser.

Every map on the page is generated ONCE. Styles, colours and patterns are then
switched in the browser, because themed-inline delivery emits each token as a
CSS custom property with the resolved value as its fallback. That is the whole
thesis made visible: geometry carries no presentation, so presentation can
change without regenerating anything.

Like preview, it is a file rather than a server — nothing to run, no port, no
process left behind. The browser is not opened when stdout is not a terminal,
when --json is set, or with --open=false, so an agent gets the path instead of a
window it cannot see.`),
		Args: exactArgs(0, CLIName+" demo [--out <file>] [--open=false] [--json]"),
		RunE: func(cmd *cobra.Command, args []string) error {
			vocab, corpus, err := vocabulary()
			if err != nil {
				return err
			}

			page, rendered, err := buildDemo(corpus, vocab)
			if err != nil {
				return err
			}
			if err := os.WriteFile(out, []byte(page), 0o644); err != nil {
				return failf(ExitFilesystem, "write_failed", "%s: %v", out, err)
			}

			// Three independent reasons not to open, any one of which is
			// sufficient. A browser window is the one side effect a caller
			// cannot undo, so the default is permissive only for a human at a
			// terminal who asked for a demo.
			opened := false
			if open && !flags.json && stdoutIsTerminal() {
				opened = openInBrowser(out) == nil
			}

			report := DemoReport{Path: out, Entities: rendered, Styles: config.StyleNames, Bytes: len(page), Opened: opened}
			human := fmt.Sprintf("wrote %s (%d bytes)\n%d entities, every style switchable in the page", out, len(page), len(rendered))
			if !opened {
				human += "\nopen it yourself: " + out
			}
			return emit(cmd, flags, report, human)
		},
	}
	cmd.Flags().StringVar(&out, "out", "country-map-demo.html", "page to write")
	cmd.Flags().BoolVar(&open, "open", true, "open the page in the default browser")
	return cmd
}

// stdoutIsTerminal reports whether output is going to a terminal rather than to
// a pipe, a file or a test harness.
//
// It reads the mode bit from the standard library rather than taking a
// dependency on x/term: one bit is all this decision needs, and a new direct
// module requirement would redden the diagnostic identity record that hashes the
// whole module graph.
func stdoutIsTerminal() bool {
	info, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeCharDevice != 0
}

// openInBrowser hands the file to the desktop's own opener. There is no browser
// detection here on purpose: which browser is default is the operating system's
// answer to give, and every platform already has one command that asks it.
func openInBrowser(path string) error {
	var command string
	var args []string
	switch runtime.GOOS {
	case "darwin":
		command, args = "open", []string{path}
	case "windows":
		command, args = "rundll32", []string{"url.dll,FileProtocolHandler", path}
	default:
		command, args = "xdg-open", []string{path}
	}
	return exec.Command(command, args...).Start()
}

func buildDemo(corpus *catalog.Corpus, vocab config.Vocabulary) (string, []string, error) {
	var cards strings.Builder
	var rendered []string

	// One generation per entity, in themed-inline, under the style whose
	// fallbacks are the most neutral starting point. Every other style is
	// reachable from the page by setting the custom properties, so generating
	// the style matrix here would produce five copies of one geometry.
	themed := "themed-inline"
	filled := "filled"
	// Markers are generated ON and hidden by CSS, because a marker is geometry —
	// a projected coordinate the page cannot compute — while its visibility is
	// presentation. Generating them off would make the page's marker toggle a
	// regeneration, which is the one thing this page exists to avoid.
	capital := "capital"

	for _, entry := range demoEntities {
		iso := entry.ISO
		name := entityName(corpus, iso)

		resolved, err := config.Resolve(nil, iso, config.Settings{
			Delivery: &themed,
			Style:    &filled,
			Marker:   &config.Marker{Mode: &capital},
		})
		if err != nil {
			return "", nil, configError(err)
		}
		if err := config.ValidateResolved(iso, resolved, vocab); err != nil {
			return "", nil, configError(err)
		}
		request, err := render.GeometryRequest(corpus, iso, resolved.Settings)
		if err != nil {
			return "", nil, failf(ExitRender, "request_unbuildable", "%v", err)
		}
		result, err := geometry.Generate(request)
		if err != nil {
			if geometry.IsNoArtifact(err) {
				// A typed absence is a first-class outcome, not a failure to
				// hide: the page says so rather than dropping the card.
				fmt.Fprintf(&cards, `<article class="card empty" data-iso="%s"><div class="art">no artifact</div><footer><b>%s</b><span>%s</span></footer></article>`,
					iso, escapeHTML(name), iso)
				continue
			}
			return "", nil, failf(ExitRender, "generation_failed", "%s: %v", iso, err)
		}
		document, err := render.SVG(result, iso, name, resolved.Settings)
		if err != nil {
			return "", nil, failf(ExitRender, "serialization_failed", "%s: %v", iso, err)
		}

		fmt.Fprintf(&cards,
			`<article class="card" data-iso="%s" tabindex="0" role="button" aria-label="%s, open full width"><div class="art">%s</div><footer><b>%s</b><span>%s — %s</span></footer></article>`,
			iso, escapeHTML(name), document, escapeHTML(name), iso, escapeHTML(entry.Why))
		rendered = append(rendered, iso)
	}

	page := strings.ReplaceAll(demoTemplate, "{{CARDS}}", cards.String())
	page = strings.ReplaceAll(page, "{{COUNT}}", fmt.Sprintf("%d", len(rendered)))
	return page, rendered, nil
}

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/ydnikolaev/country-map-svg-generator/internal/catalog"
	"github.com/ydnikolaev/country-map-svg-generator/internal/config"
	"github.com/ydnikolaev/country-map-svg-generator/internal/geometry"
	"github.com/ydnikolaev/country-map-svg-generator/internal/render"
)

// PreviewReport says what was written and what to do with it.
type PreviewReport struct {
	Path     string   `json:"path"`
	Entities []string `json:"entities"`
	Styles   []string `json:"styles"`
	Bytes    int      `json:"bytes"`
}

func newPreviewCommand(flags *globalFlags) *cobra.Command {
	local := &configFlags{}
	var out string
	cmd := &cobra.Command{
		Use:   "preview",
		Short: "Write one self-contained page to look at the output",
		Long: strings.TrimSpace(`
Write a single self-contained HTML page showing the resolved configuration
rendered in every style, and open it yourself.

It is a file, not a server. The specification is explicit that preview must not
become a hidden web server requirement, so this writes one page with the assets
inlined and nothing to run — no localhost, no port, no process left behind.

Looking at output is not optional. Both defects this pipeline shipped and then
fixed — one country rendered as a blob, another as a speck in an empty card —
were invisible to every metric and were found by rendering and looking.`),
		Args: exactArgs(0, CLIName+" preview --iso <code> [--out <file>] [--config <path>] [--json]"),
		RunE: func(cmd *cobra.Command, args []string) error {
			doc, err := local.loadDocument()
			if err != nil {
				return err
			}
			vocab, corpus, err := vocabulary()
			if err != nil {
				return err
			}
			if doc != nil {
				if err := config.ValidateDocument(doc, vocab); err != nil {
					return configError(err)
				}
			}
			isos, err := selection(local, corpus, vocab)
			if err != nil {
				return err
			}
			if len(isos) > previewLimit {
				return failf(ExitUsage, "selection_too_wide",
					"preview renders up to %d entities and this selects %d", previewLimit, len(isos)).
					withContext("hint", "pass --iso, or run generate for the whole catalog")
			}

			page, rendered, err := buildPreview(corpus, doc, local, vocab, isos)
			if err != nil {
				return err
			}
			if err := os.WriteFile(out, []byte(page), 0o644); err != nil {
				return failf(ExitFilesystem, "write_failed", "%s: %v", out, err)
			}

			report := PreviewReport{Path: out, Entities: rendered, Styles: config.StyleNames, Bytes: len(page)}
			return emit(cmd, flags, report, fmt.Sprintf("wrote %s (%d bytes)\nopen it to look at %s in %d styles",
				out, len(page), strings.Join(rendered, ", "), len(config.StyleNames)))
		},
	}
	local.bind(cmd)
	cmd.Flags().StringVar(&out, "out", "country-map-preview.html", "page to write")
	return cmd
}

// previewLimit keeps preview a thing a person looks at. A page with the whole
// catalog inlined is a different artifact with a different purpose, and P4 owns
// it.
const previewLimit = 12

func buildPreview(
	corpus *catalog.Corpus, doc *config.Document, local *configFlags,
	vocab config.Vocabulary, isos []string,
) (string, []string, error) {
	var body strings.Builder
	var rendered []string

	for _, iso := range isos {
		name := entityName(corpus, iso)
		fmt.Fprintf(&body, "<section><h2>%s <small>%s</small></h2><div class=row>", escapeHTML(name), iso)
		for _, style := range config.StyleNames {
			overrides := local.overrides()
			overrides.Style = &style
			resolved, err := config.Resolve(doc, iso, overrides)
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
					fmt.Fprintf(&body, `<figure><div class="cell empty">no artifact</div><figcaption>%s</figcaption></figure>`, style)
					continue
				}
				return "", nil, failf(ExitRender, "generation_failed", "%s: %v", iso, err)
			}
			document, err := render.SVG(result, iso, name, resolved.Settings)
			if err != nil {
				return "", nil, failf(ExitRender, "serialization_failed", "%v", err)
			}
			// The page inlines exactly what generate would write, so what is
			// looked at is the artifact rather than a rendering of it.
			fmt.Fprintf(&body, `<figure><div class=cell>%s</div><figcaption>%s <span>%d B</span></figcaption></figure>`,
				document, style, len(document))
		}
		body.WriteString("</div></section>")
		rendered = append(rendered, iso)
	}
	return previewPage(body.String()), rendered, nil
}

// previewPage wraps the assets in a page that works offline in both colour
// schemes. The styling is deliberately minimal: this is a place to look at
// silhouettes, and anything more would be presentation competing with them.
func previewPage(body string) string {
	return `<!doctype html><html lang=en><head><meta charset=utf-8>
<meta name=viewport content="width=device-width,initial-scale=1">
<title>country-map preview</title><style>
:root{color-scheme:light dark;--bg:#fff;--ink:#111;--muted:#666;--line:#e5e5e5}
@media (prefers-color-scheme:dark){:root{--bg:#111;--ink:#eee;--muted:#999;--line:#2a2a2a}}
body{margin:0;padding:2rem;background:var(--bg);color:var(--ink);
font:14px/1.5 ui-sans-serif,system-ui,-apple-system,sans-serif}
h1{font-size:1.1rem;font-weight:600;margin:0 0 1.5rem}
h2{font-size:1rem;font-weight:600;margin:2rem 0 .75rem}
h2 small{color:var(--muted);font-weight:400;margin-left:.5rem}
.row{display:flex;flex-wrap:wrap;gap:1rem}
figure{margin:0}
.cell{width:160px;height:160px;display:grid;place-items:center;
border:1px solid var(--line);border-radius:6px;padding:8px}
.cell svg{max-width:100%;max-height:100%;display:block}
.cell.empty{color:var(--muted);font-size:12px}
figcaption{margin-top:.4rem;font-size:12px;color:var(--muted)}
figcaption span{opacity:.7;margin-left:.35rem}
</style></head><body>
<h1>country-map preview — every style, rendered offline</h1>
` + body + `</body></html>`
}

func escapeHTML(value string) string {
	replacer := strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;", `"`, "&quot;")
	return replacer.Replace(value)
}

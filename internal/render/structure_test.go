package render

import (
	"strings"
	"testing"

	"github.com/ydnikolaev/country-map-svg-generator/internal/config"
)

// This file is VAL-4: "parse emitted XML/SVG and inject forbidden
// nodes/attributes/references — every forbidden mutation fails".
//
// The mutations are applied to real emitted output rather than to a hand-written
// fixture. A gate that only ever sees markup written to be caught proves that
// the gate can catch that markup, not that it guards what the tool produces.

// mutate rewrites an emitted document by replacing the first occurrence of a
// marker. Each mutation below is something a defect, a careless template edit or
// a hostile config value would actually introduce.
func mutate(t *testing.T, document []byte, old, new string) []byte {
	t.Helper()
	text := string(document)
	if !strings.Contains(text, old) {
		t.Fatalf("the mutation anchor %q is not in the output, so the mutation tests nothing:\n%s", old, text)
	}
	return []byte(strings.Replace(text, old, new, 1))
}

func emitted(t *testing.T, delivery string) ([]byte, string) {
	t.Helper()
	result, iso := sampleResult(t)
	settings := settingsFor(t, "bold-soft", delivery, config.Settings{
		Marker: &config.Marker{Mode: strptr("all-capitals")},
	})
	document, err := SVG(result, iso, "Test Entity", settings)
	if err != nil {
		t.Fatal(err)
	}
	return document, iso
}

// TestEveryForbiddenMutationFails is the gate's teeth. Every case is a way REQ-9
// can be violated, and each must be rejected — a gate that misses one is a gate
// that would have shipped that one.
func TestEveryForbiddenMutationFails(t *testing.T) {
	document, _ := emitted(t, "themed-inline")

	// The unmutated asset must pass, or every rejection below proves nothing.
	if err := Structure(document); err != nil {
		t.Fatalf("real output is rejected, so the mutations below are meaningless: %v", err)
	}

	mutations := map[string][2]string{
		"a script element":       {"<path", "<script>alert(1)</script><path"},
		"a foreignObject":        {"<path", "<foreignObject/><path"},
		"an image element":       {"<path", `<image href="x.png"/><path`},
		"a use element":          {"<path", `<use href="#x"/><path`},
		"a style element":        {"<path", "<style>*{color:red}</style><path"},
		"an event handler":       {"<path ", `<path onload="alert(1)" `},
		"an inline style":        {"<path ", `<path style="fill:red" `},
		"an xlink reference":     {"<path ", `<path xlink:href="http://x/y" `},
		"an external image href": {"<path ", `<path href="https://cdn.example/x.svg" `},
		"a data URI":             {"<path ", `<path fill="data:image/png;base64,AAA" `},
		"a javascript URL":       {"<path ", `<path fill="javascript:alert(1)" `},
		"a remote url() paint":   {"<path ", `<path fill="url(https://cdn.example/x#g)" `},
		"a protocol-relative":    {"<path ", `<path fill="url(//cdn.example/x#g)" `},
		"an editor comment":      {"<path", "<!-- Created with an editor --><path"},
		"a stylesheet PI":        {"<svg", `<?xml-stylesheet href="x.css"?><svg`},
	}
	for name, mutation := range mutations {
		mutated := mutate(t, document, mutation[0], mutation[1])
		if err := Structure(mutated); err == nil {
			t.Errorf("%s was accepted:\n%s", name, mutated)
		}
	}
}

// TestTheGateDoesNotRejectWhatItShould keeps the allowlist from being so tight
// that it refuses the tool's own legitimate output. A gate nobody can satisfy
// gets deleted, and then nothing is guarded at all.
func TestTheGateDoesNotRejectWhatItShould(t *testing.T) {
	for _, delivery := range config.Deliveries {
		for _, style := range config.StyleNames {
			result, iso := sampleResult(t)
			settings := settingsFor(t, style, delivery, config.Settings{
				Marker: &config.Marker{Mode: strptr("capital")},
			})
			document, err := SVG(result, iso, "Test Entity", settings)
			if err != nil {
				t.Fatalf("%s/%s: %v", style, delivery, err)
			}
			if err := Structure(document); err != nil {
				t.Errorf("%s/%s: legitimate output rejected: %v", style, delivery, err)
			}
		}
	}
	// The namespace declaration is an absolute URL and must stay legal: it names
	// the language rather than fetching anything.
	document, _ := emitted(t, "standalone")
	if !strings.Contains(string(document), "http://www.w3.org/2000/svg") {
		t.Fatal("the namespace declaration is missing, so the exemption below is untested")
	}
	if err := Structure(document); err != nil {
		t.Errorf("the namespace declaration was treated as an external reference: %v", err)
	}
	// A local fragment reference is REQ-5's advanced escape hatch and stays legal.
	local := mutate(t, document, `fill="currentColor"`, `fill="url(#brand)"`)
	if err := Structure(local); err != nil {
		t.Errorf("a local fragment reference was refused: %v", err)
	}
}

// TestMalformedOutputIsCaught covers the failure that would otherwise reach a
// browser as a blank box: markup that is not well-formed at all.
func TestMalformedOutputIsCaught(t *testing.T) {
	document, _ := emitted(t, "standalone")
	for name, mutated := range map[string][]byte{
		"unclosed root":     []byte(strings.TrimSuffix(string(document), "</svg>")),
		"stray open angle":  mutate(t, document, "<path", "<<path"),
		"unbalanced quotes": mutate(t, document, `fill="currentColor"`, `fill="currentColor`),
		"empty":             {},
	} {
		if err := Structure(mutated); err == nil {
			t.Errorf("%s was accepted", name)
		}
	}
}

// TestEscapingIsDefenceInDepth checks the second line, not the first. The config
// layer already refuses markup characters in every author-supplied token, so
// this asks what happens if one reached the serializer anyway — which is what
// makes the answer "it cannot close an attribute" rather than "it never gets
// there".
func TestEscapingIsDefenceInDepth(t *testing.T) {
	var out strings.Builder
	e := element{name: "path"}
	e.set("d", `M0 0`)
	e.set("data-hostile", `" onload="alert(1)`)
	e.selfClose(&out)

	if strings.Contains(out.String(), `onload="alert`) {
		t.Fatalf("an attribute value escaped its quotes:\n%s", out.String())
	}
	wrapped := "<svg xmlns=\"http://www.w3.org/2000/svg\">" + out.String() + "</svg>"
	if err := Structure([]byte(wrapped)); err != nil {
		t.Errorf("escaping produced markup the parser rejects: %v", err)
	}

	// Every character that can terminate an attribute or an element is escaped.
	var all strings.Builder
	writeEscaped(&all, "&<>\"'\n\r\t")
	if got := all.String(); got != "&amp;&lt;&gt;&quot;&apos;&#xA;&#xD;&#x9;" {
		t.Errorf("escaping = %q", got)
	}
}

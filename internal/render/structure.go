package render

import (
	"encoding/xml"
	"fmt"
	"io"
	"sort"
	"strings"
)

// This file is REQ-9 as an executable check rather than as a promise. Every
// asset the CLI publishes goes through Structure before it is written, so a
// forbidden construct fails the batch instead of reaching a page.
//
// The parse is done with encoding/xml rather than by searching the text. A
// search for "<script" is defeated by any of the dozen ways markup can spell the
// same thing; a parser sees the element regardless of how it was written.

// AllowedElements is the complete set of elements an emitted asset may contain.
// It is an allowlist rather than a denylist on purpose: a denylist has to
// enumerate every dangerous element that exists now and every one added later,
// and it is wrong the first time either list changes.
var AllowedElements = []string{"svg", "path", "circle", "g", "title", "desc"}

// ForbiddenAttributePrefixes catch whole families rather than named members.
// `on*` is every event handler there is or will be; `xlink:` is the legacy
// external-reference surface.
var ForbiddenAttributePrefixes = []string{"on", "xlink:"}

// ForbiddenSchemes are the ways a value can reach outside the document. REQ-9
// forbids external URLs and embedded raster data; `data:` covers the second and
// the rest cover the first.
var ForbiddenSchemes = []string{"http://", "https://", "//", "data:", "javascript:", "file:", "blob:"}

// StructureError names every problem in one asset. All of them are reported
// together: a caller fixing generated output wants the list.
type StructureError struct {
	Problems []string
}

func (e *StructureError) Error() string {
	return "structure: " + strings.Join(e.Problems, "; ")
}

// Structure parses an emitted asset and reports every way it violates REQ-9.
func Structure(document []byte) error {
	var problems []string
	report := func(format string, args ...any) {
		problems = append(problems, fmt.Sprintf(format, args...))
	}

	decoder := xml.NewDecoder(strings.NewReader(string(document)))
	depth, elements := 0, 0
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return &StructureError{Problems: []string{"is not well-formed XML: " + err.Error()}}
		}

		switch typed := token.(type) {
		case xml.StartElement:
			depth++
			elements++
			checkElement(typed, depth, report)
		case xml.EndElement:
			depth--
		case xml.ProcInst:
			// A processing instruction is where a stylesheet reference or an
			// editor's own directive hides.
			report("carries a processing instruction <?%s?>", typed.Target)
		case xml.Directive:
			report("carries a document directive, which no generated asset needs")
		case xml.Comment:
			// Editor metadata (REQ-9) is overwhelmingly comments. There is also no
			// reason for a generated asset to explain itself in bytes the browser
			// downloads.
			report("carries a comment, which is editor metadata rather than output")
		}
	}
	if elements == 0 {
		report("contains no elements")
	}

	if len(problems) == 0 {
		return nil
	}
	sort.Strings(problems)
	return &StructureError{Problems: problems}
}

func checkElement(start xml.StartElement, depth int, report func(string, ...any)) {
	name := start.Name.Local
	if !containsString(AllowedElements, name) {
		report("contains a <%s> element, which is not in the allowed set (%s)", name, strings.Join(AllowedElements, ", "))
	}
	if depth == 1 && name != "svg" {
		report("has <%s> as its root element rather than <svg>", name)
	}

	for _, attr := range start.Attr {
		attrName := attr.Name.Local
		if attr.Name.Space != "" {
			attrName = attr.Name.Space + ":" + attrName
		}
		lowered := strings.ToLower(attrName)
		for _, prefix := range ForbiddenAttributePrefixes {
			if strings.HasPrefix(lowered, prefix) {
				report("carries the attribute %q, which matches the forbidden prefix %q", attrName, prefix)
			}
		}
		// `style` is not forbidden as such, but an inline style block is where
		// presentation escapes the token vocabulary, and REQ-4's "no renderer
		// forks" depends on it not existing.
		if lowered == "style" {
			report("carries an inline style attribute; presentation belongs in tokens")
		}
		checkValue(attrName, attr.Value, report)
	}
}

func checkValue(name, value string, report func(string, ...any)) {
	lowered := strings.ToLower(value)
	for _, scheme := range ForbiddenSchemes {
		if strings.Contains(lowered, scheme) {
			// The namespace declaration is the one legitimate absolute URL in an
			// SVG: it names the language, it is never fetched.
			if name == "xmlns" && value == "http://www.w3.org/2000/svg" {
				continue
			}
			report("attribute %q references %q, which reaches outside the document", name, scheme)
		}
	}
	if index := strings.Index(lowered, "url("); index >= 0 && !strings.HasPrefix(lowered[index:], "url(#") {
		report("attribute %q uses url() with something other than a local fragment", name)
	}
}

func containsString(set []string, value string) bool {
	for _, item := range set {
		if item == value {
			return true
		}
	}
	return false
}

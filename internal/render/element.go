package render

import "strings"

// element is the whole markup writer: a tag name and attributes in the order
// they were set. It is deliberately not a general XML library — REQ-9 asks for
// deterministic output with no redundant markup, and a writer that can only emit
// the four element shapes this package uses is easier to keep honest than one
// that can emit anything.
//
// Attribute order is insertion order, never map order, because ARCH-INV-4 makes
// byte-identical output a contract.
type element struct {
	name  string
	attrs []attribute
}

type attribute struct{ name, value string }

func (e *element) set(name, value string) {
	for i := range e.attrs {
		if e.attrs[i].name == name {
			e.attrs[i].value = value
			return
		}
	}
	e.attrs = append(e.attrs, attribute{name: name, value: value})
}

func (e *element) open(out *strings.Builder) {
	out.WriteByte('<')
	out.WriteString(e.name)
	e.writeAttributes(out)
	out.WriteByte('>')
}

func (e *element) close(out *strings.Builder) {
	out.WriteString("</")
	out.WriteString(e.name)
	out.WriteByte('>')
}

func (e *element) selfClose(out *strings.Builder) {
	out.WriteByte('<')
	out.WriteString(e.name)
	e.writeAttributes(out)
	out.WriteString("/>")
}

func (e *element) writeAttributes(out *strings.Builder) {
	for _, attr := range e.attrs {
		out.WriteByte(' ')
		out.WriteString(attr.name)
		out.WriteString(`="`)
		writeEscaped(out, attr.value)
		out.WriteByte('"')
	}
}

// writeEscaped is defence in depth, not the primary guard.
//
// The configuration layer already refuses markup characters and external
// references in every author-supplied token, which is where a bad value gets a
// diagnostic naming the key. This is the second line: even if a value reached
// here, it cannot close an attribute or open an element.
//
// It is written out rather than delegated to xml.EscapeText because that
// function is specified for character data. Its behaviour on the quote
// characters is adjacent to what an attribute needs rather than identical to it,
// and "adjacent" is not a contract to rest markup safety on. The five XML
// entities plus the three whitespace characters that must not survive raw in an
// attribute value are exactly this, and no more.
func writeEscaped(out *strings.Builder, value string) {
	for _, r := range value {
		switch r {
		case '&':
			out.WriteString("&amp;")
		case '<':
			out.WriteString("&lt;")
		case '>':
			out.WriteString("&gt;")
		case '"':
			out.WriteString("&quot;")
		case '\'':
			out.WriteString("&apos;")
		case '\n':
			out.WriteString("&#xA;")
		case '\r':
			out.WriteString("&#xD;")
		case '\t':
			out.WriteString("&#x9;")
		default:
			out.WriteRune(r)
		}
	}
}

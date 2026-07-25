package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

// EnvelopeSchema versions the machine-readable output contract (CTR-004). It is
// a public compatibility surface: fields may be added within v1, never removed
// or repurposed.
const EnvelopeSchema = "country-map.cli/v1"

// ExitClass is the typed failure taxonomy REQ-12 requires. The numeric codes are
// part of CTR-004 — a build pipeline branches on them — so they are frozen at
// v1 and appended to, never renumbered.
type ExitClass int

const (
	ExitOK ExitClass = 0
	// ExitUsage covers flags, arguments and command selection: the caller asked
	// for something the CLI cannot parse. Distinct from ExitConfig because the
	// remediation is different — fix the invocation, not the document.
	ExitUsage ExitClass = 1
	// ExitConfig covers a document that parsed but does not satisfy the schema:
	// unknown fields, cycles, missing presets, impossible combinations.
	ExitConfig ExitClass = 2
	// ExitData covers the corpus side: an unknown ISO code, an entity with no
	// geometry, a corpus that does not match the binary.
	ExitData ExitClass = 3
	// ExitRender covers a failure inside geometry or serialization for input
	// that was otherwise legal.
	ExitRender ExitClass = 4
	// ExitValidation covers emitted output that failed its own structural gate:
	// a forbidden node, an external reference, malformed XML.
	ExitValidation ExitClass = 5
	// ExitBudget covers a byte budget exceeded. Separate from ExitValidation
	// because it is the one failure a caller can fix by asking for less detail
	// rather than by fixing a defect.
	ExitBudget ExitClass = 6
	// ExitFilesystem covers staging, publication and path resolution.
	ExitFilesystem ExitClass = 7
)

var exitClassNames = map[ExitClass]string{
	ExitOK:         "ok",
	ExitUsage:      "usage",
	ExitConfig:     "config",
	ExitData:       "data",
	ExitRender:     "render",
	ExitValidation: "validation",
	ExitBudget:     "budget",
	ExitFilesystem: "filesystem",
}

func (c ExitClass) String() string {
	if name, ok := exitClassNames[c]; ok {
		return name
	}
	return "unknown"
}

// CLIError carries an exit class alongside the message. Context is free-form
// key/value remediation detail — the "name remediation context" half of REQ-12 —
// and appears verbatim in the JSON envelope.
type CLIError struct {
	Class   ExitClass
	Code    string
	Message string
	Context map[string]string
	Err     error
}

func (e *CLIError) Error() string { return e.Message }
func (e *CLIError) Unwrap() error { return e.Err }

func failf(class ExitClass, code, format string, args ...any) *CLIError {
	return &CLIError{Class: class, Code: code, Message: fmt.Sprintf(format, args...)}
}

// withContext attaches remediation detail without a second constructor.
func (e *CLIError) withContext(pairs ...string) *CLIError {
	if len(pairs)%2 != 0 {
		panic("withContext needs key/value pairs")
	}
	if e.Context == nil {
		e.Context = map[string]string{}
	}
	for i := 0; i < len(pairs); i += 2 {
		e.Context[pairs[i]] = pairs[i+1]
	}
	return e
}

// classify maps any error onto the taxonomy. An error that never passed through
// failf is a defect rather than a user mistake, so it lands on ExitRender with
// an explicit code instead of being silently reported as usage.
func classify(err error) *CLIError {
	if err == nil {
		return nil
	}
	var cliErr *CLIError
	if errors.As(err, &cliErr) {
		return cliErr
	}
	return &CLIError{Class: ExitRender, Code: "unclassified", Message: err.Error(), Err: err}
}

// Envelope is the stable --json shape for both outcomes. Data and Error are
// mutually exclusive; Status names which one is populated so an agent can branch
// without probing for nulls.
type Envelope struct {
	Schema  string         `json:"schema"`
	Command string         `json:"command"`
	Status  string         `json:"status"`
	Exit    int            `json:"exit"`
	Data    any            `json:"data,omitempty"`
	Error   *EnvelopeError `json:"error,omitempty"`
}

type EnvelopeError struct {
	Class   string            `json:"class"`
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Context map[string]string `json:"context,omitempty"`
}

func writeEnvelope(w io.Writer, env Envelope) error {
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	return encoder.Encode(env)
}

func successEnvelope(command string, data any) Envelope {
	return Envelope{Schema: EnvelopeSchema, Command: command, Status: "ok", Exit: int(ExitOK), Data: data}
}

func errorEnvelope(command string, err *CLIError) Envelope {
	return Envelope{
		Schema: EnvelopeSchema, Command: command, Status: "error", Exit: int(err.Class),
		Error: &EnvelopeError{
			Class: err.Class.String(), Code: err.Code, Message: err.Message, Context: err.Context,
		},
	}
}

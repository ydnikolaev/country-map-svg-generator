package geometry

import "testing"

func TestSerializeRoundTrip(t *testing.T) {
	c := []Command{{Op: "M", Values: []float64{0, 0}}, {Op: "L", Values: []float64{.5, -.25}}, {Op: "L", Values: []float64{1, 1}}, {Op: "Q", Values: []float64{1.2, 1.3, 2, 2}}, {Op: "Z"}}
	path, err := SerializeCommands(c)
	if err != nil {
		t.Fatal(err)
	}
	got, err := ParsePath(path)
	if err != nil {
		t.Fatal(err)
	}
	if !commandsEqual(c, got) {
		t.Fatalf("%q did not round trip", path)
	}
}

func TestSerializeRoundsTinyNegativeToZero(t *testing.T) {
	commands := []Command{{Op: "M", Values: []float64{-0.001, -0.004}}, {Op: "L", Values: []float64{1, 1}}, {Op: "Z"}}
	path, err := SerializeCommands(commands)
	if err != nil {
		t.Fatal(err)
	}
	if path == "" {
		t.Fatal("tiny negative coordinates produced an empty path")
	}
}

func TestSerializeFallsBackWhenRelativeRoundingDrifts(t *testing.T) {
	commands := []Command{{Op: "M", Values: []float64{100, 100}}}
	for i := 1; i <= 5; i++ {
		commands = append(commands, Command{Op: "L", Values: []float64{100 + float64(i)*.004, 100}})
	}
	commands = append(commands, Command{Op: "Z"})
	path, err := SerializeCommands(commands)
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := ParsePath(path)
	if err != nil {
		t.Fatal(err)
	}
	if !commandsEqual(commands, parsed) {
		t.Fatalf("fallback path %q did not round trip", path)
	}
}

func TestSerializeRejectsAdversarialGrammar(t *testing.T) {
	for _, path := range []string{"M-0..1 2", "M 1", "X0 0", "M0,0L--1 2"} {
		if _, err := ParsePath(path); err == nil {
			t.Fatalf("accepted %q", path)
		}
	}
}

func TestSerializeMutationZeroByteFalseSuccessFails(t *testing.T) {
	if _, err := SerializeCommands(nil); err == nil {
		t.Fatal("empty command stream reported success")
	}
	if _, path, err := commandsAndPath(nil, 0); err == nil || path != "" {
		t.Fatalf("empty geometry path=%q err=%v", path, err)
	} else if pipelineErr, ok := err.(*PipelineError); !ok || pipelineErr.Code != ErrSerialization {
		t.Fatalf("error=%T %v, want typed serialization failure", err, err)
	}
}

func TestCanonicalizeQuantizationFailureDoesNotRetrySimplification(t *testing.T) {
	source := MultiPolygon{{{
		{0, 0},
		{.004, 0},
		{.004, .004},
		{0, 0},
	}}}
	got, err := canonicalizeShared(source, .01, nil)
	if err == nil || got != nil {
		t.Fatalf("canonicalization recovered from fixed-grid collapse: geometry=%v err=%v", got, err)
	}
	pipelineErr, ok := err.(*PipelineError)
	if !ok || pipelineErr.Code != ErrTopology || pipelineErr.Field != "quantization" {
		t.Fatalf("error=%T %v, want fixed-grid quantization topology failure", err, err)
	}
}

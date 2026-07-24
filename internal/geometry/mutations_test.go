package geometry

import (
	"errors"
	"testing"
)

func TestMutationTopologyBowTieFails(t *testing.T) {
	g := MultiPolygon{{{{0, 0}, {2, 2}, {0, 2}, {2, 0}, {0, 0}}}}
	err := validateTopology(g)
	var pe *PipelineError
	if !errors.As(err, &pe) || pe.Code != ErrTopology {
		t.Fatalf("bow tie passed: %v", err)
	}
}

func TestMutationEscapedHoleFails(t *testing.T) {
	g := MultiPolygon{{{{0, 0}, {4, 0}, {4, 4}, {0, 4}, {0, 0}}, {{5, 5}, {6, 5}, {6, 6}, {5, 6}, {5, 5}}}}
	if err := validateTopology(g); err == nil {
		t.Fatal("escaped hole passed")
	}
}

func TestMutationCollapsedRingFails(t *testing.T) {
	g := MultiPolygon{{{{0, 0}, {1, 0}, {2, 0}, {0, 0}}}}
	if err := validateTopology(g); err == nil {
		t.Fatal("collapsed ring passed")
	}
}

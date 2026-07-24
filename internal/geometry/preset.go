package geometry

import (
	_ "embed"
	"encoding/json"
	"fmt"
)

//go:embed presets/v1.json
var presetJSON []byte

type Preset struct {
	Name, Version  string
	Layout         Layout
	AdvisoryBytes  int
	MaxPathBytes   int
	ReferenceSizes []float64
}

type presetFile struct {
	Version       string        `json:"version"`
	QualityPolicy qualityPolicy `json:"quality_policy"`
	Presets       []struct {
		Name           string    `json:"name"`
		LongSide       float64   `json:"long_side"`
		Padding        float64   `json:"padding"`
		AdvisoryBytes  int       `json:"advisory_bytes"`
		MaxPathBytes   int       `json:"max_path_bytes"`
		ReferenceSizes []float64 `json:"reference_sizes"`
	} `json:"presets"`
}

type inverseScalePolicy struct {
	Floor, InverseRatio, Ceiling float64
}

func (p *inverseScalePolicy) UnmarshalJSON(raw []byte) error {
	var value struct {
		Floor        float64 `json:"floor"`
		InverseRatio float64 `json:"inverse_ratio"`
		Ceiling      float64 `json:"ceiling"`
	}
	if err := json.Unmarshal(raw, &value); err != nil {
		return err
	}
	*p = inverseScalePolicy(value)
	return nil
}

type relativeScalePolicy struct {
	FloorPX, RelativeRatio float64
}

func (p *relativeScalePolicy) UnmarshalJSON(raw []byte) error {
	var value struct {
		FloorPX       float64 `json:"floor_px"`
		RelativeRatio float64 `json:"relative_ratio"`
	}
	if err := json.Unmarshal(raw, &value); err != nil {
		return err
	}
	*p = relativeScalePolicy(value)
	return nil
}

type qualityPolicy struct {
	Version        string
	Flatness       inverseScalePolicy
	Simplification relativeScalePolicy
	Softening      inverseScalePolicy
	MinimumArea    inverseScalePolicy `json:"minimum_area"`
	Quantization   float64
}

var embeddedPresetFile = mustParsePresetFile(presetJSON)

func mustParsePresetFile(raw []byte) presetFile {
	f, err := parsePresetFile(raw)
	if err != nil {
		panic("invalid embedded geometry preset policy: " + err.Error())
	}
	return f
}

func parsePresetFile(raw []byte) (presetFile, error) {
	var f presetFile
	if err := json.Unmarshal(raw, &f); err != nil {
		return presetFile{}, err
	}
	if f.Version == "" || f.QualityPolicy.Version == "" {
		return presetFile{}, fmt.Errorf("missing preset or quality policy version")
	}
	p := f.QualityPolicy
	if p.Simplification.FloorPX != .80 ||
		p.Simplification.RelativeRatio < .010 || p.Simplification.RelativeRatio > .012 {
		return presetFile{}, fmt.Errorf("simplification policy outside accepted range")
	}
	for name, value := range map[string]inverseScalePolicy{
		"flatness": p.Flatness, "softening": p.Softening, "minimum_area": p.MinimumArea,
	} {
		if value.Floor <= 0 || value.InverseRatio <= 0 || value.Ceiling < value.Floor {
			return presetFile{}, fmt.Errorf("invalid %s policy", name)
		}
	}
	if p.Quantization != .01 {
		return presetFile{}, fmt.Errorf("canonical quantization must equal 0.01")
	}
	for _, preset := range f.Presets {
		if preset.Name == "" || preset.LongSide <= 0 || preset.Padding < 0 {
			return presetFile{}, fmt.Errorf("invalid preset %q", preset.Name)
		}
	}
	return f, nil
}

func Presets() (map[string]Preset, error) {
	f := embeddedPresetFile
	out := make(map[string]Preset, len(f.Presets))
	for _, p := range f.Presets {
		if p.Name == "" || p.LongSide <= 0 || p.Padding < 0 {
			return nil, fmt.Errorf("invalid preset %q", p.Name)
		}
		if _, ok := out[p.Name]; ok {
			return nil, fmt.Errorf("duplicate preset %q", p.Name)
		}
		out[p.Name] = Preset{Name: p.Name, Version: f.Version, Layout: Layout{Mode: LayoutTight, LongSide: p.LongSide, Padding: Insets{p.Padding, p.Padding, p.Padding, p.Padding}}, AdvisoryBytes: p.AdvisoryBytes, MaxPathBytes: p.MaxPathBytes, ReferenceSizes: p.ReferenceSizes}
	}
	return out, nil
}

func ApplyPreset(in Input) (Input, error) {
	if in.Preset == "" {
		return in, nil
	}
	all, err := Presets()
	if err != nil {
		return in, err
	}
	p, ok := all[in.Preset]
	if !ok {
		return in, fmt.Errorf("unknown preset %q", in.Preset)
	}
	if in.Layout.Mode == "" {
		in.Layout = p.Layout
	}
	if in.MaxPathBytes == 0 {
		in.MaxPathBytes = p.MaxPathBytes
	}
	return in, nil
}

func AutoQuality(d float64) Quality {
	p := embeddedPresetFile.QualityPolicy
	return Quality{
		Flatness:       clamp(p.Flatness.Floor+p.Flatness.InverseRatio/d, p.Flatness.Floor, p.Flatness.Ceiling),
		Simplification: mathMax(p.Simplification.FloorPX, p.Simplification.RelativeRatio*d),
		Softening:      clamp(p.Softening.Floor+p.Softening.InverseRatio/d, p.Softening.Floor, p.Softening.Ceiling),
		MinimumArea:    clamp(p.MinimumArea.Floor+p.MinimumArea.InverseRatio/d, p.MinimumArea.Floor, p.MinimumArea.Ceiling),
		Quantization:   p.Quantization, Auto: true,
	}
}

func mathMax(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func clamp(v, lo, hi float64) float64 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

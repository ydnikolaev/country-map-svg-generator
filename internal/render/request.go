// Package render owns BND-006: minimal SVG serialization and the style and pin
// contracts. This file is the one bridge between a resolved configuration and a
// geometry request; nothing else in the CLI constructs a geometry.Input.
package render

import (
	"fmt"

	"github.com/yuranikolaev/country-map-svg-generator/internal/catalog"
	"github.com/yuranikolaev/country-map-svg-generator/internal/config"
	"github.com/yuranikolaev/country-map-svg-generator/internal/geometry"
)

// GeometryRequest builds the geometry input for one entity from a resolved
// configuration.
//
// **The two axes cross here, and that is the whole reason this function
// exists in one place.** The configuration's `profile` (card, hero) becomes
// geometry's *Preset*, and the configuration's `boundary` (un, de_facto)
// becomes geometry's *Profile*. Transposed, the output would carry the wrong
// boundary posture — a DEC-002-class defect that satisfies every byte budget and
// every structural check, and is invisible to anything but a human looking at a
// disputed border.
//
// It resolves through geometry.InputFromCatalog rather than assembling an Input
// directly. That adapter refuses an unknown boundary outright, which matters
// because the pipeline's own de-facto test is a bare string comparison: an Input
// built by hand with the corpus manifest's hyphenated "de-facto" spelling would
// silently take the `un` branch.
func GeometryRequest(corpus *catalog.Corpus, iso string, settings config.Settings) (geometry.Input, error) {
	if settings.Profile == nil {
		return geometry.Input{}, fmt.Errorf("%s: profile is not resolved", iso)
	}
	if settings.Boundary == nil {
		return geometry.Input{}, fmt.Errorf("%s: boundary is not resolved", iso)
	}
	input, err := geometry.InputFromCatalog(corpus, iso, *settings.Boundary, *settings.Profile)
	if err != nil {
		return geometry.Input{}, err
	}
	layout, err := geometryLayout(settings.Layout)
	if err != nil {
		return geometry.Input{}, fmt.Errorf("%s: %w", iso, err)
	}
	input.Layout = layout
	// MaxPathBytes is deliberately left at zero. geometry.ApplyPreset fills it
	// from the preset's own budget, and restating a budget here would fork the
	// frozen 2200/7500 path and 2500/8000 file limits.
	return input, nil
}

// geometryLayout maps the configuration's layout onto geometry's.
//
// The one subtle rule: when the configuration expresses a mode but no sizing at
// all, the returned layout is left entirely zero. geometry.ApplyPreset fills a
// layout only when its mode is empty, so returning a mode with no size would
// suppress the preset's long side and produce a card sized by nothing.
func geometryLayout(layout *config.Layout) (geometry.Layout, error) {
	if layout == nil || layout.Mode == nil {
		return geometry.Layout{}, nil
	}
	padding := geometryPadding(layout.Padding)
	sized := layout.LongSide != nil || layout.MaxWidth != nil || layout.MaxHeight != nil ||
		layout.Width != nil || layout.Height != nil
	if !sized && layout.Padding == nil {
		return geometry.Layout{}, nil
	}

	switch *layout.Mode {
	case "tight":
		if !sized {
			// Padding alone cannot carry a layout past ApplyPreset either, so the
			// preset still has to supply the size; the padding is applied on top.
			return geometry.Layout{}, nil
		}
		out := geometry.Layout{Mode: geometry.LayoutTight, Padding: padding}
		if layout.LongSide != nil {
			out.LongSide = *layout.LongSide
		}
		if layout.MaxWidth != nil {
			out.MaxWidth = *layout.MaxWidth
		}
		if layout.MaxHeight != nil {
			out.MaxHeight = *layout.MaxHeight
		}
		return out, nil

	case "contain":
		if layout.Width == nil || layout.Height == nil {
			return geometry.Layout{}, fmt.Errorf("contain layout needs both width and height")
		}
		return geometry.Layout{
			Mode: geometry.LayoutContain, Width: *layout.Width, Height: *layout.Height, Padding: padding,
		}, nil
	}
	return geometry.Layout{}, fmt.Errorf("unknown layout mode %q", *layout.Mode)
}

func geometryPadding(padding *config.Padding) geometry.Insets {
	if padding == nil {
		return geometry.Insets{}
	}
	if padding.Uniform != nil {
		value := *padding.Uniform
		return geometry.Insets{Top: value, Right: value, Bottom: value, Left: value}
	}
	out := geometry.Insets{}
	if padding.Top != nil {
		out.Top = *padding.Top
	}
	if padding.Right != nil {
		out.Right = *padding.Right
	}
	if padding.Bottom != nil {
		out.Bottom = *padding.Bottom
	}
	if padding.Left != nil {
		out.Left = *padding.Left
	}
	return out
}

// Vocabulary builds the value sets the configuration layer validates against,
// enumerated from their owners rather than spelled out.
//
// The boundary set comes from geometry.AcceptedBoundaryProfiles and explicitly
// **not** from corpus.Manifest.Profiles: the manifest declares its list with a
// hyphen ("de-facto") while geometry accepts an underscore, and nothing
// translates between them, so building the enum from the manifest would offer
// authors a value the pipeline refuses.
func Vocabulary(corpus *catalog.Corpus) (config.Vocabulary, error) {
	presets, err := geometry.Presets()
	if err != nil {
		return config.Vocabulary{}, err
	}
	vocab := config.Vocabulary{
		Boundaries: append([]string(nil), geometry.AcceptedBoundaryProfiles...),
		CapitalIDs: map[string][]string{},
	}
	for name := range presets {
		vocab.Profiles = append(vocab.Profiles, name)
	}
	if corpus == nil {
		return vocab, nil
	}
	for _, entity := range corpus.Manifest.Entities {
		vocab.ISOCodes = append(vocab.ISOCodes, entity.Alpha2)
		if len(entity.Capitals) == 0 {
			continue
		}
		ids := make([]string, 0, len(entity.Capitals))
		for _, capital := range entity.Capitals {
			ids = append(ids, capital.ID)
		}
		vocab.CapitalIDs[entity.Alpha2] = ids
	}
	return vocab, nil
}

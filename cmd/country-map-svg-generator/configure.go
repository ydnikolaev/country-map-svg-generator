package main

import (
	"errors"
	"io/fs"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/ydnikolaev/country-map-svg-generator/internal/catalog"
	"github.com/ydnikolaev/country-map-svg-generator/internal/config"
	"github.com/ydnikolaev/country-map-svg-generator/internal/render"
)

// configFlags are the options every configuration-reading command shares. They
// live in one place so `validate`, `explain` and later `generate` cannot drift
// into accepting slightly different spellings of the same idea.
type configFlags struct {
	path     string
	iso      string
	profile  string
	boundary string
	style    string
	delivery string
}

func (f *configFlags) bind(cmd *cobra.Command) {
	cmd.Flags().StringVar(&f.path, "config", "", "path to a country-map configuration (.yaml, .yml or .json)")
	cmd.Flags().StringVar(&f.iso, "iso", "", "ISO alpha-2 code to resolve for")
	cmd.Flags().StringVar(&f.profile, "profile", "", "override the detail profile (card, hero)")
	cmd.Flags().StringVar(&f.boundary, "boundary", "", "override the boundary posture (un, de_facto)")
	cmd.Flags().StringVar(&f.style, "style", "", "override the visual style")
	cmd.Flags().StringVar(&f.delivery, "delivery", "", "override the delivery mode")
}

// overrides turns the selection flags into the strongest layer of the
// precedence chain. A flag that was not given contributes nothing, which is what
// keeps "not passed" from meaning "set to empty".
func (f *configFlags) overrides() config.Settings {
	settings := config.Settings{}
	for value, target := range map[string]**string{
		f.profile: &settings.Profile, f.boundary: &settings.Boundary,
		f.style: &settings.Style, f.delivery: &settings.Delivery,
	} {
		if value != "" {
			local := value
			*target = &local
		}
	}
	return settings
}

// loadDocument reads the configuration, or reports that there is none. A missing
// --config is not an error: the embedded defaults are a complete configuration,
// and refusing to run without a file would make the zero-config path — the one a
// first-time caller takes — impossible.
func (f *configFlags) loadDocument() (*config.Document, error) {
	if f.path == "" {
		return nil, nil
	}
	doc, err := config.DecodeFile(f.path)
	if err == nil {
		return doc, nil
	}
	switch {
	case errors.Is(err, fs.ErrNotExist):
		return nil, failf(ExitConfig, "config_missing", "no configuration at %s", f.path)
	case isConfigShapeError(err):
		// The typed error already carries the path, so it is not repeated in the
		// context map.
		return nil, failf(ExitConfig, "config_invalid", "%s", err.Error())
	default:
		return nil, failf(ExitFilesystem, "config_unreadable", "%s: %v", f.path, err)
	}
}

// isConfigShapeError separates "the document is wrong" from "the file could not
// be read", because REQ-12 gives those different exit classes and a caller
// branches on them: one is fixed by editing, the other by looking at the path.
//
// It asks by type, never by matching substrings of a message. The exit classes
// are a frozen public contract, and deriving one from another library's wording
// means a dependency upgrade can silently reclassify a failure.
func isConfigShapeError(err error) bool {
	var document *config.DocumentError
	if errors.As(err, &document) {
		return true
	}
	var unknown *config.UnknownKeyError
	if errors.As(err, &unknown) {
		return true
	}
	var problems config.ValidationErrors
	return errors.As(err, &problems)
}

// vocabulary builds the value sets validation checks against, enumerated from
// geometry and the corpus rather than spelled here.
func vocabulary() (config.Vocabulary, *catalog.Corpus, error) {
	corpus, err := catalog.Embedded()
	if err != nil {
		return config.Vocabulary{}, nil, failf(ExitData, "corpus_unreadable", "embedded corpus could not be decoded: %v", err)
	}
	vocab, err := render.Vocabulary(corpus)
	if err != nil {
		return config.Vocabulary{}, nil, failf(ExitData, "vocabulary_unavailable", "%v", err)
	}
	return vocab, corpus, nil
}

// selection is the set of entities a command operates on: the one named by
// --iso, or every entity the corpus carries.
func selection(flags *configFlags, corpus *catalog.Corpus, vocab config.Vocabulary) ([]string, error) {
	if flags.iso == "" {
		out := append([]string(nil), vocab.ISOCodes...)
		sort.Strings(out)
		return out, nil
	}
	iso := strings.ToUpper(flags.iso)
	for _, known := range vocab.ISOCodes {
		if known == iso {
			return []string{iso}, nil
		}
	}
	return nil, failf(ExitData, "unknown_entity", "no entity %q in corpus %s", flags.iso, short(corpus.Manifest.Identity)).
		withContext("hint", "run "+CLIName+" inspect --iso <code> to see what a code resolves to")
}

// configError converts the config layer's diagnostics into the CLI taxonomy,
// preserving every problem rather than reporting only the first.
func configError(err error) *CLIError {
	var problems config.ValidationErrors
	if errors.As(err, &problems) {
		cliErr := failf(ExitConfig, "invalid_configuration", "%s", problems.Error())
		for i, problem := range problems {
			if i >= 8 {
				// A long list is truncated in the context map but never in the
				// message: an agent reads the message, a human reads the context.
				break
			}
			cliErr = cliErr.withContext(problem.Path, problem.Message)
		}
		return cliErr
	}
	var unknown *config.UnknownKeyError
	if errors.As(err, &unknown) {
		return failf(ExitConfig, "unknown_field", "%s", unknown.Error()).
			withContext("known_fields", strings.Join(config.SchemaKeys(), ", "))
	}
	return classify(err)
}

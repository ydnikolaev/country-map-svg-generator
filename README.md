# country-map-svg-generator

Generate optimized country map SVGs from a versioned geometry corpus that is
embedded in the binary. No Node, Python, GDAL, network access or auxiliary data
files — one executable produces the whole catalog offline.

The primary operator is an AI agent: every command is non-interactive, every
option is a flag, and `--json` emits a stable envelope with a typed error class.

## Install

Download the archive for your platform from the release, verify it and put the
binary on your `PATH`:

```sh
tar -xzf country-map-svg-generator_<version>_<os>_<arch>.tar.gz
shasum -a 256 -c SHA256SUMS          # optional, from the same release
install country-map-svg-generator_*/country-map-svg-generator /usr/local/bin/
```

Or install from source (Go 1.26+):

```sh
go install github.com/ydnikolaev/country-map-svg-generator/cmd/country-map-svg-generator@latest
```

A `go install` build reports its version as `dev`; the corpus and algorithm
identities it prints are read from the embedded artifacts either way. To get a
stamped binary, clone and run `make build`, which writes
`dist/country-map-svg-generator`.

## Quickstart

```sh
country-map-svg-generator init                       # writes country-map.yaml
country-map-svg-generator validate --config country-map.yaml
country-map-svg-generator generate --config country-map.yaml --out dist
```

That writes one SVG per entity plus `dist/country-map.manifest.json`. The whole
catalog — 248 assets, ~570 KB — takes about four seconds.

Look at a single entity before committing to a batch:

```sh
country-map-svg-generator inspect --iso DE --config country-map.yaml
country-map-svg-generator preview --config country-map.yaml --out preview.html
```

## Commands

| Command | What it does |
| --- | --- |
| `init` | write a starter configuration |
| `validate` | check a configuration before it writes anything |
| `explain` | show the resolved configuration and where each value came from |
| `inspect` | report what one entity resolves to, without writing |
| `generate` | generate SVG assets and a manifest |
| `preview` | write one self-contained page to look at the output |
| `schema` | report every settable key and the values it accepts |
| `version` | report generator, corpus and algorithm identities |

### Discovering the surface

`schema` is the answer to "what may I set, and to what?". Every other command
tells you a value is wrong *after* you have written it; this one lists the
accepted values first.

```sh
country-map-svg-generator schema           # keys, vocabularies, styles, notes
country-map-svg-generator schema --json    # the same, machine-readable
```

It is derived from the definitions validation checks against — the key paths by
reflection over the schema, the vocabularies from the same variables the
validator compares to, and `profile` and `boundary` from the geometry presets
and the corpus at runtime. It cannot drift from what the binary accepts.

`explain` is its complement: `schema` says what a key *accepts*, `explain` says
what a key *resolved to* and which layer set it.

## Output

Each asset is a presentation-free SVG: geometry only, no fill, stroke or CSS.
Style it from your own stylesheet.

```html
<svg viewBox="0 0 111.59 144" class="country-map" data-country="DE"
     data-profile="card" data-boundary="un" aria-hidden="true" focusable="false">
  <path class="country-map__shape" d="…"/>
</svg>
```

```css
.country-map__shape { fill: var(--map-fill); stroke: var(--map-edge); }
```

Generation is transactional: the batch is built and structurally validated in
memory, then published atomically. A failure leaves the previous output exactly
as it was.

## Styling

Geometry is presentation-free by design — the pipeline never learns about
colour. That invariant is what makes everything below possible.

**Styles** are token bundles, not renderer modes. Run `schema` for the current
list with a sentence on what each is for; today: `outline`, `filled`,
`silhouette`, `bold-soft`, `ghost`.

**Tokens** are the individual knobs a style presets and a document may override:
`fill`, `fillOpacity`, `stroke`, `strokeOpacity`, `strokeWidth`, `lineCap`,
`lineJoin`, and `markerFill` / `markerStroke` / `markerRadius`.

**Delivery** decides how those values reach the page.

| Mode | Behaviour |
| --- | --- |
| `standalone` | visually complete with no host CSS |
| `themed-inline` | every value emitted as `var(--country-map-<token>, <resolved>)`, so the host page can restyle without regenerating |

**Markers** are capital points from the corpus registry, rendered as a circle
with its own class. `marker.mode` is `none` by default; `capital`,
`all-capitals` and `custom` turn them on. `custom` selects registry IDs rather
than coordinates, so a marker can never disagree with the corpus about where a
place is.

### Patterns, gradients and filters

`tokens.advanced` takes a **local fragment reference** — a pattern, gradient or
filter you define once on your own page. External URLs are refused.

```yaml
delivery: themed-inline
style: filled
tokens:
  advanced:
    pattern: "url(#hatch)"
```

```html
<svg width="0" height="0" style="position:absolute">
  <defs>
    <pattern id="hatch" width="6" height="6" patternUnits="userSpaceOnUse"
             patternTransform="rotate(45)">
      <line x1="0" y1="0" x2="0" y2="6" stroke="currentColor" stroke-width="1"/>
    </pattern>
  </defs>
</svg>
```

Every generated asset then fills from that one definition, at zero bytes per
asset. Swap the `<pattern>` for dots, change the angle, or restyle it from CSS
without regenerating anything.

`gradient` and `pattern` are both paint servers and both replace the fill, so
setting both is refused rather than silently resolved — name one. `filter` is
its own attribute and composes with whichever fill applies.

One limit worth stating plainly: a fragment reference points at defs **you**
supply. Under `standalone` the asset is otherwise self-contained, but a
`url(#id)` still needs that id present in the document it is inlined into. The
generator does not emit a `<defs>` block of its own; per-asset byte budgets are
frozen and a defs block would not fit inside them.

## Configuration

`country-map.yaml` only carries what differs from the preset it extends. Run
`explain` to see the resolved value of anything and the layer it came from.

```yaml
schema: country-map/v1
extends: site-default
profile: card          # card or hero — the byte budget and natural size
boundary: un           # un or de_facto
layout:
  mode: tight          # tight derives natural proportions; contain takes width/height
countries:
  US:
    profile: hero
```

Only the profile's own long side is served from the committed detail ladder.
Setting some other `longSide` falls back to full-detail source geometry, which is
over the byte ceiling for a large entity — set one only for a selection you have
generated successfully.

## Exit codes

`--json` callers branch on these; they are frozen at v1 and appended to, never
renumbered.

| Code | Class | Meaning |
| --- | --- | --- |
| 0 | ok | |
| 1 | usage | flags, arguments, command selection |
| 2 | config | the document parsed but violates the schema |
| 3 | data | unknown ISO, missing geometry, corpus mismatch |
| 4 | render | geometry or serialization failed on legal input |
| 5 | validation | emitted output failed its own structural gate |
| 6 | budget | a byte budget was exceeded — ask for less detail |
| 7 | filesystem | staging, publication, path resolution |

## Known limits

- **Scattered archipelagos render faint.** A card is fitted to everything it
  draws, so a state whose territory spans an ocean gets a frame far larger than
  any one island. Eighteen small island territories fill under 2 % of their
  frame, and **Portugal** (2.5 %, Azores and Madeira) is the one major country
  affected. Mainland-only framing is a deferred decision, not an oversight.
- **One entity has no card.** `UM` (United States Minor Outlying Islands) is
  reported as a typed skip with a reason rather than an empty file.
- `make check` runs serially and takes roughly fourteen minutes.

## Release

```sh
make release                     # VERSION from git describe
make release VERSION=v0.1.0
```

Writes `dist/*.tar.gz` for darwin and linux on arm64 and amd64, plus
`dist/SHA256SUMS`. Builds are `-trimpath` and CGO-free.

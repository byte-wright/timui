# timui

An immediate-mode terminal UI library for Go, modeled on Dear ImGui. The whole
UI is described every frame with plain function calls; widget state lives in the
caller, not in a retained widget tree.

```go
if tui.Button("ClickMe +") { count++ }
tui.Checkbox("Alpha", &checkedA)
```

Self-described pet project / experiment. ~2,500 LOC of non-test Go.

## Commands

- Build: `go build ./...`
- Test: `go test ./...`
- Vet / format: `go vet ./...`, `gofmt -l .`
- Run the demo: `go run ./cmd/example` (mouse-driven; Esc or Ctrl-C to quit)

Stdout/stderr are redirected to `stdout.log` / `stderr.log` in the working dir
while the demo runs (see [util/redirstd.go](util/redirstd.go)), because the TUI
owns the terminal. Check those files when debugging a crash.

Rendering is pinned by snapshot tests ([snapshot_test.go](snapshot_test.go)):
`snapshotBackend` accumulates cell diffs into a char grid and its string form is
checked with `expect`'s `ToBeSnapshot` against [testdata/](testdata/). On
mismatch the test writes a `.current` file next to the snapshot for diffing;
delete a snapshot file to regenerate it from the current rendering.

## Mental model

**Frame lifecycle.** The app calls widget functions, then [`Timui.Finish`](timui.go)
once per frame. `Finish` runs the deferred overlay queue, diffs the front buffer
against the back buffer, pushes only changed cells to the backend, swaps buffers,
resizes, clears the front buffer, and resets per-frame state. There is no
event-driven redraw; the caller drives a loop (the demo polls at ~30fps).

**Backend seam.** The core has no tcell dependency. Everything the renderer needs
is the 6-method `Backend` interface in [timui.go](timui.go); [tcell/](tcell/) is a
thin adapter. Keep new rendering/input concerns behind this interface.

**Area stack.** Layout is a stack of `mathi.Box2` regions. Widgets draw into
`CurrentArea` and advance a cursor via `moveCursor`. Splits, grids, rows,
columns, padding, and dialogs all work by pushing sub-areas — normally through
the closure-scoped `WithArea` / `WithAreaTranslation` / `WithClip`; the raw
`PushArea` / `PopArea` / `PushClip` / `PopClip` remain as low-level primitives
for custom widgets.

**Retained state pattern.** Stateful widgets (mouse, dropdown, draggable,
scrollarea) keep two maps — `last*` and `next*` — keyed by widget ID. Each frame
they read from `last`, write to `next`, and swap in their `finish` method. Follow
this same shape for any new stateful widget.

**IDs.** Widget identity is derived from the label plus the current ID scope stack
([internal/id.go](internal/id.go)). Identical labels in the same scope collide;
use `WithID` ([id.go](id.go)) to disambiguate (e.g. two "+" buttons in
different rows). Changing a widget's label changes its ID and drops its state.

**Deferred overlays.** Dropdowns and dialogs call `runAfter` so they paint on top
after the main pass. Z-order follows registration order in the `after` queue;
a deferred callback may itself call `runAfter` (dropdown inside a dialog) and
the nested overlay runs after — and thus above — its parent.

**Color / alpha.** Colors are packed `uint32` (`RGBColor` / `RGBAColor` in
[color.go](color.go)). The cell buffer blends alpha ([internal/screen.go](internal/screen.go)),
which is how the dialog dims its backdrop. A cell char of `0` means "leave the
existing glyph, only update color."

## Layout

- [split.go](split.go): `Split().Fixed(...).Factor(...)` computes ranges from a mix
  of fixed sizes and proportional factors. This is the tested core of layout.
- [grid.go](grid.go): bordered `Grid` whose `Rows` / `Columns` take one closure
  per split entry (`func(cell *GridCell)`) and draw divider lines between cells;
  nested splits go through the passed `GridCell`, whose area includes the
  surrounding border lines so divider glyphs merge into junctions. One
  axis-parameterized implementation (`gridSplit`) drives both directions.
- [rows.go](rows.go) / [columns.go](columns.go): borderless linear layout taking
  one cell closure per split entry (`Rows(split, cells ...func())`); panics on a
  cell/split count mismatch.

## Conventions

- No comments unless they explain a non-obvious *why*. Names carry the meaning.
- Widget methods hang off `*Timui` and take the label/id as the first arg.
- Prefer routing colors through `Theme` rather than hardcoding hex.

## Known issues & gotchas

- **Keyboard input is unimplemented.** The `Backend` interface has no key events;
  tcell only reads Esc/Ctrl-C. Everything is mouse-only — no text field, no focus
  or tab navigation. This is the largest gap (see roadmap "# 2" in [todo.md](todo.md)).

## Roadmap

Tracked loosely in [todo.md](todo.md). Phase 1 items (dropdown, scrollbars, drag
and drop, modal dialog, resizing, transparent colors) are done. Open: text input
field, split-area resize, widget style stacks. Phase 2 is the input story —
tab/focus/input tree, keyboard navigation, shortcuts.

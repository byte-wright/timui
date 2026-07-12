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

## Mental model

**Frame lifecycle.** The app calls widget functions, then [`Timui.Finish`](timui.go)
once per frame. `Finish` runs the deferred overlay queue, diffs the front buffer
against the back buffer, pushes only changed cells to the backend, swaps buffers,
resizes, clears the front buffer, and resets per-frame state. There is no
event-driven redraw; the caller drives a loop (the demo polls at ~30fps).

**Backend seam.** The core has no tcell dependency. Everything the renderer needs
is the 6-method `Backend` interface in [timui.go](timui.go); [tcell/](tcell/) is a
thin adapter. Keep new rendering/input concerns behind this interface.

**Area stack.** Layout is a stack of `mathi.Box2` regions (`PushArea` / `PopArea` /
`CurrentArea`). Widgets draw into `CurrentArea` and advance a cursor via
`moveCursor`. Splits, grids, rows, columns, padding, and dialogs all work by
pushing sub-areas.

**Retained state pattern.** Stateful widgets (mouse, dropdown, draggable,
scrollarea) keep two maps — `last*` and `next*` — keyed by widget ID. Each frame
they read from `last`, write to `next`, and swap in their `finish` method. Follow
this same shape for any new stateful widget.

**IDs.** Widget identity is derived from the label plus the current ID scope stack
([internal/id.go](internal/id.go)). Identical labels in the same scope collide;
use `PushID` / `PopID` ([id.go](id.go)) to disambiguate (e.g. two "+" buttons in
different rows). Changing a widget's label changes its ID and drops its state.

**Deferred overlays.** Dropdowns and dialogs call `runAfter` so they paint on top
after the main pass. Z-order follows registration order in the `after` queue.

**Color / alpha.** Colors are packed `uint32` (`RGBColor` / `RGBAColor` in
[color.go](color.go)). The cell buffer blends alpha ([internal/screen.go](internal/screen.go)),
which is how the dialog dims its backdrop. A cell char of `0` means "leave the
existing glyph, only update color."

## Layout

- [split.go](split.go): `Split().Fixed(...).Factor(...)` computes ranges from a mix
  of fixed sizes and proportional factors. This is the tested core of layout.
- [grid.go](grid.go): callback-based `Grid` / `GridRows` / `GridColumns` that draw
  borders automatically between cells.
- [rows.go](rows.go) / [columns.go](columns.go): a second, imperative
  `.Next()` / `.Finish()` layout API that predates or parallels the grid one.

## Conventions

- No comments unless they explain a non-obvious *why*. Names carry the meaning.
- Widget methods hang off `*Timui` and take the label/id as the first arg.
- Prefer routing colors through `Theme` rather than hardcoding hex.

## Known issues & gotchas

Findings from a read-through; fix opportunistically when touching these files.

- **Keyboard input is unimplemented.** The `Backend` interface has no key events;
  tcell only reads Esc/Ctrl-C. Everything is mouse-only — no text field, no focus
  or tab navigation. This is the largest gap (see roadmap "# 2" in [todo.md](todo.md)).
- **Dialog ignores its title.** [dialog.go](dialog.go) renders the literal string
  `"title"` instead of the `title` argument.
- **Nested deferred overlays silently drop.** `Finish` ranges over the `after`
  slice; Go fixes the length at loop start, so a deferred callback that itself
  calls `runAfter` (e.g. a dropdown inside a dialog) is appended but never run.
- **Scroll area bypasses the theme.** [scrollarea.go](scrollarea.go) hardcodes
  colors and re-parses them with `MustRGBS` every frame; it also has an unresolved
  `// why?` on the knob-visibility check.
- **Layout duplication.** The four `Rows`/`Columns` methods in [grid.go](grid.go)
  are near-identical copy-paste (X↔Y swapped), and the grid API overlaps the
  standalone [rows.go](rows.go)/[columns.go](columns.go). Prime consolidation target.
  Hex parsing in [color.go](color.go) (`RGBS`/`RGBAS`) is likewise duplicated.
- **ID keying allocates.** IDs are built by string concatenation (`parent---id`)
  each widget each frame; a label containing `---` could collide. Consider hashed
  or integer IDs if this ever matters.
- **Concurrency.** The demo's signal handler calls `screen.Fini()` from another
  goroutine while the render loop touches the backend — a data race, harmless in
  practice for the demo.

## Roadmap

Tracked loosely in [todo.md](todo.md). Phase 1 items (dropdown, scrollbars, drag
and drop, modal dialog, resizing, transparent colors) are done. Open: text input
field, split-area resize, widget style stacks. Phase 2 is the input story —
tab/focus/input tree, keyboard navigation, shortcuts.

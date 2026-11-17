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
- Run the demo: `go run ./cmd/example` (mouse and wheel driven; Esc or Ctrl-C to
  quit, PgUp/PgDn/Home/End drive the scroll list)

Stdout/stderr are redirected to `stdout.log` / `stderr.log` in the working dir
while the demo runs (see [util/redirstd.go](util/redirstd.go)), because the TUI
owns the terminal. Check those files when debugging a crash.

Rendering is pinned by snapshot tests ([snapshot_test.go](snapshot_test.go)):
`TestBackend` ([internal/test/backend.go](internal/test/backend.go)) accumulates
cell diffs into a char grid and its string form is checked with `expect`'s
`ToBeSnapshot` against [testdata/](testdata/). On
mismatch the test writes a `.current` file next to the snapshot for diffing;
delete a snapshot file to regenerate it from the current rendering.

## Mental model

**Frame lifecycle.** The app calls widget functions, then [`Timui.Finish`](timui.go)
once per frame. `Finish` runs the deferred overlay queue, diffs the front buffer
against the back buffer, pushes only changed cells to the backend, swaps buffers,
resizes, clears the front buffer, and resets per-frame state. There is no
event-driven redraw; the caller drives a loop (the demo polls at ~30fps).

**Backend seam.** The core has no tcell dependency. Everything the renderer needs
is the `Backend` interface in [timui.go](timui.go); [tcell/](tcell/) is a
thin adapter. Keep new rendering/input concerns behind this interface.

**Input.** Level state (mouse position, buttons) is polled through `Backend`.
Edge state (wheel notches, key presses) is drained once per frame with
`Backend.TakeInput` in `finish` ([input.go](input.go)). Wheel movement is
dispatched like hover, to whatever `MouseInput` area was last registered under the
cursor, and read with `MouseInput.Wheel` on the following frame. Keys have no
focus model: `Timui.Keys` / `KeyPressed` / `RunePressed` hand every widget the
whole frame's presses and each decides what it cares about. Modifiers
(`Timui.Mods`, a `ModMask` of `ModShift`/`ModCtrl`/`ModAlt`) ride along on
`Input` as state and survive `TakeInput`. The backend keeps them current: the
vaxis backend sees a lone modifier press/release through the kitty keyboard
protocol, the tcell backend only when the terminal reports another event.

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

**Text selection & clipboard.** There is no focus model, so — like dragging —
only one selection can exist at a time; it's a single global gesture in
[selection.go](selection.go), not per-widget state. A left-press that lands
where no `MouseInput` claims it arms a drag, which only becomes a selection once
the mouse has moved `selectionDragThreshold` (2) cells on either axis, so a plain
click neither highlights nor copies; releasing over a non-empty selection copies it via `Backend.WriteClipboard` (OSC 52, so it works over SSH
with no clipboard helper installed remotely). Rendering is the library's job,
not the caller's: the highlight is painted directly in `Finish`, after the
`after` queue, so it always sits on top. What a selection can reach is scoped by
a second area-like stack, `CopyRect` ([copyrect.go](copyrect.go)) — pushed
alongside `WithArea`/`WithClip` by Grid/Rows/Columns cells, `ScrollAreaV`,
`ScrollListV`, and `Dialog` — so a drag starting inside one of those can't
select past its edge; unscoped content falls back to the whole screen. A
`MouseInput` that only cares about hover/wheel (a scrollarea's wheel catcher, a
dialog's own click shield) sets its `passive` field so it keeps working
normally but doesn't also block a selection from starting there — only a
genuinely interactive widget should do that.

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

## Scrolling

Two widgets, sharing the bar drawing in `scrollBarV` ([scrollarea.go](scrollarea.go)):

- `ScrollAreaV` draws all of its content into a translated area and clips it. Its
  cost grows with the content, so it suits a page of widgets.
- `ScrollListV` ([scrolllist.go](scrolllist.go)) takes a row count and a callback
  and only calls it for visible rows, so a hundred thousand rows cost a screenful.
  Position lives in a caller owned `ScrollState`, whose `Follow` tails the end of
  the list until the user scrolls away from it.

The bar only writes the position back while its knob is dragged. Knob positions
are whole rows, so feeding them back every frame would snap away finer movement
such as a single wheel notch.

## Conventions

- No comments unless they explain a non-obvious *why*. Names carry the meaning.
- Widget methods hang off `*Timui` and take the label/id as the first arg.
- Prefer routing colors through `Theme` rather than hardcoding hex.

## Known issues & gotchas

- **There is no focus model.** Key presses are delivered to every widget at once,
  so two widgets reading the same key both act on it. Enough for shortcuts and for
  a screen with one scrollable thing, not enough for a text field or tab
  navigation (see roadmap "# 2" in [todo.md](todo.md)).
- **Ctrl-C is the only exit the tcell backend handles.** Escape is delivered as an
  ordinary key, so an app that wants to quit on it checks
  `tui.KeyPressed(timui.KeyEscape)` itself, like the demo does.

## Roadmap

Tracked loosely in [todo.md](todo.md). Phase 1 items (dropdown, scrollbars, drag
and drop, modal dialog, resizing, transparent colors) are done. Open: text input
field, split-area resize, widget style stacks. Phase 2 is the input story —
tab/focus/input tree, keyboard navigation, shortcuts.

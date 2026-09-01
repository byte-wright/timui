TIMUI
=====

**Disclaimer:** This is a personal pet project / experiment

is a terminal ui immediate mode user interface library.
inspired by https://github.com/ocornut/imgui

Focus
-----
- simple
- lightweight

Widgets
-------

are rendered using function calls like:

    if tui.Button("ClickMe +") {
      count++
    }

    if tui.Button("ClickMe -") {
      count--
    }

    tui.Checkbox("Alpha", &checkedA)


Input
-----

Mouse position and buttons are polled, wheel notches and key presses are
collected by the backend and drained once per frame. The wheel goes to the area
under the cursor, keys go to every widget because there is no focus model yet.

    for _, k := range tui.Keys() {
      if k.Key == timui.KeyPageUp {
        state.PageUp()
      }
    }

Backends
--------

The manipulating of the terminal and event handling is done via backends.

Currently there is a github.com/gdamore/tcell backend.

Example
-------

![Demo](cmd/example/demo.gif)
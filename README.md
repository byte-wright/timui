TIMUI
=====

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


Backends
--------

The manipulationg of the terminal and event handling is done via backends.

Currently there is a github.com/gdamore/tcell backend.

Example
-------

![Demo](cmd/example/demo.gif)
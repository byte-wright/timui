package timui

import (
	"os"
	"os/exec"
	"strings"
)

// WriteClipboardLocal tries local clipboard tools (wl-copy, then xclip/xsel) in
// turn. It exists because OSC 52 - the primary mechanism backends use, since it
// also works over SSH - is a silent no-op on VTE-based terminals (GNOME
// Terminal, Guake, Tilix, ...), which don't implement OSC 52 copy at all. This
// only reaches a local X11/Wayland session, so backends call it as a fallback
// alongside OSC 52, not instead of it.
//
// It only runs at all when VTE_VERSION is set (every VTE-based terminal
// exports it), since xclip/xsel have to hold a hidden window open to keep
// owning the X11 CLIPBOARD selection, and creating that window steals focus
// for an instant - harmless on most terminals, but enough to trip a drop-down
// terminal's "hide on lost focus" setting (Guake, and likely others). Terminals
// where OSC 52 already works don't need this fallback, so they don't pay for
// that risk at all.
func WriteClipboardLocal(text string) error {
	if os.Getenv("VTE_VERSION") == "" {
		return nil
	}

	var candidates [][]string

	switch {
	case os.Getenv("WAYLAND_DISPLAY") != "":
		candidates = [][]string{{"wl-copy"}}
	case os.Getenv("DISPLAY") != "":
		candidates = [][]string{
			{"xclip", "-selection", "clipboard"},
			{"xsel", "--clipboard", "--input"},
		}
	default:
		return nil
	}

	var err error

	for _, c := range candidates {
		cmd := exec.Command(c[0], c[1:]...)
		cmd.Stdin = strings.NewReader(text)

		err = cmd.Run()
		if err == nil {
			return nil
		}
	}

	return err
}

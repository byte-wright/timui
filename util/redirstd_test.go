package util_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/byte-wright/timui/util"
)

func TestRedirectStdToFilesInCreatesTheDir(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "state")

	stdout, stderr := os.Stdout, os.Stderr

	restore, err := util.RedirectStdToFilesIn(dir)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("out of the way")
	fmt.Fprintln(os.Stderr, "and this too")

	if os.Stdout == stdout {
		t.Error("stdout was not redirected")
	}

	restore.Restore()

	if os.Stdout != stdout || os.Stderr != stderr {
		t.Error("stdout and stderr were not restored")
	}

	for name, want := range map[string]string{
		"stdout.log": "out of the way",
		"stderr.log": "and this too",
	} {
		data, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			t.Fatal(err)
		}

		if !strings.Contains(string(data), want) {
			t.Errorf("%v holds %q, want it to contain %q", name, string(data), want)
		}
	}
}

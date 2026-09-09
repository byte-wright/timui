package util_test

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
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

func TestRedirectStdToFilesInWithTapsPassesLinesOnAndKeepsTheFiles(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "state")

	var lock sync.Mutex
	var out, err []string

	restore, rerr := util.RedirectStdToFilesIn(dir,
		util.TapStdout(func(l string) {
			lock.Lock()
			defer lock.Unlock()
			out = append(out, l)
		}),
		util.TapStderr(func(l string) {
			lock.Lock()
			defer lock.Unlock()
			err = append(err, l)
		}))
	if rerr != nil {
		t.Fatal(rerr)
	}

	log.Println("through the logger")
	fmt.Println("first\nsecond")
	fmt.Fprint(os.Stderr, "no newline at the end")

	// Restore drains what is still in the pipes, so nothing has to be waited for
	restore.Restore()

	lock.Lock()
	defer lock.Unlock()

	if len(out) != 3 || !strings.HasSuffix(out[0], "through the logger") ||
		out[1] != "first" || out[2] != "second" {
		t.Errorf("stdout tap got %q", out)
	}

	if len(err) != 1 || err[0] != "no newline at the end" {
		t.Errorf("stderr tap got %q", err)
	}

	data, rerr := os.ReadFile(filepath.Join(dir, "stdout.log"))
	if rerr != nil {
		t.Fatal(rerr)
	}

	if !strings.Contains(string(data), "second") {
		t.Errorf("stdout.log holds %q, want the tapped lines in it as well", string(data))
	}
}

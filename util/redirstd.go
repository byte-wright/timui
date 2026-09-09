package util

import (
	"bufio"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"sync"
)

type RestoreSTD struct {
	defaultStdout *os.File
	defaultStderr *os.File
	fileStdout    *os.File
	fileStderr    *os.File

	// pipeStdout and pipeStderr are the write ends given to os.Stdout and
	// os.Stderr while taps are installed, nil otherwise. Closing them ends the
	// reader goroutines.
	pipeStdout *os.File
	pipeStderr *os.File
	drained    sync.WaitGroup

	callbacks []func()
}

// Option configures RedirectStdToFiles / RedirectStdToFilesIn.
type Option func(*options)

type options struct {
	stdoutTap func(line string)
	stderrTap func(line string)
}

// TapStdout adds a live copy of every line written to stdout. An app that owns the
// terminal itself, like a tui, can use it to show the output of the standard logger
// and of libraries printing to stdout instead of only finding it in the file after
// the run.
//
// The tap is called from a reader goroutine, with the line stripped of its newline,
// and must not print to stdout itself, which would feed straight back into it.
// Restore waits until everything written before it has been passed on.
func TapStdout(tap func(line string)) Option {
	return func(o *options) { o.stdoutTap = tap }
}

// TapStderr is TapStdout for stderr.
func TapStderr(tap func(line string)) Option {
	return func(o *options) { o.stderrTap = tap }
}

// RedirectStdToFiles replaces stdout and stderr with file writers in the working
// directory.
// The returned struct has functions to revert replacement of stdout and err with fs files.
func RedirectStdToFiles(opts ...Option) (*RestoreSTD, error) {
	return RedirectStdToFilesIn(".", opts...)
}

// RedirectStdToFilesIn is RedirectStdToFiles writing into dir, which is created if
// it does not exist. An app that already has a place for its scratch files can keep
// the logs there instead of dropping them wherever it was started.
func RedirectStdToFilesIn(dir string, opts ...Option) (*RestoreSTD, error) {
	var o options
	for _, opt := range opts {
		opt(&o)
	}

	err := os.MkdirAll(dir, 0o755)
	if err != nil {
		return nil, err
	}

	restore := &RestoreSTD{
		defaultStdout: os.Stdout,
		defaultStderr: os.Stderr,
	}

	// Open files for stdout and stderr redirection
	restore.fileStdout, err = openLog(filepath.Join(dir, "stdout.log"))
	if err != nil {
		return nil, err
	}

	restore.fileStderr, err = openLog(filepath.Join(dir, "stderr.log"))
	if err != nil {
		restore.fileStdout.Close()
		return nil, err
	}

	outFile, outPipe, err := restore.stream(restore.fileStdout, o.stdoutTap)
	if err != nil {
		restore.fileStdout.Close()
		restore.fileStderr.Close()

		return nil, err
	}

	os.Stdout = outFile
	restore.pipeStdout = outPipe

	errFile, errPipe, err := restore.stream(restore.fileStderr, o.stderrTap)
	if err != nil {
		restore.Restore()
		return nil, err
	}

	os.Stderr = errFile
	restore.pipeStderr = errPipe

	log.Default().SetOutput(os.Stdout)

	// the crash output goes to the file directly, a pipe would need its reader to
	// still be running while the process is already dying
	err = debug.SetCrashOutput(restore.fileStderr, debug.CrashOptions{})
	if err != nil {
		restore.Restore()
		return nil, err
	}

	return restore, nil
}

// stream is the file to write the stream to and, when it is tapped, the write end
// of the pipe feeding the reader that copies into that file. Without a tap the file
// is written to directly and there is no pipe.
func (r *RestoreSTD) stream(file *os.File, tap func(line string)) (*os.File, *os.File, error) {
	if tap == nil {
		return file, nil, nil
	}

	read, write, err := os.Pipe()
	if err != nil {
		return nil, nil, err
	}

	r.drained.Add(1)

	go func() {
		defer r.drained.Done()
		defer read.Close()

		copyLines(read, file, tap)
	}()

	return write, write, nil
}

// copyLines splits r into lines, writes them to file and hands them to tap. A last
// line without a newline is passed on as it is, so output that never terminates its
// line is not swallowed.
func copyLines(r io.Reader, file *os.File, tap func(line string)) {
	buf := bufio.NewReader(r)

	for {
		line, err := buf.ReadString('\n')

		if line != "" {
			if !strings.HasSuffix(line, "\n") {
				line += "\n"
			}

			file.WriteString(line)
			tap(strings.TrimSuffix(strings.TrimSuffix(line, "\n"), "\r"))
		}

		if err != nil {
			return
		}
	}
}

func openLog(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
}

// Restore puts the real stdout and stderr back. The order matters, the pipes are
// closed only once nothing can write to them any more, and the files only once the
// readers have drained what was still in flight.
func (r *RestoreSTD) Restore() {
	os.Stdout = r.defaultStdout
	os.Stderr = r.defaultStderr

	log.Default().SetOutput(os.Stdout)

	if r.pipeStdout != nil {
		r.pipeStdout.Close()
	}

	if r.pipeStderr != nil {
		r.pipeStderr.Close()
	}

	r.drained.Wait()

	r.fileStdout.Close()
	r.fileStderr.Close()

	for _, cb := range r.callbacks {
		cb()
	}
}

func (r *RestoreSTD) Run(f func()) {
	go func() {
		defer func() {
			p := recover()
			if p != nil {
				r.Restore()
				panic(p)
			}
		}()

		f()
	}()
}

func (r *RestoreSTD) AddCallback(f func()) {
	r.callbacks = append(r.callbacks, f)
}

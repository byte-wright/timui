package util

import (
	"log"
	"os"
	"path/filepath"
	"runtime/debug"
)

type RestoreSTD struct {
	defaultStdout *os.File
	defaultStderr *os.File
	fileStdout    *os.File
	fileStderr    *os.File

	callbacks []func()
}

// RedirectStdToFiles replaces stdout and stderr with file writers in the working
// directory.
// The returned struct has functions to revert replacement of stdout and err with fs files.
func RedirectStdToFiles() (*RestoreSTD, error) {
	return RedirectStdToFilesIn(".")
}

// RedirectStdToFilesIn is RedirectStdToFiles writing into dir, which is created if
// it does not exist. An app that already has a place for its scratch files can keep
// the logs there instead of dropping them wherever it was started.
func RedirectStdToFilesIn(dir string) (*RestoreSTD, error) {
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
	os.Stdout = restore.fileStdout

	restore.fileStderr, err = openLog(filepath.Join(dir, "stderr.log"))
	if err != nil {
		restore.fileStdout.Close()
		os.Stdout = restore.defaultStdout

		return nil, err
	}
	os.Stderr = restore.fileStderr

	log.Default().SetOutput(os.Stdout)

	err = debug.SetCrashOutput(restore.fileStderr, debug.CrashOptions{})
	if err != nil {
		return nil, err
	}

	return restore, nil
}

func openLog(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
}

func (r *RestoreSTD) Restore() {
	os.Stdout = r.defaultStdout
	r.fileStdout.Close()

	os.Stderr = r.defaultStderr
	r.fileStderr.Close()

	log.Default().SetOutput(os.Stdout)

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

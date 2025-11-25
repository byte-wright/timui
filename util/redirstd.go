package util

import (
	"log"
	"os"
	"runtime/debug"
)

type RestoreSTD struct {
	defaultStdout *os.File
	defaultStderr *os.File
	fileStdout    *os.File
	fileStderr    *os.File

	callbacks []func()
}

// RedirectStdToFiles replaces stdout and stderr with file writers.
// The returned struct has functions to revert replacement of stdout and err with fs files.
func RedirectStdToFiles() (*RestoreSTD, error) {
	restore := &RestoreSTD{
		defaultStdout: os.Stdout,
		defaultStderr: os.Stderr,
	}

	var err error

	// Open files for stdout and stderr redirection
	restore.fileStdout, err = os.OpenFile("stdout.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, err
	}
	os.Stdout = restore.fileStdout

	restore.fileStderr, err = os.OpenFile("stderr.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
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

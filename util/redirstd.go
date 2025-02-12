package util

import (
	"os"
)

func RedirectStdToFiles() (func(), error) {
	// Open files for stdout and stderr redirection
	stdoutFile, err := os.OpenFile("stdout.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, err
	}
	defaultStdout := os.Stdout
	defer func() {
		stdoutFile.Close()
		os.Stdout = defaultStdout
	}()

	stderrFile, err := os.OpenFile("stderr.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, err
	}
	defaultStderr := os.Stderr
	defer func() {
		stderrFile.Close()
		os.Stderr = defaultStderr
	}()

	origStdout := os.Stdout
	origStderr := os.Stderr

	// Redirect stdout and stderr
	os.Stdout = stdoutFile
	os.Stderr = stderrFile

	return func() {
		os.Stdout = origStdout
		os.Stderr = origStderr
	}, nil
}

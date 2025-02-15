package util

import (
	"os"
)

// RedirectStdToFiles replaces stdout and stderr with file writers.
// The returned func reverts replacement of stdout and err.
func RedirectStdToFiles() (func(), error) {
	// Open files for stdout and stderr redirection
	stdoutFile, err := os.OpenFile("stdout.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, err
	}
	defaultStdout := os.Stdout
	os.Stdout = stdoutFile

	stderrFile, err := os.OpenFile("stderr.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, err
	}
	defaultStderr := os.Stderr
	os.Stderr = stderrFile

	return func() {
		os.Stdout = defaultStdout
		stdoutFile.Close()

		os.Stderr = defaultStderr
		stderrFile.Close()
	}, nil
}

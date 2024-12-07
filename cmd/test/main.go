package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/byte-wright/timui"
	"github.com/byte-wright/timui/tcell"
	"gitlab.com/bytewright/gmath/mathi"
)

func main() {
	// Open files for stdout and stderr redirection
	stdoutFile, err := os.OpenFile("stdout.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		log.Fatalf("Failed to open stdout log file: %v", err)
	}
	defaultStdout := os.Stdout
	defer func() {
		stdoutFile.Close()
		os.Stdout = defaultStdout
	}()

	stderrFile, err := os.OpenFile("stderr.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		log.Fatalf("Failed to open stderr log file: %v", err)
	}
	defaultStderr := os.Stderr
	defer func() {
		stderrFile.Close()
		os.Stderr = defaultStderr
	}()

	// Redirect stdout and stderr
	os.Stdout = stdoutFile
	os.Stderr = stderrFile

	backend, err := tcell.NewBackend()
	if err != nil {
		log.Fatal(err)
	}

	tui := timui.New(backend)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	exit := func() {
		backend.Exit()
		fmt.Println("Exit...")
	}

	defer exit()

	go func() {
		<-signalChan
		exit()
		os.Exit(0)
	}()

	count := 100

	for {
		if backend.Events() {
			break
		}

		time.Sleep(time.Millisecond * 33)

		panel := tui.Panel()
		tui.Label("AAAA")
		tui.Label("BBB")
		tui.Label("CC")

		cols := tui.Columns(timui.ByFactors([]float32{0.25, 0.25, 0.5, 1.5}).Pad(1))

		if tui.Button("ClickMe +") {
			count++
		}
		cols.Next()

		if tui.Button("ClickMe -") {
			count--
		}
		cols.Next()
		if tui.Button("Exit") {
			fmt.Println("no!")
		}
		cols.Next()
		if tui.Button("Quit") {
			fmt.Println("serioulsy!")
		}
		if tui.Button("Settings") {
			fmt.Println("serioulsy!")
		}
		if tui.Button("Start") {
			fmt.Println("serioulsy!")
		}
		cols.Finish()

		tui.Label(fmt.Sprintf("Count %v", count))

		panel.Header()
		tui.Text("[ DATA ]", mathi.Vec2{}, 0xff6666, 0x222222)

		panel.Finish()

		tui.Finish()
	}
}

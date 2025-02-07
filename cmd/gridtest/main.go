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
)

var (
	count    = 0
	tui      *timui.Timui[*tcell.TCellBackend]
	selected int
	checkedA bool
	checkedB bool
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

	tui = timui.New(backend)

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

	for {
		if backend.Events() {
			break
		}

		time.Sleep(time.Millisecond * 33)

		grid := tui.Grid()

		rows := grid.Rows(timui.Split().Fixed(1, 7).Factor(1).Fixed(1))

		header()

		rows.Next()

		{
			cols := rows.Columns(timui.Split().Fixed(12).Factor(0.3, 0.4, 0.5, 0.6, 1))
			countButtons()

			cols.Next()

			{
				revert := tui.Theme.UseBorder(timui.BorderSingle)
				subRows := cols.Rows(timui.Split().Fixed(1).Factor(1).Fixed(1))
				tui.Dropdown("sel1", 10, &selected, func(i int, s bool) {
					tui.Label(fmt.Sprintf("Item %v is my friend", i))
				})

				subRows.Next()

				tui.Label("Select the thing")

				subRows.Next()
				tui.Label(fmt.Sprintf("sel is %v", selected))

				subRows.Finish()
				revert()
			}

			cols.Next()
			checkboxes()

			cols.Next()
			options()
			cols.Next()
			tui.Label("BBBB")

			cols.Finish()
		}

		rows.Next()

		{
			subCols := rows.Columns(timui.Split().Factor(1, 1))
			leftSingleGrid()

			subCols.Next()
			rightSingleGrid()

			subCols.Finish()
		}

		rows.Next()
		footer()

		rows.Finish()

		grid.Finish()

		tui.Finish()
	}
}

func leftSingleGrid() {
	revert := tui.Theme.UseBorder(timui.BorderRoundSingle)
	pad := tui.Pad(0, 1, 0, 1)
	grid := tui.Grid()

	{
		revert2 := tui.Theme.UseBorder(timui.BorderDouble)
		rows := grid.Rows(timui.Split().Factor(1, 1))

		{
			revert3 := tui.Theme.UseBorder(timui.BorderSingle)
			subRows := rows.Rows(timui.Split().Fixed(1).Factor(1))
			tui.Label("Title")

			subRows.Next()
			tui.Label("Content...")
			tui.Label("Content...")
			tui.Label("Content...")

			subRows.Finish()
			revert3()
		}

		rows.Next()
		tui.Label("BOTTOM")

		rows.Finish()
		revert2()
	}

	grid.Finish()
	pad.Finish()
	revert()
}

func rightSingleGrid() {
	revert := tui.Theme.UseBorder(timui.BorderSingle)
	pad := tui.Pad(0, 1, 0, 1)
	grid := tui.Grid()

	{
		revert2 := tui.Theme.UseBorder(timui.BorderDouble)
		rows := grid.Columns(timui.Split().Factor(1, 1))

		{
			revert3 := tui.Theme.UseBorder(timui.BorderSingle)
			subRows := rows.Columns(timui.Split().Fixed(1).Factor(1))
			tui.Label("A")
			tui.Label("B")
			tui.Label("C")

			subRows.Next()
			tui.Label("Content...")
			tui.Label("Content...")
			tui.Label("Content...")

			subRows.Finish()
			revert3()
		}

		rows.Next()

		rows.Finish()
		revert2()
	}

	grid.Finish()
	pad.Finish()
	revert()
}

func countButtons() {
	if tui.Button("ClickMe +") {
		count++
	}

	if tui.Button("ClickMe -") {
		count--
	}

	tui.Label(fmt.Sprintf("Count: %v", count))
	tui.Label(fmt.Sprintf("Sqrd : %v", count*count))
}

func checkboxes() {
	tui.Checkbox("Alpha", &checkedA)
	tui.Checkbox("Beta.1", &checkedB)
	tui.Checkbox("Beta.2", &checkedB)
}

var selectedOption = "a"

func options() {
	og := timui.OptionGroup(tui, "aaa", &selectedOption)

	og.Option("Alpha", "a")
	og.Option("Beta", "b")
	og.Option("Gamma", "c")

	og.Finish()
}

func header() {
	tui.Label("I'm the Header!")
}

func footer() {
	tui.Label("I'm the Foooter!")
}

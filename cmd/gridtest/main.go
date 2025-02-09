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

var (
	count    = 0
	selected int
	checkedA bool
	checkedB bool
	dialogA  bool
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

	for {
		if backend.Events() {
			break
		}

		time.Sleep(time.Millisecond * 33)

		render(tui)
	}
}

func render(tui *timui.Timui) {
	grid := tui.Grid()

	rows := grid.Rows(timui.Split().Fixed(1, 7).Factor(1).Fixed(1))

	header(tui)

	rows.Next()

	{
		cols := rows.Columns(timui.Split().Fixed(12).Factor(0.3, 0.4, 0.5, 0.6, 1))
		countButtons(tui)

		cols.Next()

		tui.Theme.WithBorder(timui.BorderSingle, func() {
			subRows := cols.Rows(timui.Split().Fixed(1).Factor(1).Fixed(1))
			tui.Dropdown("sel1", 10, &selected, func(i int, s bool) {
				tui.Label(fmt.Sprintf("Item %v is my friend", i))
			})

			subRows.Next()

			tui.Label("Select the thing")

			subRows.Next()
			tui.Label(fmt.Sprintf("sel is %v", selected))

			subRows.Finish()
		})

		cols.Next()
		checkboxes(tui)

		cols.Next()
		options(tui)
		cols.Next()

		if tui.Button("Dialog...") {
			dialogA = true
		}

		tui.Dialog("Main DIalog!", &dialogA, func() {
			tui.Label("Fishfingers!")
		})

		cols.Next()
		draggable(tui)

		cols.Finish()
	}

	rows.Next()

	{
		subCols := rows.Columns(timui.Split().Factor(1, 1))
		leftSingleGrid(tui)

		subCols.Next()
		rightSingleGrid(tui)

		subCols.Finish()
	}

	rows.Next()
	footer(tui)

	rows.Finish()

	grid.Finish()

	tui.Finish()
}

var pos = mathi.Vec2{}

func draggable(tui *timui.Timui) {
	tui.Draggable("drag1", mathi.Box2{To: mathi.Vec2{X: 20, Y: 7}}, mathi.Vec2{X: 4, Y: 2}, &pos)

	cursor := tui.CurrentArea().From

	area := mathi.Box2{From: cursor.Add(pos), To: cursor.Add(pos).Add(mathi.Vec2{X: 4, Y: 2})}
	tui.PushArea(area)
	tui.HLine([3]rune{'/', '*', '\\'}, timui.MustRGBS("#ff0"), timui.MustRGBS("#550"))
	tui.PopArea()

	area.From.Y++
	tui.PushArea(area)
	tui.HLine([3]rune{'\\', '*', '/'}, timui.MustRGBS("#ff0"), timui.MustRGBS("#550"))
	tui.PopArea()
}

func leftSingleGrid(tui *timui.Timui) {
	tui.Theme.WithBorder(timui.BorderRoundSingle, func() {
		pad := tui.Pad(0, 1, 0, 1)
		grid := tui.Grid()

		tui.Theme.WithBorder(timui.BorderDouble, func() {
			rows := grid.Rows(timui.Split().Factor(1, 1))

			tui.Theme.WithBorder(timui.BorderSingle, func() {
				subRows := rows.Rows(timui.Split().Fixed(1).Factor(1))
				tui.Label("Title")

				subRows.Next()
				tui.Label("Content...")
				tui.Label("Content...")
				tui.Label("Content...")

				subRows.Finish()
			})

			rows.Next()
			tui.Label("BOTTOM")

			rows.Finish()
		})

		grid.Finish()
		pad.Finish()
	})
}

func rightSingleGrid(tui *timui.Timui) {
	tui.Theme.WithBorder(timui.BorderSingle, func() {
		pad := tui.Pad(0, 1, 0, 1)
		grid := tui.Grid()

		tui.Theme.WithBorder(timui.BorderDouble, func() {
			rows := grid.Columns(timui.Split().Factor(1, 1))

			tui.Theme.WithBorder(timui.BorderSingle, func() {
				subRows := rows.Columns(timui.Split().Fixed(1).Factor(1))
				tui.Label("A")
				tui.Label("B")
				tui.Label("C")

				subRows.Next()
				tui.Label("Content...")
				tui.Label("Content...")
				tui.Label("Content...")

				subRows.Finish()
			})

			rows.Next()

			rows.Finish()
		})

		grid.Finish()
		pad.Finish()
	})
}

func countButtons(tui *timui.Timui) {
	if tui.Button("ClickMe +") {
		count++
	}

	if tui.Button("ClickMe -") {
		count--
	}

	tui.Label(fmt.Sprintf("Count: %v", count))
	tui.Label(fmt.Sprintf("Sqrd : %v", count*count))
}

func checkboxes(tui *timui.Timui) {
	tui.Checkbox("Alpha", &checkedA)
	tui.Checkbox("Beta.1", &checkedB)
	tui.Checkbox("Beta.2", &checkedB)
}

var selectedOption = "a"

func options(tui *timui.Timui) {
	timui.OptionGroup(tui, "aaa", &selectedOption, func(og *timui.OptionGroupElement[string]) {
		og.Option("Alpha", "a")
		og.Option("Beta", "b")
		og.Option("Gamma", "c")
	})
}

func header(tui *timui.Timui) {
	tui.Label("I'm the Header!")
}

func footer(tui *timui.Timui) {
	tui.Label("I'm the Foooter!")
}

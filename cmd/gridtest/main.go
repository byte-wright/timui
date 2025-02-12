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
	"github.com/byte-wright/timui/util"
	"gitlab.com/bytewright/gmath/mathi"
)

var (
	count    = 15
	selected int
	checkedA bool
	checkedB bool
	dialogA  bool
)

func main() {
	revert, err := util.RedirectStdToFiles()
	if err != nil {
		log.Fatal(err)
	}
	defer revert()

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
	tui.Grid(func(grid *timui.Grid) {
		grid.Rows(timui.Split().Fixed(1, 7).Factor(1).Fixed(1), func(rows *timui.GridRows) {
			header(tui)

			rows.Next()

			rows.Columns(timui.Split().Fixed(12).Factor(0.3, 0.4, 0.5, 0.6, 0.5, 0.5), func(cols *timui.GridColumns) {
				countButtons(tui)

				cols.Next()

				tui.Theme.WithBorder(timui.BorderSingle, func() {
					cols.Rows(timui.Split().Fixed(1).Factor(1).Fixed(1), func(subRows *timui.GridRows) {
						tui.Dropdown("sel1", 10, &selected, func(i int, s bool) {
							tui.Label(fmt.Sprintf("Item %v is my friend", i))
						})

						subRows.Next()

						tui.Label("Select the thing")

						subRows.Next()
						tui.Label(fmt.Sprintf("sel is %v", selected))
					})
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

				cols.Next()
				scrollable(tui)
			})

			rows.Next()

			rows.Columns(timui.Split().Factor(1, 1), func(subCols *timui.GridColumns) {
				leftSingleGrid(tui)

				subCols.Next()
				rightSingleGrid(tui)
			})

			rows.Next()
			footer(tui)
		})
	})

	tui.Finish()
}

var pos = mathi.Vec2{}

func scrollable(tui *timui.Timui) {
	tui.ScrollAreaV("scroll", func() {
		for i := 0; i < count; i++ {
			tui.Label(fmt.Sprintf("Line nr %v", i))
		}
	})
}

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
		tui.Pad(0, 1, 0, 1, func() {
			tui.Grid(func(grid *timui.Grid) {
				tui.Theme.WithBorder(timui.BorderDouble, func() {
					grid.Rows(timui.Split().Factor(1, 1), func(rows *timui.GridRows) {
						tui.Theme.WithBorder(timui.BorderSingle, func() {
							rows.Rows(timui.Split().Fixed(1).Factor(1), func(rows *timui.GridRows) {
								tui.Label("Title")

								rows.Next()
								tui.Label("Content...")
								tui.Label("Content...")
								tui.Label("Content...")
							})
						})

						rows.Next()
						tui.Label("BOTTOM")
					})
				})
			})
		})
	})
}

func rightSingleGrid(tui *timui.Timui) {
	tui.Theme.WithBorder(timui.BorderSingle, func() {
		tui.Pad(0, 1, 0, 1, func() {
			tui.Grid(func(grid *timui.Grid) {
				tui.Theme.WithBorder(timui.BorderDouble, func() {
					grid.Columns(timui.Split().Factor(1, 1), func(rows *timui.GridColumns) {
						tui.Theme.WithBorder(timui.BorderSingle, func() {
							rows.Columns(timui.Split().Fixed(1).Factor(1), func(columns *timui.GridColumns) {
								tui.Label("A")
								tui.Label("B")
								tui.Label("C")

								columns.Next()
								tui.Label("Content...")
								tui.Label("Content...")
								tui.Label("Content...")
							})
						})

						rows.Next()
					})
				})
			})
		})
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
	s := tui.Size()
	tui.Label(fmt.Sprintf("I'm the Foooter! %v/%v", s.X, s.Y))
}

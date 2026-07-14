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
	count     = 15
	selected  int
	checkedA  bool
	checkedB1 bool
	checkedB2 bool
	dialogA   bool
)

func main() {
	restore, err := util.RedirectStdToFiles()
	if err != nil {
		log.Fatal(err)
	}
	defer restore.Restore()

	backend, err := tcell.NewBackend()
	if err != nil {
		log.Fatal(err)
	}

	tui := timui.New(backend)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	defer func() {
		backend.Exit()
		fmt.Println("Exit...")
	}()

	for {
		select {
		case <-signalChan:
			return
		default:
		}

		if backend.Events() {
			return
		}

		time.Sleep(time.Millisecond * 33)

		render(tui)
	}
}

func render(tui *timui.Timui) {
	tui.Grid(func(grid *timui.Grid) {
		grid.Rows(timui.Split().Fixed(1, 7).Factor(1).Fixed(1),
			func(*timui.GridCell) {
				header(tui)
			},
			func(cell *timui.GridCell) {
				cell.Columns(timui.Split().Fixed(12, 25, 12, 12, 12).Factor(0.3, 0.3),
					func(*timui.GridCell) {
						countButtons(tui)
					},
					func(cell *timui.GridCell) {
						tui.Theme.WithBorder(timui.BorderSingle, func() {
							cell.Rows(timui.Split().Fixed(1).Factor(1).Fixed(1),
								func(*timui.GridCell) {
									tui.Dropdown("sel1", 10, &selected, func(i int, s bool) {
										tui.Label(fmt.Sprintf("Item %v is my friend", i))
									})
								},
								func(*timui.GridCell) {
									tui.Label("Select the thing")
								},
								func(*timui.GridCell) {
									tui.Label(fmt.Sprintf("sel is %v", selected))
								},
							)
						})
					},
					func(*timui.GridCell) {
						checkboxes(tui)
					},
					func(*timui.GridCell) {
						options(tui)
					},
					func(*timui.GridCell) {
						if tui.Button("Dialog...") {
							dialogA = true
						}

						tui.Dialog("Dialog wow!", &dialogA, func() {
							tui.Label("I like fishfingers!")
						})
					},
					func(*timui.GridCell) {
						draggable(tui)
					},
					func(*timui.GridCell) {
						scrollable(tui)
					},
				)
			},
			func(cell *timui.GridCell) {
				cell.Columns(timui.Split().Factor(1, 1),
					func(*timui.GridCell) {
						leftSingleGrid(tui)
					},
					func(*timui.GridCell) {
						rightSingleGrid(tui)
					},
				)
			},
			func(*timui.GridCell) {
				footer(tui)
			},
		)
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
	tui.WithArea(area, func() {
		tui.HLine([3]rune{'/', '*', '\\'}, timui.MustRGBS("#ff0"), timui.MustRGBS("#550"))
	})

	area.From.Y++
	tui.WithArea(area, func() {
		tui.HLine([3]rune{'\\', '*', '/'}, timui.MustRGBS("#ff0"), timui.MustRGBS("#550"))
	})
}

func leftSingleGrid(tui *timui.Timui) {
	tui.Theme.WithBorder(timui.BorderRoundSingle, func() {
		tui.Pad(0, 1, 0, 1, func() {
			tui.Grid(func(grid *timui.Grid) {
				tui.Theme.WithBorder(timui.BorderDouble, func() {
					grid.Rows(timui.Split().Factor(1, 1),
						func(cell *timui.GridCell) {
							tui.Theme.WithBorder(timui.BorderSingle, func() {
								cell.Rows(timui.Split().Fixed(1).Factor(1),
									func(*timui.GridCell) {
										tui.Label("Title")
									},
									func(*timui.GridCell) {
										tui.Label("Content...")
										tui.Label("Content...")
										tui.Label("Content...")
									},
								)
							})
						},
						func(*timui.GridCell) {
							tui.Label("BOTTOM")
						},
					)
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
					grid.Columns(timui.Split().Factor(1, 1),
						func(cell *timui.GridCell) {
							tui.Theme.WithBorder(timui.BorderSingle, func() {
								cell.Columns(timui.Split().Fixed(1).Factor(1),
									func(*timui.GridCell) {
										tui.Label("A")
										tui.Label("B")
										tui.Label("C")
									},
									func(*timui.GridCell) {
										tui.Label("Content...")
										tui.Label("Content...")
										tui.Label("Content...")
									},
								)
							})
						},
						func(*timui.GridCell) {},
					)
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
	tui.Checkbox("Beta.1", &checkedB1)
	tui.Checkbox("Beta.2", &checkedB2)
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

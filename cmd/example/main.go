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

		tui.Finish()
	}
}

func render(tui *timui.Timui) {
	tui.Grid(func(grid *timui.Grid) {
		grid.Rows(timui.Split().Fixed(1).Factor(1).Fixed(1), func(rows *timui.GridRows) {
			header(tui)

			rows.Next()

			rows.Next()
			footer(tui)
		})
	})
}

func header(t *timui.Timui) {
	t.Label("Header")
}

func footer(t *timui.Timui) {
	t.Label("Footer")
}

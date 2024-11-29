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
	fmt.Println("goo")

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
		os.Exit(0)
	}

	defer exit()

	go func() {
		<-signalChan
		exit()
	}()

	for {
		if backend.Events() {
			break
		}

		time.Sleep(time.Millisecond * 33)

		panel := tui.Panel()
		tui.Label("AAAA")
		tui.Label("BBB")
		tui.Label("CC")

		panel.Header()
		tui.Text("[ DATA ]", mathi.Vec2{})

		panel.Finish()

		tui.Finish()
	}
}

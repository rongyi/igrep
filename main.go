package main

import (
	"github.com/nsf/termbox-go"
)

func main() {
	err := termbox.Init()

	if err != nil {
		panic(err)
	}
	defer termbox.Close()

mainloop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyCtrlC:
				break mainloop
			}
		case termbox.EventError:
			break mainloop
		}
	}
}

package igrep

import (
	"github.com/nsf/termbox-go"
)

const (
	DefaultY     int    = 1
	FilterPrompt string = "[Filter]> "
)

type Engine struct {
	queryCursorIdx int
	query          *Query
	term           *Terminal
	contentOffset  int
}

func NewEngine() (*Engine, error) {
	e := &Engine{
		queryCursorIdx: 0,
		query:          NewQuery([]rune("")),
		term:           NewTerminal(FilterPrompt, DefaultY),
		contentOffset:  0,
	}
	e.queryCursorIdx = e.query.Length()

	return e, nil
}

func (e *Engine) inputChar(ch rune) {
	_ = e.query.Insert([]rune{ch}, e.queryCursorIdx)
	e.queryCursorIdx++
}

func (e *Engine) deleteChar() {
	if i := e.queryCursorIdx - 1; i >= 0 {
		_ = e.query.Delete(i)
		e.queryCursorIdx--
	}
}

func (e *Engine) moveCursorBackward() {
	if i := e.queryCursorIdx - 1; i >= 0 {
		e.queryCursorIdx--
	}
}

func (e *Engine) moveCursorForward() {
	if e.query.Length() > e.queryCursorIdx {
		e.queryCursorIdx++
	}
}
func (e *Engine) moveCursorToTop() {
	e.queryCursorIdx = 0
}
func (e *Engine) moveCursorToEnd() {
	e.queryCursorIdx = e.query.Length()
}

func (e *Engine) Run() int {
	err := termbox.Init()

	if err != nil {
		panic(err)
	}
	defer termbox.Close()

mainloop:
	for {
		ta := &TerminalAttributes{
			Query:           e.query.StringGet(),
			CursorOffset:    e.query.IndexOffset(e.queryCursorIdx),
			Contents:        []string{"hello1", "hello2"},
			ContentsOffsetY: e.contentOffset,
		}
		err = e.term.Draw(ta)
		if err != nil {
			panic(err)
		}

		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case 0:
				e.inputChar(ev.Ch)
			case termbox.KeyBackspace, termbox.KeyBackspace2:
				e.deleteChar()
			case termbox.KeyArrowLeft, termbox.KeyCtrlB:
				e.moveCursorBackward()
			case termbox.KeyArrowRight, termbox.KeyCtrlF:
				e.moveCursorForward()
			case termbox.KeyHome, termbox.KeyCtrlA:
				e.moveCursorToTop()
			case termbox.KeyEnd, termbox.KeyCtrlE:
				e.moveCursorToEnd()
			case termbox.KeyCtrlC:
				break mainloop
			}
		case termbox.EventError:
			break mainloop
		}
	}

	return 0
}

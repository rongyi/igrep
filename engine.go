package igrep

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"io"
	"io/ioutil"
	"regexp"
	"strings"
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
	input          []string
}

func NewEngine(s io.Reader) (*Engine, error) {
	buf, err := ioutil.ReadAll(s)
	if err != nil {
		return nil, err
	}
	flows := strings.Split(string(buf), "\n")
	var fflow []string
	for _, f := range flows {
		match := FlowBeforeAction(f)
		action := FlowAction(f)
		newFlow := FlowDropStats(match) + " " + ActionLearnClean(ActionARPClean(action))
		fflow = append(fflow, newFlow)
	}
	e := &Engine{
		queryCursorIdx: 0,
		query:          NewQuery([]rune("")),
		term:           NewTerminal(FilterPrompt, DefaultY),
		contentOffset:  0,
		input:          fflow,
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

func (e *Engine) scrollToBelow() {
	e.contentOffset++
}

func (e *Engine) scrollToAbove() {
	if o := e.contentOffset - 1; o >= 0 {
		e.contentOffset = o
	}
}

func (e *Engine) scrollToBottom(rownum int) {
	e.contentOffset = rownum - 1
}

func (e *Engine) scrollToTop() {
	e.contentOffset = 0
}

func (e *Engine) getContents() []string {
	filter := e.query.StringGet()
	if filter == "" {
		return e.input
	}
	var ret []string
	for _, s := range e.input {
		match, _ := regexp.MatchString(filter, s)
		if match {
			ret = append(ret, s)
		}
	}
	return ret
}

func (e *Engine) Run() []string {
	err := termbox.Init()

	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	var contents []string
mainloop:
	for {
		bl := len(contents)
		contents = e.getContents()
		if bl != len(contents) {
			e.contentOffset = 0
		}

		ta := &TerminalAttributes{
			Query:           e.query.StringGet(),
			CursorOffset:    e.query.IndexOffset(e.queryCursorIdx),
			Contents:        contents,
			ContentsOffsetY: e.contentOffset,
		}
		err = e.term.Draw(ta)
		if err != nil {
			panic(err)
		}

		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case 0, termbox.KeySpace:
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
			case termbox.KeyCtrlK:
				e.scrollToAbove()
			case termbox.KeyCtrlJ:
				e.scrollToBelow()
			case termbox.KeyCtrlG:
				e.scrollToBottom(len(contents))
			case termbox.KeyCtrlT:
				e.scrollToTop()
			case termbox.KeyCtrlC, termbox.KeyEnter:
				break mainloop
			}
		case termbox.EventError:
			break mainloop
		}
	}

	return contents
}

func (e *Engine) RunWithOutput() int {
	filterOutput := e.Run()
	if len(filterOutput) > 0 {
		fmt.Println(strings.Join(filterOutput, "\n"))
	}

	return 0
}

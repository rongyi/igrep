package igrep

import (
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

type Terminal struct {
	defaultY   int
	prompt     string
	outputArea *[][]termbox.Cell
}

type TerminalAttributes struct {
	Query        string
	Contents     []string
	CursorOffset int
}

func NewTerminal(prompt string, defaultY int) *Terminal {
	t := &Terminal{
		prompt:     prompt,
		defaultY:   defaultY,
		outputArea: &[][]termbox.Cell{},
	}

	return t
}

func (t *Terminal) Draw(attr *TerminalAttributes) error {
	query := attr.Query
	// rows := attr.Contents

	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	// y := t.defaultY

	t.drawFilterLine(query)

	termbox.SetCursor(len(t.prompt)+attr.CursorOffset, 0)

	termbox.Flush()
	return nil
}

func (t *Terminal) drawFilterLine(qs string) error {
	fs := t.prompt + qs

	color := termbox.ColorDefault
	backgroundColor := termbox.ColorDefault

	var cells []termbox.Cell

	var c termbox.Attribute
	for _, s := range fs {
		c = color
		cells = append(cells, termbox.Cell{
			Ch: s,
			Fg: c,
			Bg: backgroundColor,
		})
	}
	t.drawCells(0, 0, cells)

	return nil
}

func (t *Terminal) drawCells(x int, y int, cells []termbox.Cell) {
	i := 0
	for _, c := range cells {
		termbox.SetCell(x+i, y, c.Ch, c.Fg, c.Bg)

		w := runewidth.RuneWidth(c.Ch)
		if w == 0 || w == 2 && runewidth.IsAmbiguousWidth(c.Ch) {
			w = 1
		}

		i += w
	}
}

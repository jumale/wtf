package wtf

import (
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"strings"
)

func NewFormatter(cfg AppConfig) *Formatter {
	return &Formatter{config: cfg}
}

type Formatter struct {
	config AppConfig
}

func (f *Formatter) CenterText(str string, width int) string {
	if width < 0 {
		width = 0
	}

	return fmt.Sprintf("%[1]*s", -width, fmt.Sprintf("%[1]*s", (width+len(str))/2, str))
}

func (f Formatter) DefaultFocusedRowColor() string {
	foreColor := tcell.Color(f.config.Colors.Highlight.Foreground)
	backColor := tcell.Color(f.config.Colors.Highlight.Background)

	return fmt.Sprintf("%s:%s", TcellToTviewColor(foreColor), TcellToTviewColor(backColor))
}

func (f Formatter) DefaultRowColor() string {
	foreColor := tcell.Color(f.config.Colors.Foreground)
	backColor := tcell.Color(f.config.Colors.Background)

	return fmt.Sprintf("%s:%s", TcellToTviewColor(foreColor), TcellToTviewColor(backColor))
}

// PadRow returns a padding for a row to make it the full width of the containing widget.
// Useful for ensurig row highlighting spans the full width (I suspect tcell has a better
// way to do this, but I haven't yet found it)
func (f Formatter) PadRow(offset int, max int) string {
	padSize := max - offset
	if padSize < 0 {
		padSize = 0
	}

	return strings.Repeat(" ", padSize)
}

func (f Formatter) SigilStr(len, pos int, view *tview.TextView) string {
	sigils := ""

	if len > 1 {
		sigils = strings.Repeat(f.config.Paging.PageSigil, pos)
		sigils = sigils + f.config.Paging.SelectedSigil
		sigils = sigils + strings.Repeat(f.config.Paging.PageSigil, len-1-pos)

		sigils = "[lightblue]" + fmt.Sprintf(f.RightAlignFormat(view), sigils) + "[white]"
	}

	return sigils
}

func (f Formatter) RightAlignFormat(view *tview.TextView) string {
	//mutex := &sync.Mutex{}
	//mutex.Lock()
	_, _, w, _ := view.GetInnerRect()
	//mutex.Unlock()

	return fmt.Sprintf("%%%ds", w-1)
}

func (f Formatter) RowColor(idx int) tcell.Color {
	if idx%2 == 0 {
		return tcell.Color(f.config.Colors.Rows.Even)
	}

	return tcell.Color(f.config.Colors.Rows.Odd)
}

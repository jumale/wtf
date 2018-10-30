package clocks

import (
	"fmt"
	"github.com/senorprogrammer/wtf/wtf"
)

func (widget *Widget) display(clocks []Clock) {
	if len(clocks) == 0 {
		widget.TextView.SetText(fmt.Sprintf("\n%s", " no timezone data available"))
		return
	}

	str := ""
	for idx, clock := range clocks {
		str = str + fmt.Sprintf(
			" [%s]%-12s %-10s %7s[white]\n",
			wtf.TcellToTviewColor(widget.formatter.RowColor(idx)),
			clock.Label,
			clock.Time(),
			clock.Date(),
		)
	}

	widget.TextView.SetText(str)
}

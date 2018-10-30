package wtf

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type TextWidget struct {
	*BaseWidget
	TextView *tview.TextView
}

func NewTextWidget(app *tview.Application, name string, config WidgetConfig, focusable bool) *TextWidget {
	widget := &TextWidget{
		BaseWidget: NewBaseWidget(name, config, focusable),
	}

	view := tview.NewTextView()
	view.SetTitle(widget.ContextualTitle(widget.Title))

	view.SetBorderColor(widget.BorderColor())
	view.SetTitleColor(tcell.Color(widget.Config.Colors.Title))
	view.SetTextColor(tcell.Color(widget.Config.Colors.Foreground))
	view.SetBackgroundColor(tcell.Color(widget.Config.Colors.Background))

	view.SetBorder(true)
	view.SetDynamicColors(true)
	view.SetWrap(false)

	// @todo: check if needed
	view.SetChangedFunc(func() {
		app.Draw()
	})

	widget.TextView = view

	return widget
}

/* -------------------- Exported Functions -------------------- */

func (widget *TextWidget) View() View {
	return widget.TextView
}

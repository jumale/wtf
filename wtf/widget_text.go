package wtf

import (
	"github.com/rivo/tview"
)

// TextWidget provides basic widget implementation based on tview.TextView.
// Use it as a base for any kind of widgets which need to display some text.
type TextWidget struct {
	*BaseWidget
	// Need for consumers of this struct to get the TableView-typed view
	TextView *tview.TextView
}

func newTextWidget(title string, app *tview.Application, config WidgetConfig, focusable bool) *TextWidget {
	view := tview.NewTextView()
	view.SetTextColor(config.Colors().Foreground.ToTcell())
	view.SetDynamicColors(true)
	view.SetWrap(false)
	view.SetChangedFunc(func() { // @todo: check if needed
		app.Draw()
	})

	return &TextWidget{
		BaseWidget: newBaseWidget(title, view, app, config, focusable),
		TextView:   view,
	}
}

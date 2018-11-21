package wtf

import (
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type BaseWidget struct {
	title     string
	config    WidgetConfig
	enabled   bool
	focusable bool
	focusKey  string
	view      View
	app       *tview.Application
	Position
}

func newBaseWidget(title string, view View, app *tview.Application, config WidgetConfig, focusable bool) *BaseWidget {
	widget := &BaseWidget{
		title:     title,
		config:    config,
		enabled:   config.Enabled(),
		focusable: focusable,
		app:       app,
		view:      view,
	}

	if config.Title() != "" {
		widget.title = config.Title()
	}
	if config.FocusKey() != "" {
		widget.focusKey = config.FocusKey()
	}

	widget.Position = NewPosition(
		config.Position().Top,
		config.Position().Left,
		config.Position().Width,
		config.Position().Height,
	)

	view.SetTitle(widget.ContextualTitle(widget.title))
	view.SetTitleColor(widget.config.Colors().Title.ToTcell())
	view.SetBackgroundColor(widget.config.Colors().Background.ToTcell())
	view.SetBorder(true)
	view.SetBorderColor(widget.BorderColor())

	return widget
}

/* -------------------- Exported Functions -------------------- */

func (widget *BaseWidget) BorderColor() tcell.Color {
	return focusableItemBorderColor(widget, widget.config.Colors().Border)
}

func (widget *BaseWidget) ContextualTitle(defaultStr string) string {
	if widget.focusKey == "" {
		return fmt.Sprintf(" %s ", defaultStr)
	}

	return fmt.Sprintf(" %s [darkgray::u]%s[::-][green] ", defaultStr, widget.focusKey)
}

func (widget *BaseWidget) Disabled() bool {
	return !widget.Enabled()
}

func (widget *BaseWidget) Enabled() bool {
	return widget.enabled
}

func (widget *BaseWidget) Focus() {
	widget.app.SetFocus(widget.view)
}

func (widget *BaseWidget) Focusable() bool {
	return widget.enabled && widget.focusable
}

func (widget *BaseWidget) FocusKey() string {
	return widget.focusKey
}

func (widget *BaseWidget) SetFocusKey(keyChar string) {
	widget.focusKey = keyChar
}

func (widget *BaseWidget) View() View {
	return widget.view
}

func (widget *BaseWidget) RefreshInterval() int {
	return widget.config.RefreshInterval()
}

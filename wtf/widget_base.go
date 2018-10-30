package wtf

import (
	"fmt"
	"github.com/gdamore/tcell"
)

type BaseWidget struct {
	enabled   bool
	focusable bool
	focusChar string

	Title  string
	Config *WidgetConfig

	Position
}

func NewBaseWidget(title string, config WidgetConfig, focusable bool) *BaseWidget {
	widget := &BaseWidget{
		Title:     title,
		Config:    &config,
		enabled:   config.Enabled,
		focusable: focusable,
	}

	if config.Title != "" {
		widget.Title = config.Title
	}

	widget.Position = NewPosition(
		config.Position.Top,
		config.Position.Left,
		config.Position.Width,
		config.Position.Height,
	)

	return widget
}

/* -------------------- Exported Functions -------------------- */

func (widget *BaseWidget) BorderColor() tcell.Color {
	return focusableItemBorderColor(widget, widget.Config.Colors.Border)
}

func (widget *BaseWidget) ContextualTitle(defaultStr string) string {
	if widget.FocusChar() == "" {
		return fmt.Sprintf(" %s ", defaultStr)
	}

	return fmt.Sprintf(" %s [darkgray::u]%s[::-][green] ", defaultStr, widget.FocusChar())
}

func (widget *BaseWidget) Disable() {
	widget.enabled = false
}

func (widget *BaseWidget) Disabled() bool {
	return !widget.Enabled()
}

func (widget *BaseWidget) Enabled() bool {
	return widget.enabled
}

func (widget *BaseWidget) Focusable() bool {
	return widget.enabled && widget.focusable
}

func (widget *BaseWidget) FocusChar() string {
	return widget.focusChar
}

func (widget *BaseWidget) SetFocusChar(char string) {
	widget.focusChar = char
}

func (widget *BaseWidget) RefreshInterval() int {
	return widget.Config.RefreshInterval
}

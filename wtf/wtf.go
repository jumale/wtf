package wtf

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

// Widget is the main interface of the project.
// It must be implemented by any widget, in order to be added to the applications
type Widget interface {
	// Returns
	Name() string
	Enabler
	Focuser
	Initializer
	Positioner
	Refresher
	Viewer
}

type Viewer interface {
	View() View
}

type View interface {
	tview.Primitive
	Boxable
}

type Boxable interface {
	SetBorderPadding(top, bottom, left, right int) *tview.Box
	SetBackgroundColor(color tcell.Color) *tview.Box
	SetBorderColor(color tcell.Color) *tview.Box
}

type Enabler interface {
	Disabled() bool
	Enabled() bool
}

type UnmarshalFunc func(widgetConfig interface{}) error

type Initializer interface {
	Init(configure UnmarshalFunc, context *AppContext) error
}

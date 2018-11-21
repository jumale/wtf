package wtf

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

// Represents a widget which provides its view implementation to the app.
type Viewer interface {
	View() View
}

// Represents a tview element, which is a mixture
// of the primitive interface and bunch of functions
// implemented in the tview.Box (because the Box is
// already composed into all tview elements which we
// would use for WTF)
type View interface {
	tview.Primitive
	Boxable
}

// Boxable represents interface of tview primitives
// based on tview.Box. This interface covers only a
// limited set of functions, used in this project.
type Boxable interface {
	SetTitle(title string) *tview.Box
	SetTitleColor(color tcell.Color) *tview.Box
	SetBackgroundColor(color tcell.Color) *tview.Box

	SetBorder(show bool) *tview.Box
	SetBorderColor(color tcell.Color) *tview.Box
	SetBorderPadding(top, bottom, left, right int) *tview.Box

	SetInputCapture(capture func(event *tcell.EventKey) *tcell.EventKey) *tview.Box
}

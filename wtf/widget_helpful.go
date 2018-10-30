package wtf

import (
	"github.com/rivo/tview"
)

type HelpfulWidget struct {
	app      *tview.Application
	helpText string
	pages    *tview.Pages
	view     tview.Primitive
}

func NewHelpfulWidget(app *tview.Application, pages *tview.Pages, view tview.Primitive, helpText string) *HelpfulWidget {
	return &HelpfulWidget{
		app:      app,
		pages:    pages,
		view:     view,
		helpText: helpText,
	}
}

func (widget *HelpfulWidget) ShowHelp() {
	closeFunc := func() {
		widget.pages.RemovePage("help")
		widget.app.SetFocus(widget.view)
	}

	modal := NewBillboardModal(widget.helpText, closeFunc)

	widget.pages.AddPage("help", modal, false, true)
	widget.app.SetFocus(modal)
	widget.app.Draw()
}

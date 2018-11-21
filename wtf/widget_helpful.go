package wtf

import (
	"github.com/rivo/tview"
)

type HelpfulWidgetTrait struct {
	app      *tview.Application
	helpText string
	pages    *tview.Pages
	view     tview.Primitive
}

func newHelpfulWidgetTrait(
	app *tview.Application,
	pages *tview.Pages,
	view tview.Primitive,
	helpText string,
) *HelpfulWidgetTrait {

	return &HelpfulWidgetTrait{
		app:      app,
		pages:    pages,
		view:     view,
		helpText: helpText,
	}
}

func (widget *HelpfulWidgetTrait) ShowHelp() {
	closeFunc := func() {
		widget.pages.RemovePage("help")
		widget.app.SetFocus(widget.view)
	}

	modal := NewBillboardModal(widget.helpText, closeFunc)

	widget.pages.AddPage("help", modal, false, true)
	widget.app.SetFocus(modal)
	widget.app.Draw()
}

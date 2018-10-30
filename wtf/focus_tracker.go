package wtf

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type FocusState int

const (
	itemFocused FocusState = iota
	appBoardFocused
	neverFocused
)

type Focuser interface {
	Focusable() bool
	FocusChar() string
	SetFocusChar(string)
}

type Focusable interface {
	Focuser
	Viewer
}

// FocusTracker is used by the app to track which onscreen item currently has focus,
// and to move focus between screen items.
type FocusTracker struct {
	app    *tview.Application
	idx    int
	config *AppConfig
	items  []Focusable
}

func NewFocusTracker(app *tview.Application, cfg *AppConfig, items []Focusable) *FocusTracker {
	return &FocusTracker{
		idx:    -1,
		app:    app,
		config: cfg,
		items:  items,
	}
}

/* -------------------- Exported Functions -------------------- */

// AssignHotKeys assigns an alphabetic keyboard character to each focusable
// item so that the item can be brought into focus by pressing that keyboard key
func (tracker *FocusTracker) AssignHotKeys() {
	if !tracker.config.Navigation.Shortcuts {
		return
	}

	i := 1

	for _, focusable := range tracker.focusables() {
		// Don't have nav characters > "9"
		if i >= 10 {
			break
		}

		focusable.SetFocusChar(string('0' + i))
		i++
	}
}

func (tracker *FocusTracker) FocusOn(char string) bool {
	if !tracker.config.Navigation.Shortcuts {
		return false
	}

	if tracker.focusState() == appBoardFocused {
		return false
	}

	hasFocusable := false

	for idx, focusable := range tracker.focusables() {
		if focusable.FocusChar() == char {
			tracker.blur(tracker.idx)
			tracker.idx = idx
			tracker.focus(tracker.idx)

			hasFocusable = true
			break
		}
	}

	return hasFocusable
}

// Next sets the focus on the next item in the items list. If the current item is
// the last item, sets focus on the first item.
func (tracker *FocusTracker) Next() {
	if tracker.focusState() == appBoardFocused {
		return
	}

	tracker.blur(tracker.idx)
	tracker.increment()
	tracker.focus(tracker.idx)
}

// None removes focus from the currently-focused item.
func (tracker *FocusTracker) None() {
	if tracker.focusState() == appBoardFocused {
		return
	}

	tracker.blur(tracker.idx)
}

// Prev sets the focus on the previous item in the teims list. If the current item is
// the last item, sets focus on the last item.
func (tracker *FocusTracker) Prev() {
	if tracker.focusState() == appBoardFocused {
		return
	}

	tracker.blur(tracker.idx)
	tracker.decrement()
	tracker.focus(tracker.idx)
}

func (tracker *FocusTracker) Refocus() {
	tracker.focus(tracker.idx)
}

/* -------------------- Unexported Functions -------------------- */

func (tracker *FocusTracker) blur(idx int) {
	item := tracker.focusableAt(idx)
	if item == nil {
		return
	}

	view := item.View()
	view.Blur()
	view.SetBorderColor(focusableItemBorderColor(item, tracker.config.Colors.Border))
}

func (tracker *FocusTracker) decrement() {
	tracker.idx = tracker.idx - 1

	if tracker.idx < 0 {
		tracker.idx = len(tracker.focusables()) - 1
	}
}

func (tracker *FocusTracker) focus(idx int) {
	item := tracker.focusableAt(idx)
	if item == nil {
		return
	}

	view := item.View()
	view.SetBorderColor(tcell.Color(tracker.config.Colors.Border.Focused))

	tracker.app.SetFocus(view)
	tracker.app.Draw()
}

func (tracker *FocusTracker) focusables() []Focusable {
	var focusable []Focusable

	for _, item := range tracker.items {
		if item.Focusable() {
			focusable = append(focusable, item)
		}
	}

	return focusable
}

func (tracker *FocusTracker) focusableAt(idx int) Focusable {
	if idx < 0 || idx >= len(tracker.focusables()) {
		return nil
	}

	return tracker.focusables()[idx]
}

func (tracker *FocusTracker) focusState() FocusState {
	if tracker.idx < 0 {
		return neverFocused
	}

	for _, item := range tracker.items {
		if item.View() == tracker.app.GetFocus() {
			return itemFocused
		}
	}

	return appBoardFocused
}

func (tracker *FocusTracker) increment() {
	tracker.idx = tracker.idx + 1

	if tracker.idx == len(tracker.focusables()) {
		tracker.idx = 0
	}
}

func focusableItemBorderColor(item Focuser, cnf ColorsBorderConfig) tcell.Color {
	borderColor := cnf.Normal
	if item.Focusable() {
		borderColor = cnf.Focusable
	}
	return tcell.Color(borderColor)
}

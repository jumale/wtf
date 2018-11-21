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

// Represents a widget item which is able to focus itself, tell if it's
// focusable, and provides getter/setter for a binded hot-key which suppose
// to focuse the widget on press.
type Focuser interface {
	//
	Focus()
	Focusable() bool
	FocusKey() string
	SetFocusKey(string)
}

// Represents an item, which is trackable by the FocusTracker.
type FocusTrackedItem interface {
	Focuser
	Viewer
}

// FocusTracker is used by the appView to track which onscreen item currently has focus,
// and to move focus between screen items.
type FocusTracker struct {
	app    *tview.Application
	idx    int
	config *AppConfig
	items  []FocusTrackedItem
	logger Logger
}

func NewFocusTracker(
	app *tview.Application,
	cfg *AppConfig,
	items []FocusTrackedItem,
	logger Logger,
) *FocusTracker {
	return &FocusTracker{
		idx:    -1,
		app:    app,
		config: cfg,
		items:  items,
		logger: logger,
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
	for _, focusable := range tracker.focusableItems() {
		// Don't have nav characters > "9"
		if i >= 10 {
			break
		}

		tracker.logger.Debugf("FocusTracker: set focus key %s", string('0'+i))
		focusable.SetFocusKey(string('0' + i))
		i++
	}
}

// FocusOn sets the focus on the specified
func (tracker *FocusTracker) FocusOn(keyChar string) bool {
	if !tracker.config.Navigation.Shortcuts {
		return false
	}

	if tracker.focusState() == appBoardFocused {
		return false
	}

	hasFocusable := false

	for idx, focusable := range tracker.focusableItems() {
		if focusable.FocusKey() == keyChar {
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
		tracker.idx = len(tracker.focusableItems()) - 1
	}
}

func (tracker *FocusTracker) focus(idx int) {
	item := tracker.focusableAt(idx)
	if item == nil {
		return
	}

	view := item.View()
	view.SetBorderColor(tracker.config.Colors.Border.Focused.ToTcell())

	item.Focus()
	tracker.app.Draw()
}

func (tracker *FocusTracker) focusableItems() []FocusTrackedItem {
	var focusable []FocusTrackedItem

	for _, item := range tracker.items {
		if item.Focusable() {
			focusable = append(focusable, item)
		}
	}

	return focusable
}

func (tracker *FocusTracker) focusableAt(idx int) FocusTrackedItem {
	if idx < 0 || idx >= len(tracker.focusableItems()) {
		return nil
	}

	return tracker.focusableItems()[idx]
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

	if tracker.idx == len(tracker.focusableItems()) {
		tracker.idx = 0
	}
}

func focusableItemBorderColor(item Focuser, cnf ColorsBorderConfig) tcell.Color {
	borderColor := cnf.Normal
	if item.Focusable() {
		borderColor = cnf.Focusable
	}
	return borderColor.ToTcell()
}

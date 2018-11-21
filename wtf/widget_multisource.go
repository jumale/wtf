package wtf

type MultiSourceWidgetTrait struct {
	module   string
	singular string
	plural   string

	DisplayFunction func()
	Idx             int
	Sources         []string
}

func newMultiSourceSourceWidgetTrait(singular string, plural []string) *MultiSourceWidgetTrait {
	widget := MultiSourceWidgetTrait{}
	widget.Sources = append(widget.Sources, singular)
	widget.Sources = append(widget.Sources, plural...)

	return &widget
}

/* -------------------- Exported Functions -------------------- */

func (widget *MultiSourceWidgetTrait) CurrentSource() string {
	if widget.Idx >= len(widget.Sources) {
		return ""
	}

	return widget.Sources[widget.Idx]
}

func (widget *MultiSourceWidgetTrait) NextSource() {
	widget.Idx = widget.Idx + 1
	if widget.Idx == len(widget.Sources) {
		widget.Idx = 0
	}

	if widget.DisplayFunction != nil {
		widget.DisplayFunction()
	}
}

func (widget *MultiSourceWidgetTrait) PrevSource() {
	widget.Idx = widget.Idx - 1
	if widget.Idx < 0 {
		widget.Idx = len(widget.Sources) - 1
	}

	if widget.DisplayFunction != nil {
		widget.DisplayFunction()
	}
}

func (widget *MultiSourceWidgetTrait) SetDisplayFunction(displayFunc func()) {
	widget.DisplayFunction = displayFunc
}

package wtf

type MultiSourceWidget struct {
	module   string
	singular string
	plural   string

	DisplayFunction func()
	Idx             int
	Sources         []string
}

func NewMultiSourceWidget(singular string, plural []string) *MultiSourceWidget {
	widget := MultiSourceWidget{}
	widget.Sources = append(widget.Sources, singular)
	widget.Sources = append(widget.Sources, plural...)

	return &widget
}

/* -------------------- Exported Functions -------------------- */

func (widget *MultiSourceWidget) CurrentSource() string {
	if widget.Idx >= len(widget.Sources) {
		return ""
	}

	return widget.Sources[widget.Idx]
}

func (widget *MultiSourceWidget) NextSource() {
	widget.Idx = widget.Idx + 1
	if widget.Idx == len(widget.Sources) {
		widget.Idx = 0
	}

	if widget.DisplayFunction != nil {
		widget.DisplayFunction()
	}
}

func (widget *MultiSourceWidget) PrevSource() {
	widget.Idx = widget.Idx - 1
	if widget.Idx < 0 {
		widget.Idx = len(widget.Sources) - 1
	}

	if widget.DisplayFunction != nil {
		widget.DisplayFunction()
	}
}

func (widget *MultiSourceWidget) SetDisplayFunction(displayFunc func()) {
	widget.DisplayFunction = displayFunc
}

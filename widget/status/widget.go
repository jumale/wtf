package status

import (
	"github.com/senorprogrammer/wtf/wtf"
)

type Config struct {
	wtf.BaseWidgetConfig `yaml:",inline"`
}

type Widget struct {
	*wtf.TextWidget
	CurrentIcon int
	config      *Config
}

func New(configure wtf.UnmarshalFunc, app *wtf.AppContext) (wtf.Widget, error) {
	// Initialise
	widget := &Widget{}

	// Define default configs
	widget.config = &Config{}
	// Load configs from config file
	if err := configure(widget.config); err != nil {
		return nil, err
	}

	// Initialise the base widget implementation
	widget.TextWidget = app.TextWidget("Status", widget.config, false)

	// Initialise data and services
	widget.CurrentIcon = 0

	return widget, nil
}

/* -------------------- Exported Functions -------------------- */

func (widget *Widget) Refresh() {
	widget.TextView.SetText(widget.animation())
}

func (widget *Widget) Close() error {
	return nil
}

/* -------------------- Unexported Functions -------------------- */

func (widget *Widget) animation() string {
	icons := []string{"|", "/", "-", "\\", "|"}
	next := icons[widget.CurrentIcon]

	widget.CurrentIcon = widget.CurrentIcon + 1
	if widget.CurrentIcon == len(icons) {
		widget.CurrentIcon = 0
	}

	return next
}

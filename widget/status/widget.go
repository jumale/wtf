package status

import (
	"github.com/senorprogrammer/wtf/wtf"
)

type Config struct {
	wtf.WidgetConfig `yaml:",inline"`
}

type Widget struct {
	*wtf.TextWidget

	CurrentIcon int
	config      *Config
}

func (widget *Widget) Name() string {
	return "status"
}

func (widget *Widget) Init(configure wtf.UnmarshalFunc, context *wtf.AppContext) error {
	context.Logger.Debug("Status: init")

	widget.config = &Config{}
	if err := configure(widget.config); err != nil {
		return err
	}

	widget.TextWidget = wtf.NewTextWidget(context.App, "Status", widget.config.WidgetConfig, false)
	widget.CurrentIcon = 0

	return nil
}

/* -------------------- Exported Functions -------------------- */

func (widget *Widget) Refresh() {
	widget.TextView.SetText(widget.animation())
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

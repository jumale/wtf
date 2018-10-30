package system

import (
	"fmt"
	"time"

	"github.com/senorprogrammer/wtf/wtf"
)

type Config struct {
	wtf.WidgetConfig `yaml:",inline"`
	InputDateFormat  string `yaml:"inputDateFormat"`
	OutputDateFormat string `yaml:"outputDateFormat"`
}

type Widget struct {
	Date    string
	Version string

	*wtf.TextWidget
	systemInfo *SystemInfo
	config     *Config
	logger     wtf.Logger
}

func (widget *Widget) Name() string {
	return "system"
}

func (widget *Widget) Init(configure wtf.UnmarshalFunc, context *wtf.AppContext) error {
	context.Logger.Debug("System: init")

	widget.config = &Config{
		InputDateFormat:  "2006-01-02T15:04:05-0700",
		OutputDateFormat: "Jan _2, 15:04",
	}
	if err := configure(widget.config); err != nil {
		return err
	}

	widget.TextWidget = wtf.NewTextWidget(context.App, "System", widget.config.WidgetConfig, false)

	widget.systemInfo = NewSystemInfo()

	widget.logger = context.Logger

	return nil
}

func (widget *Widget) Refresh() {
	widget.TextView.SetText(
		fmt.Sprintf(
			"%8s: %s\n%8s: %s\n\n%8s: %s\n%8s: %s",
			"Built",
			widget.reformatDate(widget.Date),
			"Vers",
			widget.Version,
			"OS",
			widget.systemInfo.ProductVersion,
			"Build",
			widget.systemInfo.BuildVersion,
		),
	)
}

func (widget *Widget) reformatDate(date string) string {
	str, err := time.Parse(widget.config.InputDateFormat, date)

	if err != nil {
		widget.logger.Error(err.Error())
		return date
	}

	return str.Format(widget.config.OutputDateFormat)
}

package system

import (
	"fmt"
	"time"

	"github.com/senorprogrammer/wtf/wtf"
)

type Config struct {
	wtf.BaseWidgetConfig `yaml:",inline"`
	InputDateFormat      string `yaml:"inputDateFormat"`
	OutputDateFormat     string `yaml:"outputDateFormat"`
}

type Widget struct {
	Date    string
	Version string

	*wtf.TextWidget
	systemInfo *SystemInfo
	config     *Config
	logger     wtf.Logger
}

func CreateConstructor(date string, version string) wtf.WidgetConstructor {
	return func(configure wtf.UnmarshalFunc, app *wtf.AppContext) (wtf.Widget, error) {
		// Initialise
		widget := &Widget{
			Date:    date,
			Version: version,
		}

		// Define default configs
		widget.config = &Config{
			InputDateFormat:  "2006-01-02T15:04:05-0700",
			OutputDateFormat: "Jan _2, 15:04",
		}
		// Load configs from config file
		if err := configure(widget.config); err != nil {
			return nil, err
		}

		// Initialise the base widget implementation
		widget.TextWidget = app.TextWidget("System", widget.config, false)

		// Initialise data and services
		widget.systemInfo = NewSystemInfo()
		widget.logger = app.Logger

		return widget, nil
	}
}

/* -------------------- Exported Functions -------------------- */

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

func (widget *Widget) Close() error {
	return nil
}

/* -------------------- Unexported Functions -------------------- */

func (widget *Widget) reformatDate(date string) string {
	str, err := time.Parse(widget.config.InputDateFormat, date)

	if err != nil {
		widget.logger.Error(err.Error())
		return date
	}

	return str.Format(widget.config.OutputDateFormat)
}

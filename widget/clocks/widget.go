package clocks

import (
	"strings"
	"time"

	"github.com/senorprogrammer/wtf/wtf"
)

type Config struct {
	wtf.BaseWidgetConfig `yaml:",inline"`
	Sort                 SortType          `yaml:"sort"`
	Locations            map[string]string `yaml:"locations"`
}

type Widget struct {
	*wtf.TextWidget
	clockColl ClockCollection
	config    *Config
	formatter *wtf.Formatter
}

func New(configure wtf.UnmarshalFunc, app *wtf.AppContext) (wtf.Widget, error) {
	// Initialise
	widget := &Widget{}

	// Define default configs
	widget.config = &Config{
		Sort: SortAlphabetical,
	}
	// Load configs from config file
	if err := configure(widget.config); err != nil {
		return nil, err
	}

	// Initialise the base widget implementation
	widget.TextWidget = app.TextWidget("World Clocks", widget.config, false)

	// Initialise data and services
	widget.clockColl = widget.buildClockCollection(widget.config.Locations)
	widget.formatter = &app.Formatter

	return widget, nil
}

/* -------------------- Exported Functions -------------------- */

func (widget *Widget) Refresh() {
	widget.display(widget.clockColl.Sorted())
}

func (widget *Widget) Close() error {
	return nil
}

/* -------------------- Unexported Functions -------------------- */

func (widget *Widget) buildClockCollection(locData map[string]string) ClockCollection {
	clockColl := ClockCollection{Sort: widget.config.Sort}

	for label, locStr := range locData {
		timeLoc, err := time.LoadLocation(widget.sanitizeLocation(locStr))
		if err != nil {
			continue
		}

		clockColl.Clocks = append(clockColl.Clocks, NewClock(label, timeLoc))
	}

	return clockColl
}

func (widget *Widget) sanitizeLocation(locStr string) string {
	return strings.Replace(locStr, " ", "_", -1)
}

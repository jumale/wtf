package clocks

import (
	"strings"
	"time"

	"github.com/senorprogrammer/wtf/wtf"
)

type Config struct {
	wtf.WidgetConfig `yaml:",inline"`
	Sort             SortType          `yaml:"sort"`
	Locations        map[string]string `yaml:"locations"`
}

type Widget struct {
	*wtf.TextWidget
	clockColl ClockCollection
	config    *Config
	formatter *wtf.Formatter
}

func (widget *Widget) Name() string {
	return "clocks"
}

func (widget *Widget) Init(configure wtf.UnmarshalFunc, context *wtf.AppContext) error {
	context.Logger.Debug("Clocks: init")

	widget.config = &Config{
		Sort: SortAlphabetical,
	}
	if err := configure(widget.config); err != nil {
		return err
	}

	widget.TextWidget = wtf.NewTextWidget(context.App, "World Clocks", widget.config.WidgetConfig, false)

	widget.clockColl = widget.buildClockCollection(widget.config.Locations)

	widget.formatter = &context.Formatter

	return nil
}

/* -------------------- Exported Functions -------------------- */

func (widget *Widget) Refresh() {
	widget.display(widget.clockColl.Sorted())
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

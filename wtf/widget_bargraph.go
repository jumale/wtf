package wtf

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/rivo/tview"
)

type GraphEntry interface {
	Title() string
	Value() int64
}

type BarGraphConfig struct {
	*BaseWidgetConfig `yaml:",inline"`
	GraphIcon         string `yaml:"graphIcon"`
	GraphStars        uint8  `yaml:"graphStars"`
}

//BarGraph lets make graphs
type BarGraphWidget struct {
	*TextWidget
	Config *BarGraphConfig
}

func newBarGraphWidget(app *tview.Application, name string, config BarGraphConfig, focusable bool) BarGraphWidget {
	graph := BarGraphWidget{}

	if config.GraphIcon == "" {
		config.GraphIcon = "*"
	}
	if config.GraphStars == 0 {
		config.GraphStars = 20
	}

	graph.Config = &config
	graph.TextWidget = newTextWidget(name, app, config, focusable)

	return graph
}

// BuildBars will build a string of * to represent your data of [time][value]
// time should be passed as a int64
func (widget *BarGraphWidget) BuildBars(data []GraphEntry) {
	widget.TextView.SetText(BuildStars(data, widget.Config.GraphStars, widget.Config.GraphIcon))
}

//BuildStars build the string to display
func BuildStars(data []GraphEntry, maxStars uint8, starChar string) string {
	if len(data) == 0 {
		return ""
	}

	var buffer bytes.Buffer

	//store the max value from the array
	maxValue := data[0].Value()

	//store the min value from the array
	minValue := data[0].Value()

	//just getting min and max values
	for i := range data {

		var val = data[i].Value()

		//update max value
		if val > maxValue {
			maxValue = val
		}

		//update minValue
		if val < minValue {
			minValue = val
		}

	}

	// each number = how many stars?
	var starRatio = float64(maxStars) / float64(maxValue-minValue)

	//build the stars
	for i := range data {
		var val = data[i].Value()

		//how many stars for this one?
		var starCount = int(float64(val-minValue) * starRatio)

		if starCount == 0 {
			starCount = 1
		}
		//build the actual string
		var stars = strings.Repeat(starChar, starCount)

		//write the line
		buffer.WriteString(fmt.Sprintf("%s -\t [red]%s[white] - (%d)\n", data[i].Title(), stars, val))
	}

	return buffer.String()
}

/* -------------------- GraphEntry implementations -------------------- */

type TimeGraphEntry [2]int64

func (tge TimeGraphEntry) Title() string {
	t := time.Unix(int64(tge[1]/1000), 0)
	return t.Format("Jan 02, 2006")
}

func (tge TimeGraphEntry) Value() int64 {
	return tge[0]
}

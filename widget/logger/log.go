package logger

import (
	"fmt"
	"io"
	"log"
	"path"
	"time"

	"os"
	"strings"

	"github.com/senorprogrammer/wtf/wtf"
)

type Config struct {
	wtf.BaseWidgetConfig `yaml:",inline"`
	NumLines             int    `yaml:"numLines"`
	DateFormat           string `yaml:"dateFormat"`
}

type Widget struct {
	*wtf.TextWidget
	filePath         string
	config           *Config
	logger           wtf.Logger
	originDateFormat string
}

func New(configure wtf.UnmarshalFunc, app *wtf.AppContext) (wtf.Widget, error) {
	// Initialise
	widget := &Widget{}

	// Define default configs
	widget.config = &Config{
		NumLines:   10,
		DateFormat: "15:04:05",
	}
	// Load configs from config file
	if err := configure(widget.config); err != nil {
		return nil, err
	}

	// Initialise the base widget implementation
	widget.TextWidget = app.TextWidget("Logs", widget.config, false)

	// Initialise data and services
	widget.originDateFormat = app.Config.Log.DateFormat
	widget.filePath = path.Join(app.AppDir, app.Config.Log.File)
	widget.logger = app.Logger

	return widget, nil
}

/* -------------------- Exported Functions -------------------- */

func (widget *Widget) Refresh() {
	widget.TextView.SetTitle(widget.config.Title())

	logLines, err := widget.tailFile()
	if err != nil {
		log.Println(err)
	}
	widget.TextView.SetText(widget.contentFrom(logLines))
}

func (widget *Widget) Close() error {
	return nil
}

/* -------------------- Unexported Functions -------------------- */

func (widget *Widget) contentFrom(logLines []string) string {
	str := ""

	for _, line := range logLines {
		chunks := strings.SplitN(line, "] ", 3)
		date := ""
		level := ""
		msg := ""

		if len(chunks) == 3 {
			date = strings.Trim(chunks[0], "[")
			level = strings.Trim(chunks[1], "[")
			msg = chunks[2]
		} else if len(chunks) == 2 {
			level = strings.Trim(chunks[0], "[")
			msg = chunks[1]
		} else {
			msg = line
		}

		if date != "" {
			t, err := time.Parse(widget.originDateFormat, date)
			if err != nil {
				widget.logger.Errorf("could not parse log date: %s", err)
			} else {
				date = t.Format(widget.config.DateFormat)
			}

			date = fmt.Sprintf("[cadetblue]%s[white] ", date)
		}
		if level != "" {
			level = fmt.Sprintf("[%s]%s[white] ", widget.levelColor(level), level)
		}

		str = str + date + level + msg + "\n"
	}

	return str
}

func (widget *Widget) levelColor(level string) string {
	switch strings.ToLower(level) {
	case "debug":
		return "dodgerblue"
	case "info":
		return "green"
	case "warn", "warning":
		return "orange"
	case "error":
		return "red"
	}
	return "white"
}

func (widget *Widget) tailFile() (lines []string, err error) {
	file, err := os.Open(widget.filePath)
	if err != nil {
		return lines, err
	}
	defer file.Close()

	_, err = file.Seek(0, 2)
	if err != nil {
		return lines, err
	}
	linebreak := []byte("\n")[0]

	for len(lines) < widget.config.NumLines {
		var line []byte
		for {
			_, err = file.Seek(-2, 1)
			if err != nil {
				return lines, err
			}
			buff := make([]byte, 1)
			n, err := file.Read(buff)
			if err != nil && err != io.EOF {
				return lines, err
			}

			if 0 == n || buff[0] == linebreak {
				break
			}
			line = append(buff, line...)
		}
		lines = append([]string{string(line)}, lines...)
	}

	return lines, err
}

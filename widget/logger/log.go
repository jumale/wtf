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
	wtf.WidgetConfig `yaml:",inline"`
	NumLines         int    `yaml:"numLines"`
	DateFormat       string `yaml:"dateFormat"`
}

type Widget struct {
	*wtf.TextWidget
	filePath         string
	config           *Config
	logger           wtf.Logger
	originDateFormat string
}

func (widget *Widget) Name() string {
	return "logger"
}

func (widget *Widget) Init(configure wtf.UnmarshalFunc, context *wtf.AppContext) error {
	widget.config = &Config{
		NumLines:   10,
		DateFormat: "15:04:05",
	}
	if err := configure(widget.config); err != nil {
		return err
	}

	widget.logger = context.Logger
	widget.originDateFormat = context.Config.Log.DateFormat
	widget.TextWidget = wtf.NewTextWidget(context.App, "Logs", widget.config.WidgetConfig, false)
	widget.filePath = path.Join(context.ConfigDir, context.Config.Log.File)

	return nil
}

/* -------------------- Exported Functions -------------------- */

func (widget *Widget) Refresh() {
	widget.TextView.SetTitle(widget.Title)

	logLines, err := widget.tailFile()
	if err != nil {
		log.Println(err)
	}
	widget.TextView.SetText(widget.contentFrom(logLines))
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
	case "error", "fatal":
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

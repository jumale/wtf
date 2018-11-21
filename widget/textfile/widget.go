package textfile

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/gdamore/tcell"
	"github.com/radovskyb/watcher"
	"github.com/rivo/tview"
	"github.com/senorprogrammer/wtf/wtf"
)

const HelpText = `
  Keyboard commands for Textfile:

    /: Show/hide this help window
    h: Previous text file
    l: Next text file
    o: Open the text file in the operating system

    arrow left:  Previous text file
    arrow right: Next text file
`

type Config struct {
	wtf.BaseWidgetConfig `yaml:",inline"`
	FilePath             string   `yaml:"filePath"`
	FilePaths            []string `yaml:"filePaths"`
	Format               bool     `yaml:"format"`
	FormatStyle          string   `yaml:"formatStyle"`
}

type Widget struct {
	*wtf.HelpfulWidgetTrait
	*wtf.MultiSourceWidgetTrait
	*wtf.TextWidget
	config    *Config
	formatter *wtf.Formatter
	fs        *wtf.FileSystem
	logger    wtf.Logger
	watcher   *watcher.Watcher
}

func New(configure wtf.UnmarshalFunc, app *wtf.AppContext) (wtf.Widget, error) {
	// Initialise
	widget := &Widget{}

	// Define default configs
	widget.config = &Config{
		FormatStyle: "vim",
	}
	// Load configs from config file
	if err := configure(widget.config); err != nil {
		return nil, err
	}
	// Don't use a timer for this widget, watch for filesystem changes instead
	widget.config.BaseWidgetConfig.ParamRefreshInterval = 0

	// Initialise the base widget implementation
	widget.TextWidget = app.TextWidget("Textfile", widget.config.BaseWidgetConfig, true)
	widget.HelpfulWidgetTrait = app.HelpfulWidgetTrait(widget.View(), HelpText)
	widget.MultiSourceWidgetTrait = app.MultiSourceWidgetTrait(widget.config.FilePath, widget.config.FilePaths)

	// Initialise data and services
	widget.formatter = &app.Formatter
	widget.fs = &app.FS
	widget.logger = app.Logger
	widget.watcher = watcher.New()

	// Adjust view settings
	widget.TextView.SetWrap(true)
	widget.TextView.SetWordWrap(true)
	widget.TextView.SetInputCapture(widget.keyboardIntercept)

	// Enable file watcher
	widget.SetDisplayFunction(widget.display)
	go widget.watchForFileChanges()

	return widget, nil
}

/* -------------------- Exported Functions -------------------- */

// Refresh is only called once on start-up. Its job is to display the
// text files that first time. After that, the watcher takes over
func (widget *Widget) Refresh() {
	widget.display()
}

func (widget *Widget) Close() error {
	widget.watcher.Close()

	return nil
}

/* -------------------- Unexported Functions -------------------- */

func (widget *Widget) display() {
	title := fmt.Sprintf("[green]%s[white]", widget.CurrentSource())
	title = widget.ContextualTitle(title)

	text := widget.formatter.SigilStr(len(widget.Sources), widget.Idx, widget.TextView) + "\n"

	if widget.config.Format {
		text = text + widget.formattedText()
	} else {
		text = text + widget.plainText()
	}

	//widget.TableView.Lock()
	widget.TextView.SetTitle(title) // <- Writes to TableView's title
	widget.TextView.SetText(text)   // <- Writes to TableView's text
	//widget.TableView.Unlock()
}

func (widget *Widget) fileName() string {
	return filepath.Base(widget.CurrentSource())
}

func (widget *Widget) formattedText() string {
	filePath, _ := widget.fs.ExpandHomeDir(widget.CurrentSource())

	file, err := os.Open(filePath)
	if err != nil {
		return err.Error()
	}

	lexer := lexers.Match(filePath)
	if lexer == nil {
		lexer = lexers.Fallback
	}

	style := styles.Get(widget.config.FormatStyle)
	if style == nil {
		style = styles.Fallback
	}
	formatter := formatters.Get("terminal256")
	if formatter == nil {
		formatter = formatters.Fallback
	}

	contents, _ := ioutil.ReadAll(file)
	iterator, _ := lexer.Tokenise(nil, string(contents))

	var buf bytes.Buffer
	err = formatter.Format(&buf, style, iterator)
	if err != nil {
		return err.Error()
	}

	return tview.TranslateANSI(buf.String())
}

func (widget *Widget) plainText() string {
	filePath, _ := widget.fs.ExpandHomeDir(widget.CurrentSource())

	text, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err.Error()
	}
	return string(text)
}

func (widget *Widget) keyboardIntercept(event *tcell.EventKey) *tcell.EventKey {
	switch string(event.Rune()) {
	case "/":
		widget.ShowHelp()
		return nil
	case "h":
		widget.PrevSource()
		return nil
	case "l":
		widget.NextSource()
		return nil
	case "o":
		widget.fs.OpenFile(widget.CurrentSource())
		return nil
	}

	switch event.Key() {
	case tcell.KeyLeft:
		widget.PrevSource()
		return nil
	case tcell.KeyRight:
		widget.NextSource()
		return nil
	default:
		return event
	}
}

func (widget *Widget) watchForFileChanges() {
	widget.watcher.FilterOps(watcher.Write)

	go func() {
		for {
			select {
			case <-widget.watcher.Event:
				widget.display()
			case err := <-widget.watcher.Error:
				log.Fatalln(err)
			case <-widget.watcher.Closed:
				return
			}
		}
	}()

	// Watch each text-file for changes
	for _, source := range widget.Sources {
		fullPath, err := widget.fs.ExpandHomeDir(source)
		if err == nil {
			if err := widget.watcher.Add(fullPath); err != nil {
				log.Fatalln(err)
			}
		}
	}

	// Start the watching process - it'll check for changes every 100ms.
	if err := widget.watcher.Start(time.Millisecond * 100); err != nil {
		log.Fatalln(err)
	}
}

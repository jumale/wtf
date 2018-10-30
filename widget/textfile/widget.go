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
	wtf.WidgetConfig `yaml:",inline"`
	FilePath         string   `yaml:"filePath"`
	FilePaths        []string `yaml:"filePaths"`
	Format           bool     `yaml:"format"`
	FormatStyle      string   `yaml:"formatStyle"`
}

type Widget struct {
	*wtf.HelpfulWidget
	*wtf.MultiSourceWidget
	*wtf.TextWidget
	config    *Config
	formatter *wtf.Formatter
	fs        *wtf.FileSystem
	logger    wtf.Logger
}

func (widget *Widget) Name() string {
	return "textfile"
}

func (widget *Widget) Init(configure wtf.UnmarshalFunc, context *wtf.AppContext) error {
	context.Logger.Debug("Textfile: init")

	widget.config = &Config{
		FormatStyle: "vim",
	}
	if err := configure(widget.config); err != nil {
		return err
	}
	// Don't use a timer for this widget, watch for filesystem changes instead
	widget.config.WidgetConfig.RefreshInterval = 0

	widget.TextWidget = wtf.NewTextWidget(context.App, "Textfile", widget.config.WidgetConfig, true)
	widget.HelpfulWidget = wtf.NewHelpfulWidget(context.App, context.Pages, widget.TextView, HelpText)
	widget.MultiSourceWidget = wtf.NewMultiSourceWidget(widget.config.FilePath, widget.config.FilePaths)

	widget.formatter = &context.Formatter
	widget.fs = &context.FS
	widget.logger = context.Logger

	widget.SetDisplayFunction(widget.display)
	widget.TextView.SetWrap(true)
	widget.TextView.SetWordWrap(true)
	widget.TextView.SetInputCapture(widget.keyboardIntercept)

	go widget.watchForFileChanges()

	return nil
}

/* -------------------- Exported Functions -------------------- */

// Refresh is only called once on start-up. Its job is to display the
// text files that first time. After that, the watcher takes over
func (widget *Widget) Refresh() {
	widget.display()
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

	//widget.TextView.Lock()
	widget.TextView.SetTitle(title) // <- Writes to TextView's title
	widget.TextView.SetText(text)   // <- Writes to TextView's text
	//widget.TextView.Unlock()
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

	fmt.Println(filePath)

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
	watch := watcher.New()
	watch.FilterOps(watcher.Write)

	go func() {
		for {
			select {
			case <-watch.Event:
				widget.display()
			case err := <-watch.Error:
				log.Fatalln(err)
			case <-watch.Closed:
				return
			}
		}
	}()

	// Watch each textfile for changes
	for _, source := range widget.Sources {
		fullPath, err := widget.fs.ExpandHomeDir(source)
		if err == nil {
			if err := watch.Add(fullPath); err != nil {
				log.Fatalln(err)
			}
		}
	}

	// Start the watching process - it'll check for changes every 100ms.
	if err := watch.Start(time.Millisecond * 100); err != nil {
		log.Fatalln(err)
	}
}

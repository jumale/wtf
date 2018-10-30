package wtf

import (
	"github.com/gdamore/tcell"
	"github.com/pkg/errors"
	"github.com/rivo/tview"
	"math"
	"os"
	"os/exec"
	"regexp"
	"strconv"
)

type Displayable interface {
	Viewer
	Enabler
	Positioner
}

type Display struct {
	Grid   *tview.Grid
	logger Logger
}

func NewDisplay(widgets []Displayable, cnf AppConfig, logger Logger) *Display {
	display := Display{
		Grid:   tview.NewGrid(),
		logger: logger,
	}

	display.Grid.SetBackgroundColor(tcell.Color(cnf.Colors.Background))
	display.Grid.SetColumns(display.columns(cnf)...)
	display.Grid.SetRows(cnf.Grid.Rows...)
	display.Grid.SetBorder(false)

	for _, widget := range widgets {
		display.add(widget)
	}

	return &display
}

/* -------------------- Unexported Functions -------------------- */

func (display *Display) add(widget Displayable) {
	if widget.Disabled() {
		return
	}

	display.Grid.AddItem(
		widget.View(),
		widget.Top(),
		widget.Left(),
		widget.Height(),
		widget.Width(),
		0,
		0,
		false,
	)
}

func (display *Display) columns(cnf AppConfig) []int {
	if cnf.Grid.NumCols == 0 {
		return cnf.Grid.Columns
	}

	size, err := getTerminalSize()
	if err != nil {
		display.logger.Error(err.Error())
		return cnf.Grid.Columns
	}

	var cols []int
	colSize := float64(size.width) / float64(cnf.Grid.NumCols)
	minColSize := int(math.Floor(colSize))
	maxColSize := int(math.Ceil(colSize))
	maxSizedCols := size.width - cnf.Grid.NumCols*minColSize
	minSizedCols := cnf.Grid.NumCols - maxSizedCols

	for i := 0; i <= maxSizedCols; i++ {
		cols = append(cols, maxColSize)
	}
	for i := 0; i <= minSizedCols; i++ {
		cols = append(cols, minColSize)
	}

	return cols
}

type size struct {
	width  int
	height int
}

func getTerminalSize() (s size, err error) {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		return s, err
	}

	r := regexp.MustCompile(`^(\d+)\s(\d+)`)
	match := r.FindSubmatch(out)
	if len(match) == 0 {
		return s, errors.Errorf("could not parse window size from system response '%s'", out)
	}

	s.height, err = strconv.Atoi(string(match[1]))
	if err != nil {
		return s, errors.Errorf("could not parse window height from system response '%s', error: %s", out, err)
	}

	s.width, _ = strconv.Atoi(string(match[2]))
	if err != nil {
		return s, errors.Errorf("could not parse window width from system response '%s', error: %s", out, err)
	}

	return s, nil
}

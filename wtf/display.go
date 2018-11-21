package wtf

import (
	"github.com/pkg/errors"
	"github.com/rivo/tview"
	"math"
	"os"
	"os/exec"
	"regexp"
	"strconv"
)

// Displayable represents an instance which is able to be
// added to the application display
type Displayable interface {
	Viewer
	Enabler
	Positioner
}

// Display represents an application grid area,
// which displays the list of registered Displayable items.
type Display struct {
	Grid   *tview.Grid
	logger Logger
}

// Creates a new Display instance with initialized list of displayable items.
func NewDisplay(widgets []Displayable, cnf AppConfig, logger Logger) *Display {
	display := Display{
		Grid:   tview.NewGrid(),
		logger: logger,
	}

	display.Grid.SetBackgroundColor(cnf.Colors.Background.ToTcell())
	display.Grid.SetColumns(display.columns(cnf)...)
	display.Grid.SetRows(display.rows(cnf)...)
	display.Grid.SetBorder(false)

	for _, widget := range widgets {
		display.add(widget)
	}

	return &display
}

/* -------------------- Unexported Functions -------------------- */

// Adds widget to the current display
func (display *Display) add(widget Displayable) {
	if false == widget.Enabled() {
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

// Returns columns map for grid, based on the provided configs.
// If "numCols" is configured, then the columns map is calculated
// automatically by dividing your terminal width into equal
// desired amount of columns.
// If "numCols" is not set, or if it's not possible to get the
// terminal window size - then the value of "columns" custom map
// is returned.
func (display *Display) columns(cnf AppConfig) []int {
	if cnf.Grid.NumCols == 0 {
		return cnf.Grid.ColumnsMap
	}

	size, err := getTerminalSize()
	if err != nil {
		display.logger.Error(err.Error())
		return cnf.Grid.ColumnsMap
	}

	return splitSize(size.width, cnf.Grid.NumCols)
}

// Returns columns map for grid, based on the provided configs.
// If "numCols" is configured, then the columns map is calculated
// automatically by dividing your terminal width into equal
// desired amount of columns.
// If "numCols" is not set, or if it's not possible to get the
// terminal window size - then the value of "columns" custom map
// is returned.
func (display *Display) rows(cnf AppConfig) []int {
	if cnf.Grid.NumRows == 0 {
		return cnf.Grid.RowsMap
	}

	size, err := getTerminalSize()
	if err != nil {
		display.logger.Error(err.Error())
		return cnf.Grid.RowsMap
	}

	return splitSize(size.height, cnf.Grid.NumRows)
}

type size struct {
	width  int
	height int
}

// Returns window size of the current terminal session,
// or error if it's not possible to get the size
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

// Splits a size (width/height) into array of equal sizes. If the "size" can
// not be divided evenly by the "count", then it decrements some results so
// that the sum of them is equal to the original size.
func splitSize(size int, count int) []int {
	var cols []int
	colSize := float64(size) / float64(count)
	minColSize := int(math.Floor(colSize))
	maxColSize := int(math.Ceil(colSize))
	maxSizedCols := size - count*minColSize
	minSizedCols := count - maxSizedCols

	for i := 0; i <= maxSizedCols; i++ {
		cols = append(cols, maxColSize)
	}
	for i := 0; i <= minSizedCols; i++ {
		cols = append(cols, minColSize)
	}

	return cols
}

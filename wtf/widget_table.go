package wtf

import (
	"fmt"
	"github.com/rivo/tview"
)

// TableWidget provides basic widget implementation based on tview.Table.
// Use it as a base for any kind of widgets which need to display tables.
// The widget provides RenderTable methods which you can use to send table
// data to the app view.
type TableWidget struct {
	*BaseWidget
	// Need for consumers of this struct to get the TableView-typed view
	TableView *tview.Table
}

func newTableWidget(title string, app *tview.Application, config WidgetConfig, focusable bool) *TableWidget {
	view := tview.NewTable()
	return &TableWidget{
		BaseWidget: newBaseWidget(title, view, app, config, focusable),
		TableView:  view,
	}
}

// RenderTable iterates over the table data and sends it to the app view.
func (widget *TableWidget) RenderTable(t Table) {
	for row, r := range t {
		widget.renderRow(row, r)
	}
}

func (widget *TableWidget) renderRow(row int, r Row) {
	for col, cell := range r {
		widget.renderCell(row, col, cell)
	}
}

func (widget *TableWidget) renderCell(row int, col int, c Cell) {
	if c != nil {
		cell := tview.NewTableCell(c.String())
		cell.SetAlign(c.Align())
		cell.SetExpansion(1)
		widget.TableView.SetCell(row, col, cell)
	}
}

// Table is a matrix of abstract Cell interfaces.
type Table []Row
type Row []Cell

// You can implement your custom cell types which know how they should be
// formatted to a cell string value.
type Cell interface {
	// String returns formatted data of the cell
	String() string

	// Raw returns raw data of the cell
	Raw() string

	// Align returns text align rule of the cell.
	// Possible values:
	//     tview.AlignLeft
	//     tview.AlignCenter
	//     tview.AlignRight
	Align() int
}

// Implementation of Cell interface, which makes text bold. Can be used for
// displaying headers.
type Header string

func (h Header) Raw() string {
	return string(h)
}

func (h Header) String() string {
	return fmt.Sprintf("[::b]%s[-:-:-]", string(h))
}

func (h Header) Align() int {
	return tview.AlignCenter
}

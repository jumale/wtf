package wtf

import (
	"github.com/rivo/tview"
)

// AppContext represents a container of service instances and configurations,
// shared from application level to every widget.
// If there is any functionality, which depends on global app configs or functions, and needed to appear in widgets,
// then this functionality should be passed vie this context.
type AppContext struct {
	ConfigDir string
	Config    AppConfig
	Formatter Formatter
	FS        FileSystem
	Logger    Logger
	App       *tview.Application
	Pages     *tview.Pages
}

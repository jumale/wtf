package wtf

import (
	"github.com/rivo/tview"
)

// AppContext represents a container of service-instances and configurations,
// shared from application level to every widget. If there is any functionality,
// which depends on global appView configs or functions, and needed to appear in
// widgets, then this functionality should be passed via this context.
type AppContext struct {
	// Path to the directory with all the application files.
	AppDir string

	// Global appView configs
	Config AppConfig

	// Text-formatter service with all the string formatting functions
	Formatter Formatter

	// File-system service with all the file/dir related functions
	FS FileSystem

	// A service for logging and debugging
	Logger Logger

	appView *tview.Application
	pages   *tview.Pages
}

func (app AppContext) BaseWidget(title string, view View, config WidgetConfig, focusable bool) *BaseWidget {
	return newBaseWidget(title, view, app.appView, config, focusable)
}

func (app AppContext) TextWidget(title string, config WidgetConfig, focusable bool) *TextWidget {
	return newTextWidget(title, app.appView, config, focusable)
}

func (app AppContext) TableWidget(title string, config WidgetConfig, focusable bool) *TableWidget {
	return newTableWidget(title, app.appView, config, focusable)
}

func (app AppContext) HelpfulWidgetTrait(view tview.Primitive, helpText string) *HelpfulWidgetTrait {
	return newHelpfulWidgetTrait(app.appView, app.pages, view, helpText)
}

func (app AppContext) MultiSourceWidgetTrait(singular string, plural []string) *MultiSourceWidgetTrait {
	return newMultiSourceSourceWidgetTrait(singular, plural)
}

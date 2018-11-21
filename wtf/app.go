package wtf

import (
	"github.com/gdamore/tcell"
	"github.com/pkg/errors"
	"github.com/rivo/tview"
	"log"
	"os"
)

// The main WTF application. It is responsible for boot-strapping configuration
// and all services, initialising and configuring all widgets, and finally
// running and displaying the result.
type App struct {
	appDir            string
	configLoader      *ConfigLoader
	watchConfig       bool
	registeredWidgets map[WidgetType]WidgetConstructor
	activeWidgets     *WidgetsList

	appView  *tview.Application
	pageView *tview.Pages

	display      *Display
	focusTracker *FocusTracker
	refresher    *RefreshScheduler

	fs        *FileSystem
	formatter *Formatter
	logger    *FileLogger
}

// NewApp creates a new instance of a WTF application. Returns an error if
// it is failed to initialise the application config. The appDir defines
// a path where application look for a config file, as well as store all
// logs and caches. The configFile defines name of the yaml config file in
// the appDir. If the watchConfig is set true, then the application will
// be dynamically reloaded every time when the config file changes.
func NewApp(appDir string, configFile string, watchConfig bool) (*App, error) {
	// prepare config dir: expand relative path to absolute path
	appDir, err := expandConfigDir(appDir)
	if err != nil {
		return nil, err
	}

	// load configs from yaml file
	configLoader, err := NewConfigLoader(configFile)
	if err != nil {
		return nil, err
	}
	err = configLoader.LoadConfig()
	if err != nil {
		return nil, err
	}

	// initialize logger
	cfg := configLoader.appConfig
	logger, err := NewFileLogger(appDir, *cfg)
	if err != nil {
		return nil, err
	}
	logger.Info("----------------------------------")
	logger.Info("APP: logger is initialized")
	logger.Infof("APP: loaded %d widget configs", len(configLoader.widgetsConfig))

	// initialize application
	app := App{
		configLoader:      configLoader,
		watchConfig:       watchConfig,
		activeWidgets:     &WidgetsList{},
		registeredWidgets: make(map[WidgetType]WidgetConstructor),

		appDir:    appDir,
		formatter: NewFormatter(*cfg),
		fs:        &FileSystem{OpenFileUtil: cfg.OpenFileUtil},
		logger:    logger,
		appView:   tview.NewApplication(),
		pageView:  tview.NewPages(),
	}

	// start watching for config changes, if needed
	if watchConfig {
		go configLoader.WatchChanges(app.update)
	}

	return &app, nil
}

// RegisterWidget adds a widget to a list of registered widgets. The registered
// widgets are used as prototypes for initializing active widgets based on
// configuration. The "widgetType" - is a name of widget which is used to recognize
// this widget in a config file (see the type definition for more info).
// The "constructor" is a function which creates a new instance of the widget
// (see the type definition for more info).
func (app *App) RegisterWidget(widgetType WidgetType, constructor WidgetConstructor) {
	app.registeredWidgets[widgetType] = constructor
}

// Run starts application based on its loaded configuration and the list of
// registered widgets.
func (app *App) Run() error {
	defer app.close()
	app.logger.Infof("APP: registered %d widgets", len(app.registeredWidgets))

	err := os.Setenv("TERM", app.configLoader.appConfig.Term)
	if err != nil {
		return err
	}

	err = app.update()
	if err != nil {
		return err
	}

	return app.appView.
		SetInputCapture(app.keyboardIntercept).
		SetRoot(app.pageView, true).
		Run()
}

// Closes all appView services.
func (app *App) close() {
	err := app.logger.Close()
	if err != nil {
		log.Fatal(err)
	}
}

// Re-initialises all application services and widgets.
func (app *App) update() error {
	err := app.initWidgets()
	if err != nil {
		return err
	}
	app.logger.Infof("APP: configured %d widgets", len(app.activeWidgets.items))
	app.logger.Infof("APP: enabled %d widgets", len(app.activeWidgets.enabled()))

	app.initApp()

	return nil
}

// Re-initialises application widgets
func (app *App) initWidgets() error {
	err := app.clearWidgets()
	if err != nil {
		return err
	}

	for _, widgetConfig := range app.configLoader.widgetsConfig {
		if !widgetConfig.Enabled() {
			continue
		}
		constructor, ok := app.registeredWidgets[widgetConfig.Type()]
		if !ok {
			return errors.Errorf("could not initialize widget configured as type '%s', the widget with such type is not registered in application.")
		}

		err := app.initWidget(constructor, widgetConfig)
		if err != nil {
			return err
		}
	}
	return nil
}

// Sends "close" to all widgets and resets the active widgets list
func (app *App) clearWidgets() (err error) {
	for _, widget := range app.activeWidgets.items {
		err = widget.Close()
	}
	app.activeWidgets = &WidgetsList{}
	return err
}

// Initialises a single widget
func (app *App) initWidget(constructor WidgetConstructor, cfg interimWidgetConfig) error {

	unmarshalFunc := func(target WidgetConfig) error {
		return app.configLoader.unmarshalWidgetConfig(cfg, target, app.logger)
	}

	widget, err := constructor(
		unmarshalFunc,
		&AppContext{
			AppDir:    app.appDir,
			Config:    *app.configLoader.appConfig,
			Formatter: *app.formatter,
			FS:        *app.fs,
			Logger:    app.logger,
			appView:   app.appView,
			pages:     app.pageView,
		},
	)
	if err != nil {
		return err
	}

	app.activeWidgets.add(widget)
	return nil
}

// Initialises application and its services.
func (app *App) initApp() {
	app.refresher = NewRefresher(app.activeWidgets.asRefreshable())
	app.refresher.ScheduleAutoRefresh()

	app.display = NewDisplay(app.activeWidgets.asDisplayable(), *app.configLoader.appConfig, app.logger)
	app.pageView.AddPage("grid", app.display.Grid, true, true)

	bgColor := app.configLoader.appConfig.Colors.Background
	app.display.Grid.SetBackgroundColor(bgColor.ToTcell())
	app.pageView.SetBackgroundColor(bgColor.ToTcell())

	app.focusTracker = NewFocusTracker(
		app.appView,
		app.configLoader.appConfig,
		app.activeWidgets.asFocusable(),
		app.logger,
	)
	app.focusTracker.AssignHotKeys()
}

// Adds global appView-level key bindings.
func (app *App) keyboardIntercept(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyCtrlR:
		app.refresher.Refresh()
	case tcell.KeyTab:
		app.focusTracker.Next()
	case tcell.KeyBacktab:
		app.focusTracker.Prev()
	case tcell.KeyEsc:
		app.focusTracker.None()
	}

	if app.focusTracker.FocusOn(string(event.Rune())) {
		return nil
	}

	return event
}

package app

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"github.com/senorprogrammer/wtf/wtf"
	"log"
	"os"
)

type App struct {
	configDir    string
	configLoader *ConfigLoader
	appView      *tview.Application
	pageView     *tview.Pages
	display      *wtf.Display
	focusTracker *wtf.FocusTracker
	refresher    *wtf.RefreshScheduler
	widgets      *WidgetsList
	fs           *wtf.FileSystem
	formatter    *wtf.Formatter
	logger       *wtf.FileLogger
	watchConfig  bool
}

func NewApp(configDir string, configFile string, watchConfig bool) (*App, error) {
	// prepare config dir: expand relative path to absolute path
	configDir, err := expandConfigDir(configDir)
	if err != nil {
		return nil, err
	}

	// load configs from yaml file
	configLoader, err := NewConfigLoader(configDir, configFile)
	if err != nil {
		return nil, err
	}
	err = configLoader.LoadConfig()
	if err != nil {
		return nil, err
	}

	// initialize logger
	cfg := configLoader.appConfig
	logger, err := wtf.NewFileLogger(configDir, *cfg)
	if err != nil {
		return nil, err
	}
	logger.Info("----------------------------------")
	logger.Info("APP: logger is initialized")
	logger.Infof("APP: loaded %d widget configs", len(configLoader.widgetsConfig))

	// initialize application
	app := App{
		configLoader: configLoader,
		watchConfig:  watchConfig,
		widgets:      &WidgetsList{},

		configDir: configDir,
		formatter: wtf.NewFormatter(*cfg),
		fs:        &wtf.FileSystem{OpenFileUtil: cfg.OpenFileUtil},
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

// AddWidget adds a widget to a list of registered widgets.
// All registered widgets will be initialized on application start.
func (app *App) AddWidget(w wtf.Widget) {
	app.widgets.add(w)
}

func (app *App) Run() error {
	defer app.close()

	app.logger.Infof("APP: registered %d widgets", len(*app.widgets))

	err := os.Setenv("TERM", app.configLoader.appConfig.Term)
	if err != nil {
		return err
	}

	app.initWidgets()
	app.initApp()

	app.focusTracker = wtf.NewFocusTracker(app.appView, app.configLoader.appConfig, app.widgets.asFocusable())
	app.focusTracker.AssignHotKeys()

	return app.appView.
		SetInputCapture(app.keyboardIntercept).
		SetRoot(app.pageView, true).
		Run()
}

func (app *App) close() {
	err := app.logger.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func (app *App) update() {
	app.initWidgets()
	app.initApp()
}

func (app *App) initApp() {
	app.refresher = wtf.NewRefresher(app.widgets.asRefreshable())
	app.refresher.ScheduleAutoRefresh()

	app.display = wtf.NewDisplay(app.widgets.asDisplayable(), *app.configLoader.appConfig, app.logger)
	app.pageView.AddPage("grid", app.display.Grid, true, true)

	bgColor := app.configLoader.appConfig.Colors.Background
	app.display.Grid.SetBackgroundColor(tcell.Color(bgColor))
	app.pageView.SetBackgroundColor(tcell.Color(bgColor))
}

func (app *App) initWidgets() {
	for _, widget := range *app.widgets {
		app.logger.Debugf("APP: init widget '%s'", widget.Name())
		err := widget.Init(
			func(cnf interface{}) error {
				return app.configLoader.unmarshalWidgetConfig(widget.Name(), cnf, app.logger)
			},
			&wtf.AppContext{
				ConfigDir: app.configDir,
				Config:    *app.configLoader.appConfig,
				Formatter: *app.formatter,
				FS:        *app.fs,
				Logger:    app.logger,
				App:       app.appView,
				Pages:     app.pageView,
			},
		)
		checkWidgetInitError(err, widget.Name())
	}
}

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

func checkWidgetInitError(err error, name string) {
	if err != nil {
		log.Panicf("failed while initilazing widget '%s': %s", name, err)
	}
}

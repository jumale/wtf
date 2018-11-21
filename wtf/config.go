package wtf

import (
	"github.com/gdamore/tcell"
	"os"
)

// Provides default configurations recommended for an application. It can be
// used to fill undefined values in a loaded configuration.
func GetDefaultAppConfig() AppConfig {
	color := Color("")
	return AppConfig{
		Term:            os.Getenv("TERM"),
		RefreshInterval: 300, // 5 minutes
		OpenFileUtil:    "open",
		Log: LogConfig{
			File:       "log.txt",
			DateFormat: "2006-01-02 15:04:05",
			Level:      LogLevelInfo,
		},
		Colors: ColorsConfig{
			Title: color.FromTcell(tcell.ColorDefault),

			ColorsTextConfig: ColorsTextConfig{
				Foreground: color.FromTcell(tcell.ColorDefault),
				Background: color.FromTcell(tcell.ColorDefault),
			},
			Highlight: ColorsTextConfig{
				Foreground: color.FromTcell(tcell.ColorBlack),
				Background: color.FromTcell(tcell.ColorOrange),
			},
			Border: ColorsBorderConfig{
				Focusable: color.FromTcell(tcell.ColorRed),
				Focused:   color.FromTcell(tcell.ColorGray),
				Normal:    color.FromTcell(tcell.ColorGray),
			},
			Rows: ColorsRowsConfig{
				Even: color.FromTcell(tcell.ColorLightBlue),
				Odd:  color.FromTcell(tcell.ColorWhite),
			},
		},
		Navigation: NavConfig{
			Shortcuts: true,
		},
		Paging: PagingConfig{
			PageSigil:     "*",
			SelectedSigil: "_",
		},
	}
}

// ------------------------------ ROOT CONFIGS ---------------------------- //

type Config struct {
	App           AppConfig             `yaml:"appView"`
	WidgetsConfig []interimWidgetConfig `yaml:"widgets"`
}

// ---------------------------- GENERIC CONFIGS --------------------------- //

// Defines all color configurations, used in appView and its widgets.
type ColorsConfig struct {
	// The composed text color config represents bg and fg colors for a
	// normal text.
	ColorsTextConfig `yaml:",inline"`

	// Color of the widget title (displayed in the header).
	Title Color `yaml:"title"`

	// Text color config for focused (highlighted) text.
	Highlight ColorsTextConfig `yaml:"highlight"`

	// Color configs for widget borders.
	Border ColorsBorderConfig `yaml:"border"`

	// Color configs for row-like data.
	Rows ColorsRowsConfig `yaml:"rows"`
}

// A common text color configuration, which specifies background and foreground
// color of text.
type ColorsTextConfig struct {
	Background Color `yaml:"background"`
	Foreground Color `yaml:"foreground"`
}

// Defines colors for widget borders.
type ColorsBorderConfig struct {
	// Color for regular borders.
	Normal Color `yaml:"normal"`

	// Color for borders of a widget which can be focused.
	Focusable Color `yaml:"focusable"`

	// Color for borders of a widget which is focused.
	Focused Color `yaml:"focused"`
}

// Configures colors for data which can be represented as rows (e.g. lists,
// tables). Note: it's decided in every widget, whether to use these configs
// or not. If a widget displays rows-like data, but does not use these values
// to configure its rows colors, then obviously, changing these values will
// not help you to manipulate those rows colors.
type ColorsRowsConfig struct {
	// Background colors for even rows.
	Even Color `yaml:"even"`

	// Background colors for odd rows.
	Odd Color `yaml:"odd"`
}

// Defines pagination parameters for widgets which support switching between
// multiple data sources. E.g. switching between repositories in Github widget.
type PagingConfig struct {
	// The character which will be displayed for every single data source (page).
	PageSigil string `yaml:"pageSigil"`

	// The character which will be displayed for the current (selected) data
	// source (page).
	SelectedSigil string `yaml:"selectedSigil"`
}

// ------------------------------ APP CONFIG ------------------------------ //

// AppConfig is a root configuration, which contains two kinds of config
// parameters:
//     - any parameters required by application itself and its services
//
//     - global parameters for widgets, which are merged with widgets' own
//       configs and provide kind of default state
//
// For example Colors: if there is no specified config parameters for colors
// for a specific widget, then the parameters from the appView level will be taken
// as defaults. This functionality is provided by the GlobalsConfigMerger,
// interface, which is implemented by the BaseWidgetConfig config.
type AppConfig struct {
	// Default refresh interval for all widgets. Set it to -1 to disable
	// refreshing at all.
	RefreshInterval int `yaml:"refreshInterval"`

	// The command on host machine which is responsible for opening files.
	// For example "open" in OSX.
	OpenFileUtil string `yaml:"openFileUtil"`

	// Value for the "TERM" env variable in OS. By default should be equal
	// to the current "TERM" env value. If configured another value, then
	// that value should be set to the "TERM" env variable by an application.
	Term string `yaml:"term"`

	// Configures a grid layout for appView.
	Grid GridConfig `yaml:"grid"`

	// Configures navigation options for appView.
	Navigation NavConfig `yaml:"navigation"`

	// Configures logger.
	Log LogConfig `yaml:"log"`

	// Global color configs. It is used on different levels of the WTF application,
	// to defined all background-, border- and text-colors. As well it is used
	// as default set of values in every widget configuration.
	Colors ColorsConfig `yaml:"colors"`

	// Defines global pagination configs, which is used as default values in
	// every widget configuration.
	Paging PagingConfig `yaml:"paging"`
}

// Defines a grid layout, which gonna contain all widgets
type GridConfig struct {
	// By specifying numCols you telling application to divide your terminal
	// into the X equal columns. This value is based on a current width value
	// of your terminal window. Note that it might be not possible to get the
	// current window size - in that case this values is skipped.
	NumCols int `yaml:"numCols"`

	// Configures the number of grid rows. Works in the same way as the
	// "numCols" parameter.
	NumRows int `yaml:"numRows"`

	// Configures a custom columns map. By specifying an array of integers,
	// you describe how many columns you need, and what with (in characters)
	// each column should have. This value is taken if "numColumns" is not
	// set or if it's not possible to automatically calculate widths of columns.
	ColumnsMap []int `yaml:"columns"`

	// Configures a custom map of rows. Works in the same way as the "columns"
	// parameter.
	RowsMap []int `yaml:"rows"`
}

// Defines navigation-related configurations.
type NavConfig struct {
	// Specifies whether shortcuts (hot-keys) are enables for the whole
	// application, or not.
	Shortcuts bool `yaml:"shortcuts"`
}

// Defines configs for internal logger
type LogConfig struct {
	// Where to store logs. It's a path, relative to some application folder.
	File string `yaml:"file"`

	// Specifies a date format for every log-entry stored in the log file.
	DateFormat string `yaml:"dateFormat"`

	// Specifies level of logging.
	// Possible values:
	//     - debug
	//     - info
	//     - warn (or warning)
	//     - error
	//
	// Logs which are higher than the specified level will not be written to
	// the log file.
	Level LogLevel `yaml:"level"`
}

// ------------------------- BASE WIDGET CONFIG --------------------------- //

// Represents a widget position in appView grid.
type WidgetPositionConfig struct {
	// Number of grid-row for widget's top-left corner (starting from 0)
	Top int `yaml:"top"`

	// Number of grid-column for widget's top-left corner (starting from 0)
	Left int `yaml:"left"`

	// How many rows the widget takes.
	Height int `yaml:"height"`

	// How many columns the widget takes
	Width int `yaml:"width"`
}

// A base widget config, which suppose to be implemented by any widget,
// registered in an application. You can compose this struct into your custom
// widget's config struct, to have a basic set of config parameters. Do not
// forget to specify `yaml:",inline"` for composed fields.
type BaseWidgetConfig struct {
	// The type tells us what kind of widget we are going to load and configure
	// (e.g. "textfile", "github", etc). The list of available widget-types
	// is defined somewhere in application high-level code, where all widgets
	// are registered into appView.
	ParamType WidgetType `yaml:"type"`

	// Specifies if the widget enabled or not. A disabled widget is initialized,
	// but never refreshed (displayed).
	ParamEnabled bool `yaml:"enabled"`

	// Widget title, which is displayed in the widget's header.
	ParamTitle string `yaml:"title"`

	// Custom refresh interval for the particular widget. If not specified
	// or set to 0, then the refreshInterval from application config is taken
	// as a default value. Set it to -1 to disable refreshing at all.
	ParamRefreshInterval int `yaml:"refreshInterval"`

	// Binds a key, responsible for focusing this widget. Expected a string
	// representation of one of the keyboard keys, e.g. "h", "=", "4", etc.
	//
	// If this parameter is not set, then the focus-key is equal to an index
	// of the fidget in a list of focusable widgets (index starts from 1).
	//
	// For example there is the next list of widgets:
	//     - "foo", nonFocusable
	//     - "bar", focusable, index 1, focus-key "1"
	//     - "baz", nonFocusable
	//     - "nod", focusable, index 2, focus-key "2"
	//     - "sed", focusable, index 3, focus-key "3"
	//
	// So that, to focus the "nod" widget you would need to press num-key "2"
	ParamFocusKey string `yaml:"focusKey"`

	// Configures the widget position in an appView grid.
	ParamPosition WidgetPositionConfig `yaml:"position"`

	// Custom color configs for the widget. All missing values will be filled
	// with global appView color configs.
	ParamColors ColorsConfig `yaml:"colors"`

	// Custom pagination configs for the widget. All missing values will be
	// filled with global appView pagination configs.
	ParamPaging PagingConfig `yaml:"paging"`
}

// ----------------- IMPLEMENTING WIDGET CONFIG INTERFACE ----------------- //

func (w BaseWidgetConfig) Type() WidgetType {
	return w.ParamType
}

func (w BaseWidgetConfig) Title() string {
	return w.ParamTitle
}

func (w BaseWidgetConfig) Enabled() bool {
	return w.ParamEnabled
}

func (w BaseWidgetConfig) FocusKey() string {
	return w.ParamFocusKey
}

func (w BaseWidgetConfig) RefreshInterval() int {
	return w.ParamRefreshInterval
}

func (w BaseWidgetConfig) Colors() ColorsConfig {
	return w.ParamColors
}

func (w BaseWidgetConfig) Position() WidgetPositionConfig {
	return w.ParamPosition
}

func (w BaseWidgetConfig) Paging() PagingConfig {
	return w.ParamPaging
}

// ------------------ IMPLEMENTING DEFAULT CONFIG SETTER ------------------ //

// The defaultConfigSetter interface subscribes a widget config for receiving
// a default configuration before the actual values from a config file are
// unmarshaled into the object.
type defaultConfigSetter interface {
	// Sets default values for a widget config, considering the global pre-loaded
	// application config.
	setDefaults(cfg AppConfig)
}

func (w *BaseWidgetConfig) setDefaults(cfg AppConfig) {
	w.ParamEnabled = true
	w.ParamRefreshInterval = cfg.RefreshInterval
	w.ParamColors = cfg.Colors
	w.ParamPaging = cfg.Paging
	w.ParamPosition = WidgetPositionConfig{
		Width:  1,
		Height: 1,
		Top:    0,
		Left:   0,
	}
}

package wtf

import (
	"github.com/gdamore/tcell"
	"os"
)

func GetDefaultAppConfig() AppConfig {
	return AppConfig{
		Term:            os.Getenv("TERM"),
		RefreshInterval: -1,
		OpenFileUtil:    "open",
		Log: LogConfig{
			File:       "log.txt",
			DateFormat: "2006-01-02 15:04:05",
			Level:      LogLevelInfo,
		},
		Colors: ColorsConfig{
			Title: Color(tcell.ColorDefault),

			ColorsTextConfig: ColorsTextConfig{
				Foreground: Color(tcell.ColorDefault),
				Background: Color(tcell.ColorDefault),
			},
			Highlight: ColorsTextConfig{
				Foreground: Color(tcell.ColorBlack),
				Background: Color(tcell.ColorOrange),
			},
			Border: ColorsBorderConfig{
				Focusable: Color(tcell.ColorRed),
				Focused:   Color(tcell.ColorGray),
				Normal:    Color(tcell.ColorGray),
			},
			Rows: ColorsRowsConfig{
				Even: Color(tcell.ColorLightBlue),
				Odd:  Color(tcell.ColorWhite),
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

// ------------------------------------- GENERIC CONFIGS ------------------------------------- //

type ColorsConfig struct {
	ColorsTextConfig `yaml:",inline"`
	Title            Color              `yaml:"title"`
	Highlight        ColorsTextConfig   `yaml:"highlight"`
	Border           ColorsBorderConfig `yaml:"border"`
	Rows             ColorsRowsConfig   `yaml:"rows"`
}

type ColorsTextConfig struct {
	Background Color `yaml:"background"`
	Foreground Color `yaml:"foreground"`
}

type ColorsBorderConfig struct {
	Focusable Color `yaml:"focusable"`
	Focused   Color `yaml:"focused"`
	Normal    Color `yaml:"normal"`
}

type ColorsRowsConfig struct {
	Even Color `yaml:"even"`
	Odd  Color `yaml:"odd"`
}

type PagingConfig struct {
	PageSigil     string `yaml:"pageSigil"`
	SelectedSigil string `yaml:"selectedSigil"`
}

// --------------------------------------- APP CONFIG --------------------------------------- //

type AppConfig struct {
	Grid            GridConfig `yaml:"grid"`
	Navigation      NavConfig  `yaml:"navigation"`
	RefreshInterval int        `yaml:"refreshInterval"`
	OpenFileUtil    string     `yaml:"openFileUtil"`
	Term            string     `yaml:"term"`
	Log             LogConfig  `yaml:"log"`

	Colors ColorsConfig `yaml:"colors"`
	Paging PagingConfig `yaml:"paging"`
}

type GridConfig struct {
	NumCols int   `yaml:"numCols"`
	Columns []int `yaml:"columns"`
	Rows    []int `yaml:"rows"`
}

type NavConfig struct {
	Shortcuts bool `yaml:"shortcuts"`
}

type LogConfig struct {
	File       string   `yaml:"file"`
	DateFormat string   `yaml:"dateFormat"`
	Level      LogLevel `yaml:"level"`
}

// ------------------------------------- WIDGET CONFIG -------------------------------------- //

type DefaultConfigSetter interface {
	SetDefaults(cfg AppConfig)
}

type WidgetConfig struct {
	// widget own configs
	Title           string               `yaml:"title"`
	Enabled         bool                 `yaml:"enabled"`
	Position        WidgetPositionConfig `yaml:"position"`
	RefreshInterval int                  `yaml:"refreshInterval"`

	// configs, merged with WTF global configs
	Colors ColorsConfig `yaml:"colors"`
	Paging PagingConfig `yaml:"paging"`
}

func (w *WidgetConfig) SetDefaults(cfg AppConfig) {
	w.Enabled = true
	w.RefreshInterval = cfg.RefreshInterval
	w.Colors = cfg.Colors
	w.Paging = cfg.Paging
	w.Position = WidgetPositionConfig{
		Width:  1,
		Height: 1,
		Top:    0,
		Left:   0,
	}
}

type WidgetPositionConfig struct {
	Top    int `yaml:"top"`
	Left   int `yaml:"left"`
	Height int `yaml:"height"`
	Width  int `yaml:"width"`
}

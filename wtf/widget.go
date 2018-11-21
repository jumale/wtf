package wtf

// Widget is the main interface of the project. It must be implemented by
// any widget, in order to be added to the application.
type Widget interface {
	Enabler
	Focuser
	Positioner
	Refresher
	Viewer
	Close() error
}

// WidgetConstructor represents a function which knows how to create a new
// instance of a widget. It must be implemented by any widget, in order to
// be added to the application. The widget constructor should create a new
// instance of a widget, configure it and initialise all its services.
// The "configure" is an unmarshal function, which parses yaml configs into
// a config object. The "app" is an application container which you can use
// to access global functions, services and configurations.
type WidgetConstructor func(configure UnmarshalFunc, app *AppContext) (Widget, error)

// UnmarshalFunc is a function which abstracts from a widget implementation
// the way how configs are encoded from a file to a config object.
type UnmarshalFunc func(cnf WidgetConfig) error

// Any widget config should provide the basic config params, needed on different
// levels of the WTF application. Also it should provide an ability to be
// merged with global application config.
type WidgetConfig interface {
	Type() WidgetType
	Title() string
	Enabled() bool
	FocusKey() string
	RefreshInterval() int
	Colors() ColorsConfig
	Position() WidgetPositionConfig
	Paging() PagingConfig
}

// Represents a type-name of a widget. It is used to correctly point a widget
// configuration in config.yml file to a widget implementation.
// Example:
//     considering there is a config with a single widget:
//         wtf: {widgets: [{type: "github", "enabled": true}]}
//     and we want to tell the appView, that this widget config should be received
//     by github.Widget type, then we would use var t WidgetType = "github"
//     to register github.Widget in the appView.
type WidgetType string

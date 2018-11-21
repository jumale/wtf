package hello_world

import (
	"fmt"
	"github.com/senorprogrammer/wtf/wtf"
)

// -------------------------------- CONFIG -------------------------------- //
// First you need to define a config struct for your new widget. It must
// implement at least the `wtf.WidgetConfig` interface. You can also add any
// extra parameters you need in your widget.
//
// There is already a `wtf.BaseWidgetConfig` struct which implements the
// `wtf.WidgetConfig` interface, so if your widget is not expecting any extra
// parameters then you can just use this struct as your widget config. If
// your widget needs some extra config parameters, then you can create a new
// struct and compose `wtf.BaseWidgetConfig` into it:
type Config struct {
	// Compose the BaseWidgetConfig to implement the required interface
	*wtf.BaseWidgetConfig `yaml:",inline"`

	// Specify your custom config parameters. For example:

	// Configure a custom color:
	ContentColor wtf.Color `yaml:"contentColor"`
	// Note: to configure colors, use the `wtf.Color` type - it is based on
	// string type and represents a color label (e.g. "reg", "green", etc).
	// This type can convert itself to a `tcell.Color` constant, which you
	// will need if you want to manipulate the view colors.

	// Configure some boolean parameter:
	ShowColors bool `yaml:"showColors"`
}

// If your widget does not need any extra parameters at the moment, maybe it
// still makes sense to wrap `wtf.BaseWidgetConfig` into a custom struct -
// this will make your life easier, if you realise you need some extra configs
// in future.

// -------------------------------- WIDGET -------------------------------- //
// Any widget must implement the `wtf.Widget` interface.
//
// You have 3 ways to achieve this:
//     - you can implement the whole interface on your own
//     - you can compose your widget with `wtf.BaseWidget` and implement
//       missing `Refresh()` and `Close()` methods, and specify which
//       implementation of `wtf.View` your widget is going to use
//     - you can compose your widget with `wtf.TextWidget` or `wtf.TableWidget`
//       and only implement missing `Refresh()` and `Close()` (these widgets
//       already have the `wtf.View` implementations)
//
// Let's start from TextWidget, because we just want to display some text:
type Widget struct {
	// Compose the TextWidget into your struct
	*wtf.TextWidget

	// Let's keep the config inside the widget, to make it accessible in all
	// methods.
	config *Config

	// Just as an example, we are going to inject one of the app service and
	// use it in our widget:
	logger wtf.Logger
}

// Now we need to implement all the missing required logic. Because we use
// the TextWidget, we already have everything except Refresh and Close methods.
// Let's implement them:

// Refresh is responsible for refreshing the widget content. It runs for the
// first time when the widget is initialised and every time when refresh
// scheduler issues a refreshing event.
func (widget *Widget) Refresh() {
	// Here is an example how we can use the logger service
	widget.logger.Debug("HelloWorld: refresh")

	// Now you can do any steps to collect the latest data, format it, and
	// send it to the app view.

	// If your widget is small and simple, then you could write everything
	// in this function. But in more complex widgets it's better to separate
	// this logic into different functions, or even different files.
	// For example if your widget is going to do some complex fetching of
	// remote data, and do some complex formatting, then it makes sense to
	// create 2 extra files:
	//     - data.go (just an example name) will be responsible for defining
	//       data structs and fetching data from a remote source.
	//     - display.go will contain e.g "display()" method of our widget,
	//       which will format the data structs into text and display it.

	// We are writing just a simple "Hello World", so we can do it in this
	// method.

	// Let's pretend to do some text formatting
	content := "Hello World"
	// Let's colorize our content with configured color, if colors are
	// enabled in our configuration
	if widget.config.ShowColors {
		// All tview primitives, which are used as view implementations in
		// WTF app, support color tags in their content.
		// For example "[red]foo" will colorize foo to red color.
		// See https://godoc.org/github.com/rivo/tview for more information.

		color := widget.config.ContentColor
		content = fmt.Sprintf("[%s]%s[-]", color, content)
	}

	// The way how you display the content depends on which kind of view
	// you have - we use TextWidget which uses tview.TextView which provides
	// SetText method.
	widget.TextView.SetText(content)
}

// Close is triggered when the widget is going to be detached from the app.
func (widget *Widget) Close() error {
	// If your widget have any open files or connections, then you can close
	// then here.
	return nil
}

// ------------------------------ CONSTRUCTOR ----------------------------- //
// In order to be able to register our widget in a WTF application, we need
// to provide a specific constructor function, which matches type wtf.WidgetConstructor.
// The function should create and initialise a new instance of our widget,
// or return an error if something went wrong. The function accepts two
// arguments: "configure" is a function which we gonna use to unmarshal yaml
// configs into our Config struct, "app" is a container with some global
// configs, services and functions, which are shared between all widgets.
func New(configure wtf.UnmarshalFunc, app *wtf.AppContext) (wtf.Widget, error) {
	// Initialise a new widget instance
	widget := &Widget{}

	// Initialise a new config instance with default configs
	widget.config = &Config{
		ContentColor: wtf.Color("green"),
		ShowColors:   true,
	}
	// Use the wtf.UnmarshalFunc to apply configs from config file on top of
	// the default configs
	if err := configure(widget.config); err != nil {
		return nil, err
	}

	// Now it's time to initialise our composed widget implementation.
	// All available types of basic widgets as well as widget traits can be
	// created by the "app" container.
	//
	// Let's create a new TextWidget instance.
	title := "GitHub" // will be displayed in the widget's header
	focusable := true // make widget focusable (you can change focus by Tab)
	widget.TextWidget = app.TextWidget(title, widget.config, focusable)

	// OUr widget expects to have a logger, which is provided in the "app"
	// container. Let's pass it:
	widget.logger = app.Logger
	// Check the wtf.AppContext to learn more about which services and functions
	// are available there.

	// Here you can also define any extra initialising steps, e.g. opening
	// files, connections, etc.

	return widget, nil
}

// Now the widget is ready to be registered in the app:
//
//     package main
//     import "github.com/senorprogrammer/wtf/widget/hello-world"
//
//     // considering you have an initialised WTF application
//     var app wtf.App
//
//     //you can add your widget
//     app.RegisterWidget("hello-world", hello_world.New)

// Finally, you just need to configure it in the config:
//     widgets:
//
//         # note: the type here matches the string which you specified while
//         # registering your widget in the app. That is how the app will be
//         # able to realise that this config belongs to your widget.
//       - type: hello-world
//
//         # all widgets are enabled by default, so you do not need to enable
//         # it explicitly, but you can set it to false to disable it.
//         enabled: true
//
//         # set the refresh (in seconds) to -1 to disable refreshing (we have
//         # just static content, does not make sense to refresh it)
//         refreshInterval: -1
//
//         # defile how the widget is displayed - first column, second row,
//         # two columns wide, three rows tall
//         position:
//           left: 0
//           top: 1
//           width: 2
//           height: 3
//
//         # define our custom config params
//         contentColor: yellow
//         showColors: true

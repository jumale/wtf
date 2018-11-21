##### This document contains explanations about the changes which have been done in order to improve the app customization and creating new modules

## Main change notes
- migrated vendors to "go mod"

- changed the config structure: now there is "app" and "widgets" config on the root level

- now it's possible to configure and enable several widgets of the same type (e.g. 2 "textfile" widgets)

- you can configure rows and columns in grid by specifying just a desired number of rows/columns, and their sizes will be 
  automatically calculated based on the terminal window size

- you can specify on configuration level a focus-char for focusable elements. It allows you to customize the way how 
  would you switch between focusable widgets

- by default there is "default" color set for window/app-view/widget background, as well as default text color. Which 
  makes WTF matching your default terminal view (e.g. if a terminal is configured with a half-transparent background, 
  then WTF application will get it by default. This potentially can cause some artifacts while refreshing dynamic content
  in some widgets, but for such cases we always have override the background in config.yml with e.g. black color.
  We just need to let users know about these options and let them decide which behaviour to choose

- the whole config is described in structs, which makes it more discoverable

- no globals, all variables accessed by arguments, or via new AppContext struct

- all widgets moved to their own new directory (before it was quite hard to understand if a folder contains a widget inside,
  or some util package or something else)

- the whole app instantiation is implemented in the new App component (main.go file now contains much less code)

- refactored the way how colors work: now most of responsibilities for colors conversion are moved to tcell library, also
  created a new Color type which is convertable between different color formats.
  
- re-defined widget interfaces and widgets' responsibilities, defined different levels of base widget implementations

- BarGraph: data structure `[][2]int64` replaced with interface `[]GraphEntry`, which provides more flexible way of 
  displaying different kinds of data

- separated the logger implementation from the logger widget. Now logger service in the "wtf" package is responsible for 
  logging into a file, while the logger widget is just responsible for displaying those logs.

## New structure

### main.go

[The main function](main.go) is responsible for:
- defining initial arguments: where the application directory is, and which config should be loaded, enable or disable
  the config watcher
- instantiating a new App
- registering all needed widgets, and defining by which type-names those widgets are recognisable in the config 

### Application

[The application](wtf/app.go) is responsible for:
- loading configs from a config file using [the config loader](wtf/config_loader.go)
- initialising all the services:
    - [file system utils](wtf/service_fs.go)
    - [text formatter](wtf/service_formatter.go)
    - [logger](wtf/service_logger.go)
    - [display](wtf/display.go)
    - [focus tracker](wtf/focus_tracker.go)
    - [scheduled app refresher](wtf/refresher.go)
- registering and initialising widgets
- initialising the app grid and the display layout
- connecting all those together and finally running the app  

### AppContext

While the application is the guy who knows everything about everyone, [the AppContext](wtf/app_context.go) is a container
with a limited amount of functions, configs and services, which is passed from the application level to every widget in
order to provided access to the shared values. It should provide a read-only access, and should not share any application
specific components to not violate a limited responsibility of widgets.

### Config

[The config](wtf/config.go) is now completely defined in structs with definitions for unmarshaling from "yaml".
The config consist of 3 different kinds of configs:
- application specific configs: those a config parameters, which define how the whole application looks and behaves
- widget specific configs: the config parameters which define a generic configuration for every widget
- generic configs: some config-structs which are reusable in both widget's and app's configs.

[The colors](wtf/colors.go) now provide a specific `Color` type used everywhere in configurations, which represents
a color label (e.g "red", "green") and knows how to convert itself from/to `tcell` constants.

[There is a default app config](wtf/config.go) which is used to fulfill all missing values in the config yaml.

[There is also a sample](wtf/config_sample.go) of the yaml config, it is used by [the config loader](wtf/config_loader.go)
to create a sample config file, if it's not created yet.  

### Widget

[The widget](wtf/widget.go) is a type which aggregates a bunch of interfaces, required by different parts of the application:
- [the viewer interface](wtf/viewer.go) tells us that the widget can provide its view implementation
- [the enabler interface](wtf/enabler.go) tells us if the widget enabled or not
- [the focuser interface](wtf/focus_tracker.go) connects the widget to the focus tracker
- [the positioner interface](wtf/position.go) connects the widget to the [display](wtf/display.go)
- [the refresher interface](wtf/refresher.go) connects the widget to the scheduled app refresher
- and finally a Closer interface tells the widget to properly react on removing from the app

There are also the definitions of the main widget config, and a widget constructor (a function which must be implemented
by every widget and which defines how to create an initialise a new instance of the widget)

#### Base Widget

[The base widget](wtf/widget_base.go) provides a very basic implementation of the Widget interface. It implements only
those functions which are common for any kind of widget, also it performs a basic view configuration. The base widget
can be extended by a more high level implementation which provides the view and the way how the widget displays itself.

#### View Specific widgets

The widgets like [TextWidget](wtf/widget_text.go) or [TableWidget](wtf/widget_table.go) provide a higher level of the
widget implementation. They provide an additional information about what kind of view.

#### Widget Traits

Some types, like [HelpfulWidgetTrait](wtf/widget_helpful.go) or [MultiSourceWidgetTrait](wtf/widget_multisource.go) do
not implement themselves any part of the `Widget` interface, but instead they provide some additional set of methods to
extend functionality of the widgets.

### Creating New Widgets

Creating a new widget is a trivial task, but still requires a few steps.

Check [the hello-world widget](widget/hello-world) as an example of how to create a new widget. There you will find
detailed explanations of every step.

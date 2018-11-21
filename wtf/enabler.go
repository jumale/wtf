package wtf

// Represents a widget item which can tell if it is enabled or not.
// Disabled widgets will be skipped from displaying as well as refreshing.
type Enabler interface {
	Enabled() bool
}

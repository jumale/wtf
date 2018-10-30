package app

import "github.com/senorprogrammer/wtf/wtf"

type WidgetsList []wtf.Widget

func (list *WidgetsList) add(w wtf.Widget) {
	*list = append(*list, w)
}

func (list WidgetsList) asFocusable() []wtf.Focusable {
	var result []wtf.Focusable
	for _, widget := range list {
		result = append(result, widget)
	}
	return result
}

func (list WidgetsList) asRefreshable() []wtf.Refreshable {
	var result []wtf.Refreshable
	for _, widget := range list {
		result = append(result, widget)
	}
	return result
}

func (list WidgetsList) asDisplayable() []wtf.Displayable {
	var result []wtf.Displayable
	for _, widget := range list {
		result = append(result, widget)
	}
	return result
}

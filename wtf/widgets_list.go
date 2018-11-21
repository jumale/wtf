package wtf

type WidgetsList struct {
	items []Widget
}

func (list *WidgetsList) add(w Widget) {
	list.items = append(list.items, w)
}

func (list WidgetsList) enabled() []Widget {
	var result []Widget
	for _, item := range list.items {
		if item.Enabled() {
			result = append(result, item)
		}
	}
	return result
}

func (list WidgetsList) asFocusable() []FocusTrackedItem {
	var result []FocusTrackedItem
	for _, widget := range list.items {
		result = append(result, widget)
	}
	return result
}

func (list WidgetsList) asRefreshable() []Refreshable {
	var result []Refreshable
	for _, widget := range list.items {
		result = append(result, widget)
	}
	return result
}

func (list WidgetsList) asDisplayable() []Displayable {
	var result []Displayable
	for _, widget := range list.items {
		result = append(result, widget)
	}
	return result
}

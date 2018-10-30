package wtf

import (
	"github.com/rivo/tview"
)

type TableWidget struct {
	enabled   bool
	focusable bool
	focusChar string

	Name       string
	RefreshInt int
	TableView  *tview.Table

	Position
}

//func NewTableWidget(app *tview.Application, name string, configKey string, focusable bool) *TableWidget {
//	widget := &TableWidget{
//		enabled:   AppConfig.UBool(fmt.Sprintf("wtf.mods.%s.enabled", configKey), false),
//		focusable: focusable,
//
//		Name:       AppConfig.UString(fmt.Sprintf("wtf.mods.%s.title", configKey), name),
//		RefreshInt: AppConfig.UInt(fmt.Sprintf("wtf.mods.%s.refreshInterval", configKey)),
//	}
//
//	widget.Position = NewPosition(
//		AppConfig.UInt(fmt.Sprintf("wtf.mods.%s.position.top", configKey)),
//		AppConfig.UInt(fmt.Sprintf("wtf.mods.%s.position.left", configKey)),
//		AppConfig.UInt(fmt.Sprintf("wtf.mods.%s.position.width", configKey)),
//		AppConfig.UInt(fmt.Sprintf("wtf.mods.%s.position.height", configKey)),
//	)
//
//	widget.addView(app, configKey)
//
//	return widget
//}
//
///* -------------------- Exported Functions -------------------- */
//
//func (widget *TableWidget) BorderColor() string {
//	if widget.Focusable() {
//		return AppConfig.UString("wtf.colors.border.focusable", "red")
//	}
//
//	return AppConfig.UString("wtf.colors.border.normal", "gray")
//}
//
//func (widget *TableWidget) ContextualTitle(defaultStr string) string {
//	if widget.FocusChar() == "" {
//		return fmt.Sprintf(" %s ", defaultStr)
//	}
//
//	return fmt.Sprintf(" %s [darkgray::u]%s[::-][green] ", defaultStr, widget.FocusChar())
//}
//
//func (widget *TableWidget) Disable() {
//	widget.enabled = false
//}
//
//func (widget *TableWidget) Disabled() bool {
//	return !widget.Enabled()
//}
//
//func (widget *TableWidget) Enabled() bool {
//	return widget.enabled
//}
//
//func (widget *TableWidget) Focusable() bool {
//	return widget.enabled && widget.focusable
//}
//
//func (widget *TableWidget) FocusChar() string {
//	return widget.focusChar
//}
//
//func (widget *TableWidget) RefreshInterval() int {
//	return widget.RefreshInt
//}
//
//func (widget *TableWidget) SetFocusChar(char string) {
//	widget.focusChar = char
//}
//
//func (widget *TableWidget) View() View {
//	return widget.TableView
//}
//
///* -------------------- Unexported Functions -------------------- */
//
//
//func (widget *TableWidget) addView(app *tview.Application, configKey string) {
//	view := tview.NewTable()
//
//	view.SetBackgroundColor(colorFor(
//		AppConfig.UString(fmt.Sprintf("wtf.mods.%s.colors.background", configKey),
//			AppConfig.UString("wtf.colors.background", "black"),
//		),
//	))
//
//	view.SetTitleColor(colorFor(
//		AppConfig.UString(
//			fmt.Sprintf("wtf.mods.%s.colors.title", configKey),
//			AppConfig.UString("wtf.colors.title", "white"),
//		),
//	))
//
//	view.SetBorder(true)
//	view.SetBorderColor(colorFor(widget.BorderColor()))
//	view.SetTitle(widget.ContextualTitle(widget.Name))
//
//	widget.TableView = view
//}

package checklist

// ChecklistItem is a module for creating generic checklist implementations
// See 'Todo' for an implementation example
type ChecklistItem struct {
	CheckedIcon string
	Checked     bool
	Text        string
}

// CheckMark returns the string used to indicate a ChecklistItem is checked or unchecked
func (item *ChecklistItem) CheckMark() string {
	if item.Checked {
		return item.CheckedIcon
		//return wtf.AppConfig.UString("wtf.mods.todo.checkedIcon", "x")
	}

	return " "
}

// Toggle changes the checked state of the ChecklistItem
// If checked, it is unchecked. If unchecked, it is checked
func (item *ChecklistItem) Toggle() {
	item.Checked = !item.Checked
}

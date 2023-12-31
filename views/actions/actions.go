package actions

import (
	"github.com/z3phyro/termui"
	"github.com/z3phyro/termui/widgets"
)

var wActionsList *widgets.List
var focused bool = false

func InitActions() *widgets.List {
	wActionsList = widgets.NewList()

	wActionsList.Title = " Requests "
	wActionsList.SelectedRowStyle = termui.NewStyle(termui.ColorGreen)

	return wActionsList
}

func ToggleFocus() {
	focused = !focused

	if focused {
		wActionsList.Block.BorderStyle = termui.NewStyle(termui.ColorGreen)
		wActionsList.TitleStyle = termui.NewStyle(termui.ColorGreen, termui.ColorClear, termui.ModifierBold)
	} else {
		wActionsList.Block.BorderStyle = termui.NewStyle(termui.ColorWhite)
		wActionsList.TitleStyle = termui.NewStyle(termui.ColorWhite)
	}
}

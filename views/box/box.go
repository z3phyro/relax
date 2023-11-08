package box

import (
	"github.com/z3phyro/termui"
	"github.com/z3phyro/termui/widgets"
)

var (
	wTextBox *widgets.ScrollBox
	focused  bool
)

func InitBox() *widgets.ScrollBox {
	wTextBox = widgets.NewScrollBox()

	wTextBox.Title = " Request "
	wTextBox.WrapText = true

	return wTextBox
}

func SetTitleAndContent(file string, text string) {
	wTextBox.Text = text
	wTextBox.ScrollYPosition = 0
	wTextBox.Title = file
}

func ToggleFocus() {
	focused = !focused

	if focused {
		wTextBox.Block.BorderStyle = termui.NewStyle(termui.ColorGreen)
		wTextBox.TitleStyle = termui.NewStyle(termui.ColorGreen, termui.ColorClear, termui.ModifierBold)
	} else {
		wTextBox.Block.BorderStyle = termui.NewStyle(termui.ColorWhite)
		wTextBox.TitleStyle = termui.NewStyle(termui.ColorWhite)
	}
}

package box

import (
	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

var (
	wTextBox *widgets.Paragraph
	focused  bool
)

func InitBox() *widgets.Paragraph {
	wTextBox = widgets.NewParagraph()

	wTextBox.Title = " Request "
	// wTextBox.WrapText = true

	return wTextBox
}

func SetTitleAndContent(file string, text string) {
	wTextBox.Text = text
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

package explorer

import (
	"os"
	"strings"

	"github.com/gizak/termui/v3"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type nodeValue string

var (
	focused   bool
	wFileTree *widgets.Tree
)

func (nv nodeValue) String() string {
	return string(nv)
}

func GetNodes(route string) ([]*widgets.TreeNode, error) {
	nodes := []*widgets.TreeNode{}

	files, err := os.ReadDir(route)

	if err != nil {
		return nodes, err
	}

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".rest") {
			continue
		}

		nodes = append(nodes, &widgets.TreeNode{
			Value: nodeValue(file.Name()),
		})
	}

	return nodes, nil
}

func InitTree(route string) *widgets.Tree {
	nodes, _ := GetNodes(route)

	wFileTree = widgets.NewTree()
	wFileTree.WrapText = false
	wFileTree.SetNodes(nodes)
	wFileTree.Title = " Requests "
	wFileTree.SelectedRowStyle = termui.NewStyle(ui.ColorGreen)
	wFileTree.TextStyle = termui.NewStyle(termui.ColorWhite)

	return wFileTree
}

func ToggleFocus() {
	focused = !focused

	if focused {
		wFileTree.Block.BorderStyle = termui.NewStyle(termui.ColorGreen)
		wFileTree.TitleStyle = termui.NewStyle(termui.ColorGreen, termui.ColorClear, termui.ModifierBold)
	} else {
		wFileTree.Block.BorderStyle = termui.NewStyle(termui.ColorWhite)
		wFileTree.TitleStyle = termui.NewStyle(termui.ColorWhite)
	}
}

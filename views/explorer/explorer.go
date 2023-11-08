package explorer

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/gizak/termui/v3"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type nodeValue string

type NodeLocation struct {
	Node     *widgets.TreeNode
	Location string
}

var (
	focused       bool
	wFileTree     *widgets.Tree
	NodeLocations []NodeLocation
)

func (nv nodeValue) String() string {
	return string(nv)
}

func GetNodes(route string, nodeLocations *[]NodeLocation) ([]*widgets.TreeNode, error) {
	nodes := []*widgets.TreeNode{}

	files, err := os.ReadDir(route)

	if err != nil {
		return nodes, err
	}

	for _, file := range files {
		if file.IsDir() {
			fullPath := filepath.Join(route, file.Name())
			if dirFiles, errDir := GetNodes(fullPath, nodeLocations); errDir == nil {
				if len(dirFiles) > 0 {
					dirNode := &widgets.TreeNode{
						Value: nodeValue(file.Name()),
						Nodes: dirFiles,
					}
					nodes = append(nodes, dirNode)
				}
			}
		}

		if !strings.HasSuffix(file.Name(), ".rest") {
			continue
		}

		node := &widgets.TreeNode{
			Value: nodeValue(file.Name()),
		}

		*nodeLocations = append(*nodeLocations, NodeLocation{
			Node:     node,
			Location: filepath.Join(route, file.Name()),
		})
		nodes = append(nodes, node)
	}

	return nodes, nil
}

func InitTree(route string) *widgets.Tree {
	NodeLocations = []NodeLocation{}
	nodes, _ := GetNodes(route, &NodeLocations)

	wFileTree = widgets.NewTree()
	wFileTree.WrapText = false
	wFileTree.SetNodes(nodes)
	wFileTree.Title = " Files "
	wFileTree.SelectedRowStyle = termui.NewStyle(ui.ColorGreen)
	wFileTree.TextStyle = termui.NewStyle(termui.ColorWhite)

	return wFileTree
}

func GetNodeLocation(node *widgets.TreeNode) string {
	for _, n := range NodeLocations {
		if n.Node == node {
			return n.Location
		}
	}

	return ""
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

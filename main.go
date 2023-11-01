package main

import (
	"fmt"
	"log"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	tb "github.com/nsf/termbox-go"
	"stoicdynamics.com/relax/config"
	"stoicdynamics.com/relax/controllers/client"
	"stoicdynamics.com/relax/controllers/parser"
	"stoicdynamics.com/relax/views/actions"
	box "stoicdynamics.com/relax/views/box"
	fe "stoicdynamics.com/relax/views/explorer"
)

var (
	tabFocus     int
	wFileTree    *widgets.Tree
	wActionsList *widgets.List
	wMainBox     *widgets.Paragraph
)

func refreshTabFocus() {
	switch tabFocus {
	case 0:
		fe.ToggleFocus()
	case 1:
		actions.ToggleFocus()
	case 2:
		box.ToggleFocus()
	}
}

func switchFocus(tab int) {
	refreshTabFocus()
	tabFocus = tab
	refreshTabFocus()
}

func openFile() {
	boxFile := wFileTree.SelectedNode().Value.String()
	switch tabFocus {
	case 0:
		boxContent := parser.OpenFile(config.GetRootPath(), boxFile)
		box.SetTitleAndContent(fmt.Sprintf(" %s ", boxFile), boxContent)
		parser.ParseRequestText(boxContent)

		actionsList := []string{}
		for _, request := range parser.Requests {
			actionsList = append(actionsList, request.Name)
		}
		wActionsList.Rows = actionsList
	case 1:
		request := parser.Requests[wActionsList.SelectedRow]
		box.SetTitleAndContent(fmt.Sprintf(" %s - %s ", request.Name, boxFile), request.Raw)
	}

	ui.Render(wActionsList)
	ui.Render(wMainBox)
}

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	// To allow Mouse Events
	tb.SetInputMode(tb.InputEsc)

	config.InitConfig()

	grid := ui.NewGrid()

	wFileTree = fe.InitTree(config.GetRootPath())
	wMainBox = box.InitBox()
	wActionsList = actions.InitActions()

	tabFocus = 0

	grid.Set(
		ui.NewCol(1.0/4,
			ui.NewRow(1.0/2, wFileTree),
			ui.NewRow(1.0/2, wActionsList)),
		ui.NewCol(3.0/4, wMainBox),
	)

	resize := func() {
		x, y := ui.TerminalDimensions()
		grid.SetRect(0, 0, x, y)
	}

	resize()

	openFile()
	fe.ToggleFocus()
	ui.Render(grid)

	previousKey := ""
	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		case "j", "<Down>":
			switch tabFocus {
			case 0:
				wFileTree.ScrollDown()
			case 1:
				wActionsList.ScrollDown()
			}
			openFile()
		case "k", "<Up>":
			switch tabFocus {
			case 0:
				wFileTree.ScrollUp()
			case 1:
				wActionsList.ScrollUp()
			}
			openFile()
		case "g":
			if previousKey == "g" {
				wFileTree.ScrollTop()
			}
		case "<Resize>":
			resize()
		case "1":
			switchFocus(0)
			openFile()
		case "2":
			switchFocus(1)
			openFile()
		case "3":
			switchFocus(2)
			openFile()
		case "<Space>", "l":
			switch tabFocus {
			case 0:
				switchFocus(1)
			case 1:
				request := parser.Requests[wActionsList.SelectedRow]
				response := client.MakeRequest(request)

				text := fmt.Sprintf("%s\nRESPONSE BODY\n\n%s", request.Raw, parser.ParseResponse(response))
				box.SetTitleAndContent(request.Name, text)
				switchFocus(2)
			}
		case "h":
			refreshTabFocus()
			tabFocus -= 1
			if tabFocus < 0 {
				tabFocus = 2
			}
			refreshTabFocus()
		case "<Tab>":
			refreshTabFocus()
			tabFocus += 1
			if tabFocus > 2 {
				tabFocus = 0
			}
			refreshTabFocus()
			if tabFocus == 1 {
				openFile()
			}
		}

		if previousKey == "g" {
			previousKey = ""
		} else {
			previousKey = e.ID
		}

		ui.Render(grid)
	}
}

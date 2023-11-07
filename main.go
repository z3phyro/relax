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
	"stoicdynamics.com/relax/views/history"
)

var (
	wFileTree        *widgets.Tree
	wActionsList     *widgets.List
	wMainBox         *widgets.ScrollBox
	wHistoryList     *widgets.List
	grid             *ui.Grid
	isHistoryVisible bool = true
	viewsChanged     bool = false
	tabFocus         int
)

func refreshTabFocus() {
	switch tabFocus {
	case 0:
		fe.ToggleFocus()
	case 1:
		actions.ToggleFocus()
	case 2:
		box.ToggleFocus()
	case 3:
		history.ToggleFocus()
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
		parser.ParseRequestText(boxContent, boxFile)

		actionsList := []string{}
		for _, request := range parser.Requests {
			actionsList = append(actionsList, request.Name)
		}
		wActionsList.Rows = actionsList
	case 1:
		request := parser.Requests[wActionsList.SelectedRow]
		box.SetTitleAndContent(fmt.Sprintf(" %s - %s ", request.Name, boxFile), request.Raw)
	case 3:
		if wHistoryList.SelectedRow >= len(history.RequestsHistory) {
			break
		}

		box.SetTitleAndContent(
			fmt.Sprintf(" %s ", wHistoryList.Rows[wHistoryList.SelectedRow]),
			history.RequestsHistory[wHistoryList.SelectedRow].Response,
		)
	}
}

func render() {
	if isHistoryVisible {
		grid.Set(
			ui.NewCol(1.0/5,
				ui.NewRow(1.0/2, wFileTree),
				ui.NewRow(1.0/2, wActionsList)),
			ui.NewCol(3.0/5, wMainBox),
			ui.NewCol(1.0/5, wHistoryList),
		)
	} else {
		grid.Set(
			ui.NewCol(1.0/4,
				ui.NewRow(1.0/2, wFileTree),
				ui.NewRow(1.0/2, wActionsList)),
			ui.NewCol(3.0/4, wMainBox),
		)
	}

	ui.Render(grid)
	ui.Render(wHistoryList)
	ui.Render(wMainBox)
	ui.Render(wFileTree)
	ui.Render(wActionsList)

	viewsChanged = false
}

func resize() {
	x, y := ui.TerminalDimensions()
	grid.SetRect(0, 0, x, y)
}

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	// To allow Mouse Events
	tb.SetInputMode(tb.InputEsc)

	config.InitConfig()

	grid = ui.NewGrid()

	wFileTree = fe.InitTree(config.GetRootPath())
	wMainBox = box.InitBox()
	wActionsList = actions.InitActions()
	wHistoryList = history.InitHistory()

	tabFocus = 0

	openFile()
	fe.ToggleFocus()
	resize()
	render()
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
			case 2:
				wMainBox.ScrollDown()
			case 3:
				wHistoryList.ScrollDown()
			}
			openFile()
		case "k", "<Up>":
			switch tabFocus {
			case 0:
				wFileTree.ScrollUp()
			case 1:
				wActionsList.ScrollUp()
			case 2:
				wMainBox.ScrollUp()
			case 3:
				wHistoryList.ScrollUp()
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
		case "<Space>":
			switch tabFocus {
			case 0:
				switchFocus(1)
			case 1, 2:
				if wActionsList.SelectedRow >= len(parser.Requests) {
					break
				}
				request := parser.Requests[wActionsList.SelectedRow]
				response, err := client.MakeRequest(request)
				var text string
				if err != nil {
					text = err.Error()
				} else {
					text = fmt.Sprintf(
						"%s \n\nResponse:\nStatus %s\n\nBody: \n%s\n\nHeaders:\n%s",
						request.Raw,
						response.Status,
						parser.ParseResponse(response.Body),
						response.Header,
					)
				}
				history.LogRequest(request, text)

				box.SetTitleAndContent(fmt.Sprintf("%s - %s", request.Name, request.FileName), text)
			}
		case "l":
			newPos := tabFocus + 1
			if newPos > 3 {
				newPos = 0
			}
			switchFocus(newPos)
		case "h":
			newPos := tabFocus - 1
			if newPos < 0 {
				newPos = 3
			}
			switchFocus(newPos)
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
		case "y":
			if isHistoryVisible {
				isHistoryVisible = false
			} else {
				isHistoryVisible = true
			}
			viewsChanged = true
		}

		if previousKey == "g" {
			previousKey = ""
		} else {
			previousKey = e.ID
		}

		render()
	}
}
